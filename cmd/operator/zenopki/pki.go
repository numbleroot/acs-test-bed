package zenopki

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// Endpoint bundles all information sent by a
// node to register under various categories
// with the PKI:
// Category == 0  =>  mix intent,
// category == 1  =>  client.
type Endpoint struct {
	Category       uint8
	Name           string
	PubAddr        string
	PubKey         *[32]byte
	PubCertPEM     []byte
	ContactAddr    string
	ContactCertPEM []byte
}

// PKI maintains the mappings of aliases to
// the public keys they registered with.
type PKI struct {
	Lis              net.Listener
	LisAddr          string
	EpochTicker      *time.Ticker
	AcceptMixRegs    int32
	AcceptClientRegs int32
	MuNodes          *sync.RWMutex
	Nodes            map[string]*Endpoint
	EvalCtrlChan     chan struct{}
}

// SendDataToNode accepts all arguments required
// to securely contact one node in isolation and
// transmit supplied data.
func SendDataToNode(wg *sync.WaitGroup, node *Endpoint, data string) {

	// Create new empty cert pool.
	certPool := x509.NewCertPool()

	// Attempt to add received certificate to pool.
	ok := certPool.AppendCertsFromPEM(node.ContactCertPEM)
	if !ok {
		fmt.Printf("[ZENO PKI] Failed to append PEM certificate to empty pool.\n")
		return
	}

	fmt.Printf("[ZENO PKI] Broadcasting to node %s\n", node.ContactAddr)

	tlsConfig := &tls.Config{
		RootCAs:            certPool,
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS13,
		CurvePreferences:   []tls.CurveID{tls.X25519},
	}

	// Contact node.
	tried := 1
	connWrite, err := tls.Dial("tcp", node.ContactAddr, tlsConfig)
	for err != nil && tried <= 20 {

		fmt.Printf("[ZENO PKI] Failed %d times already to dial node at %s (will try again)\n", tried, node.ContactAddr)

		tried++
		time.Sleep(200 * time.Millisecond)

		connWrite, err = tls.Dial("tcp", node.ContactAddr, tlsConfig)
	}
	if err != nil {
		fmt.Printf("[ZENO PKI] Error dialing node %s (tried %d times): %v\n", node.ContactAddr, tried, err)
		return
	}

	// Send data.
	fmt.Fprintf(connWrite, "%s\n", data)

	wg.Done()
}

// BroadcastData sends out a properly formatted
// representation of either all mix candidates
// for the upcoming epoc, all clients for the
// upcoming epoch, or a simple epoch rotation
// signal to all nodes the PKI is aware of.
func (pki *PKI) BroadcastData(dataToSend string) {

	wg := &sync.WaitGroup{}

	pki.MuNodes.RLock()
	defer pki.MuNodes.RUnlock()

	// Make sure to contact all nodes.
	wg.Add(len(pki.Nodes))

	var data string

	if dataToSend == "mixes" {

		data = "mixes"
		for n := range pki.Nodes {

			// Prepare string containing all candidate mix nodes
			// that registered in time for this epoch.
			if pki.Nodes[n].Category == 0 {
				data = fmt.Sprintf("%s;%s,%s,%x,%x", data, pki.Nodes[n].Name, pki.Nodes[n].PubAddr, *pki.Nodes[n].PubKey, pki.Nodes[n].PubCertPEM)
			}
		}

	} else if dataToSend == "clients" {

		data = "clients"
		for n := range pki.Nodes {

			// Prepare string containing all client nodes
			// that registered in time for this epoch.
			if pki.Nodes[n].Category == 1 {
				data = fmt.Sprintf("%s;%s,%s,%x,%x", data, pki.Nodes[n].Name, pki.Nodes[n].PubAddr, *pki.Nodes[n].PubKey, pki.Nodes[n].PubCertPEM)
			}
		}

	} else if dataToSend == "epoch" {
		data = "epoch;"
	}

	for n := range pki.Nodes {

		node := pki.Nodes[n]

		// Contact node and send data
		// off the hot path.
		go SendDataToNode(wg, node, data)
	}

	// Wait for all background transmissions
	// to finish before returning.
	wg.Wait()

	fmt.Printf("[ZENO PKI] Broadcast finished\n\n")
}

