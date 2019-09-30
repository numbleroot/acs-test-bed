package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/emicklei/go-restful"
	"github.com/numbleroot/acs-test-bed/cmd/operator/zenopki"
)

// Operator describes the node in the
// public cloud that orchestrates the
// execution of experiments from start
// to finish.
type Operator struct {
	sync.Mutex

	GCloudServiceAcc  string
	GCloudProject     string
	GCloudBucket      string
	GCloudAccessToken string

	TLSCertPath string
	TLSKeyPath  string

	InternalListenAddr   string
	InternalSrv          *restful.WebService
	InternalRegisterChan chan *RegisterReq
	InternalReadyChan    chan string
	InternalFinishedChan chan string
	InternalFailedChan   chan *FailedReq

	PublicListenAddr    string
	PublicSrv           *restful.WebService
	PublicNewChan       chan string
	PublicTerminateChan chan struct{}

	ExpInProgress string
	Exps          map[string]*Exp

	ZenoPKI *zenopki.PKI
}

// Exp contains all information relevant
// for monitoring an experiment.
type Exp struct {
	ID           string             `json:"id"`
	Created      string             `json:"created"`
	System       string             `json:"system"`
	Concluded    bool               `json:"concluded"`
	ResultFolder string             `json:"resultFolder"`
	Progress     []string           `json:"progress"`
	ProgressChan chan string        `json:"-"`
	Servers      []*Worker          `json:"servers"`
	ServersMap   map[string]*Worker `json:"-"`
	Clients      []*Worker          `json:"clients"`
	ClientsMap   map[string]*Worker `json:"-"`
}

// Worker describes one compute instance
// exhaustively for reproducibility.
type Worker struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Address         string `json:"address"`
	Status          string `json:"status"`
	Zone            string `json:"zone"`
	MinCPUPlatform  string `json:"minCPUPlatform"`
	MachineType     string `json:"machineType"`
	TypeOfNode      string `json:"typeOfNode"`
	BinaryName      string `json:"binaryName"`
	SourceImage     string `json:"sourceImage"`
	DiskType        string `json:"diskType"`
	DiskSize        string `json:"diskSize"`
	NetTroubles     string `json:"netTroubles"`
	ZenoMixesKilled int    `json:"zenoMixesKilled"`
}

func init() {

	// Enable TLS 1.3.
	if os.Getenv("GODEBUG") == "" {
		os.Setenv("GODEBUG", "tls13=1")
	} else {
		os.Setenv("GODEBUG", fmt.Sprintf("%s,tls13=1", os.Getenv("GODEBUG")))
	}
}

func main() {

	// Command-line options.
	publicListenAddrFlag := flag.String("publicAddr", "0.0.0.0:20443", "Specify HTTPS address for receiving experiment instructions.")
	internalListenAddrFlag := flag.String("internalAddr", "127.0.0.1:443", "Specify HTTPS address for administrating experiments.")
	gcloudServiceAccFlag := flag.String("gcloudServiceAcc", "", "Supply the GCloud service account to use for the experiments.")
	gcloudProjectFlag := flag.String("gcloudProject", "", "Supply the GCloud project ID to use for the experiments.")
	gcloudBucketFlag := flag.String("gcloudBucket", "", "Supply the GCloud Storage Bucket to use for the experiments.")
	certPathFlag := flag.String("certPath", "/root/operator-cert.pem", "Supply file system location of the operator's TLS certificate.")
	keyPathFlag := flag.String("keyPath", "/root/operator-key.pem", "Supply file system location of the operator's TLS key.")

	flag.Parse()

	if *gcloudServiceAccFlag == "" || *gcloudProjectFlag == "" || *gcloudBucketFlag == "" {
		fmt.Printf("Missing argument(s), please provide values for all flags: '-gcloudServiceAcc', '-gcloudProject', '-gcloudBucket'.\n")
		os.Exit(1)
	}

	op := &Operator{
		GCloudServiceAcc: *gcloudServiceAccFlag,
		GCloudProject:    *gcloudProjectFlag,
		GCloudBucket:     *gcloudBucketFlag,

		TLSCertPath: *certPathFlag,
		TLSKeyPath:  *keyPathFlag,

		InternalListenAddr:   *internalListenAddrFlag,
		InternalRegisterChan: make(chan *RegisterReq),
		InternalReadyChan:    make(chan string),
		InternalFinishedChan: make(chan string),
		InternalFailedChan:   make(chan *FailedReq),

		PublicListenAddr:    *publicListenAddrFlag,
		PublicNewChan:       make(chan string),
		PublicTerminateChan: make(chan struct{}),

		ExpInProgress: "",
		Exps:          make(map[string]*Exp),
	}

	// Create goroutine that completely
	// handles experiment procedure.
	go op.RunExperiments()

	// Prepare and listen for API calls on the
	// internal network endpoint (worker nodes).
	op.PrepareInternalSrv()
	go op.RunInternalSrv()

	// Prepare and listen for API calls on the
	// Internet-facing endpoint (start experiments).
	op.PreparePublicSrv()

	fmt.Printf("[PUBLIC] Listening on https://%s/public/experiments for API calls regarding experiments...\n", op.PublicListenAddr)

	err := http.ListenAndServeTLS(op.PublicListenAddr, op.TLSCertPath, op.TLSKeyPath, nil)
	if err != nil {
		fmt.Printf("Failed handling public experiment requests: %v\n", err)
		os.Exit(1)
	}
}