// HandleReq is responsible for handling a
// delegated PKI request. Request structure:
//   CATEGORY PUB_ADDR PKI_ADDR PUB_KEY PUB_CERT
// Category is either 'mixes' or 'clients', the
// second argument is the network address to store
// the fourth (Curve25519 public key as hex) and
// fifth argument (PEM-encoded TLS certificate),
// under. Additionally, as third argument, we
// require an address for the PKI to contact the
// node on.
func (pki *PKI) HandleReq(connWrite net.Conn) {

	decoder := gob.NewDecoder(connWrite)

	// Read and parse registration message.
	var reg Endpoint
	err := decoder.Decode(&reg)
	if err != nil {
		fmt.Printf("[ZENO PKI] Failed decoding registration message from node: %v\n", err)
		return
	}

	success := 1

	pki.MuNodes.Lock()

	if reg.Category == 0 {

		if atomic.LoadInt32(&pki.AcceptMixRegs) == 0 {

			// If the atomic flag is set to allow
			// mix intention registrations, do it.
			pki.Nodes[reg.Name] = &reg
			success = 0

		} else {
			success = 2
		}

	} else if reg.Category == 1 {

		if atomic.LoadInt32(&pki.AcceptClientRegs) == 0 {

			// If the atomic flag is set to allow
			// client registrations, do it.
			pki.Nodes[reg.Name] = &reg
			success = 0

		} else {
			success = 2
		}
	}

	pki.MuNodes.Unlock()

	// Respond to client with status.
	// 0 == success, all good.
	// 1 == failure, no good.
	// 2 == currently not accepting, try again later.
	fmt.Fprintf(connWrite, "%d\n", success)
}

// AcceptRegistrations is the main dispatcher
// function accepting incoming PKI requests
// and setting them up in separate routine.
func (pki *PKI) AcceptRegistrations() {

	for {

		fmt.Printf("[ZENO PKI] Waiting for next registration message.\n")

		// Accept incoming new requests.
		connWrite, err := pki.Lis.Accept()
		if err != nil {
			fmt.Printf("[ZENO PKI] Connection error: %v\n", err)
			continue
		}

		go pki.HandleReq(connWrite)
	}
}

// Run initializes and operates the PKI reduced
// in functionality we use in order to operate zeno.
func (pki *PKI) Run(cert string, key string) {

	// Load TLS server certificate and key.
	tlsCert, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		fmt.Printf("[ZENO PKI] Failed loading TLS certificate and private key: %v\n", err)
		os.Exit(1)
	}

	// Prepare TLS configuration.
	conf := &tls.Config{
		Certificates:           []tls.Certificate{tlsCert},
		InsecureSkipVerify:     false,
		MinVersion:             tls.VersionTLS13,
		CurvePreferences:       []tls.CurveID{tls.X25519},
		SessionTicketsDisabled: true,
	}

	// Listen on specified address for incoming
	// requests over TLS connection.
	pki.Lis, err = tls.Listen("tcp", pki.LisAddr, conf)
	if err != nil {
		fmt.Printf("[ZENO PKI] Error listening for PKI requests on TLS endpoint: %v\n", err)
		os.Exit(1)
	}
	defer pki.Lis.Close()

	fmt.Printf("[ZENO PKI] Listening on %s for PKI requests...\n", pki.LisAddr)

	// Handle incoming mix intentions and
	// client registrations.
	go pki.AcceptRegistrations()

	// Configure to be the time slot building
	// block for all epoch timers.
	epochBrick := 5 * time.Second

	fmt.Printf("[ZENO PKI] Waiting for start signal.\n")

	// Wait for start signal from operator.
	<-pki.EvalCtrlChan

	fmt.Printf("[ZENO PKI] Start signal received!\n")

	for {

		fmt.Printf("\n[ZENO PKI] Mixes and clients can register now\n")

		// First time period: accept declarations of
		// intent by nodes wanting to become mixes.
		pki.EpochTicker = time.NewTicker(2 * epochBrick)

		// Registration closed.
		<-pki.EpochTicker.C

		// Block further mix registrations.
		atomic.StoreInt32(&pki.AcceptMixRegs, 1)

		fmt.Printf("[ZENO PKI] Mixes registration closed, broadcasting...\n")

		// Broadcast candidates to all nodes.
		pki.BroadcastData("mixes")

		// Second time period: all nodes deterministically
		// determine the cascades locally.
		pki.EpochTicker = time.NewTicker(4 * epochBrick)

		// Cascades election done. Also, no more clients
		// can register for the upcoming epoch.
		<-pki.EpochTicker.C

		// Block further client registrations.
		atomic.StoreInt32(&pki.AcceptClientRegs, 1)

		fmt.Printf("[ZENO PKI] Clients registration closed, broadcasting...\n")

		// Inform all nodes about the set of clients.
		pki.BroadcastData("clients")

		// Third time period: regular epoch execution
		// minus the time it takes for the subsequent
		// cascade matrix election to complete.
		pki.EpochTicker = time.NewTicker(500 * epochBrick)

		// Regular epoch finished.
		<-pki.EpochTicker.C

		fmt.Printf("\n[ZENO PKI] Epoch closing, broadcasting...\n")

		// Inform nodes about epoch rotation.
		pki.BroadcastData("epoch")

		// Reset state.
		pki.AcceptMixRegs = 0
		pki.AcceptClientRegs = 0
		pki.MuNodes = &sync.RWMutex{}
		pki.Nodes = make(map[string]*Endpoint)
	}
}
