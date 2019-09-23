package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/numbleroot/vuvuzela"
)

// MixConf contains all mix-identifying
// fields that a Vuvuzela mix expects to
// find in its local configuration file.
type MixConf struct {
	ServerName string           `json:"ServerName"`
	PublicKey  *vuvuzela.BoxKey `json:"PublicKey"`
	PrivateKey *vuvuzela.BoxKey `json:"PrivateKey"`
	ConvoMu    float64          `json:"ConvoMu"`
	ConvoB     float64          `json:"ConvoB"`
}

// ClientConf contains all client-identifying
// fields that a Vuvuzela client expects to
// find in its local configuration file.
type ClientConf struct {
	MyName       string           `json:"MyName"`
	MyPublicKey  *vuvuzela.BoxKey `json:"MyPublicKey"`
	MyPrivateKey *vuvuzela.BoxKey `json:"MyPrivateKey"`
}

// PKIServer holds all information relevant
// for contacting a server found in the PKI file.
type PKIServer struct {
	Address   string `json:"Address"`
	PublicKey string `json:"PublicKey"`
}

// PKI will get written to a file and thus
// replace an actual node being a PKI.
type PKI struct {
	People      map[string]string    `json:"People"`
	Servers     map[string]PKIServer `json:"Servers"`
	ServerOrder []string             `json:"ServerOrder"`
	EntryServer string               `json:"EntryServer"`
}

func generateVuvuzelaMixConfs(mixes []Config, confsPath string) error {

	// Create configuration files folder
	// if it does not exist.
	err := os.MkdirAll(confsPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create configurations folder %s: %v", confsPath, err)
	}

	for i := range mixes {

		// Generate fresh public-private key pair.
		pubKey, privKey, err := vuvuzela.GenerateBoxKey(rand.Reader)
		if err != nil {
			return fmt.Errorf("failed to generate vuvuzela.BoxKey: %v", err)
		}

		conf := &MixConf{
			ServerName: mixes[i].Name,
			PublicKey:  pubKey,
			PrivateKey: privKey,
			ConvoMu:    300000,
			ConvoB:     13800,
		}

		// Marshal mix config values to JSON.
		data, err := json.MarshalIndent(conf, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to encode fresh mix config to JSON: %v", err)
		}

		// Write formatted JSON to file.
		err = ioutil.WriteFile(filepath.Join(confsPath, fmt.Sprintf("%s.conf", mixes[i].Name)), data, 0600)
		if err != nil {
			return fmt.Errorf("could not write out marshalled mix config: %v", err)
		}
	}

	return nil
}

func generateVuvuzelaClientConfs(clients []Config, confsPath string) error {

	// Create configuration files folder
	// if it does not exist.
	err := os.MkdirAll(confsPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create configurations folder %s: %v", confsPath, err)
	}

	for i := range clients {

		// Generate fresh public-private key pair.
		pubKey, privKey, err := vuvuzela.GenerateBoxKey(rand.Reader)
		if err != nil {
			return fmt.Errorf("failed to generate vuvuzela.BoxKey: %v", err)
		}

		conf := &ClientConf{
			MyName:       clients[i].Name,
			MyPublicKey:  pubKey,
			MyPrivateKey: privKey,
		}

		// Marshal client config values to JSON.
		data, err := json.MarshalIndent(conf, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to encode fresh client config to JSON: %v", err)
		}

		// Write formatted JSON to file.
		err = ioutil.WriteFile(filepath.Join(confsPath, fmt.Sprintf("%s.conf", clients[i].Name)), data, 0600)
		if err != nil {
			return fmt.Errorf("could not write out marshalled client config: %v", err)
		}
	}

	return nil
}

func preparePKIFile(confsPath string) error {

	pki := &PKI{
		People:      make(map[string]string),
		Servers:     make(map[string]PKIServer),
		ServerOrder: make([]string, 0, 3),
	}

	// Walk configurations path and read files.
	err := filepath.Walk(confsPath, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if !strings.Contains(path, ".conf") || strings.Contains(path, "pki.conf") {
			return nil
		}

		// Prepare machine name.
		nodeName := strings.Split(filepath.Base(path), ".conf")[0]

		// Read content of machine configuration.
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		if strings.Contains(nodeName, "client") {

			// Extract public key from file.
			configSplit := strings.Split(string(data), "\"MyPublicKey\": \"")
			pubKey := strings.Split(configSplit[1], "\",")[0]

			// Add client to People list of PKI.
			pki.People[nodeName] = pubKey

		} else {

			// Extract public key from file.
			configSplit := strings.Split(string(data), "\"PublicKey\": \"")
			pubKey := strings.Split(configSplit[1], "\",")[0]

			if strings.Contains(nodeName, "00001") {

				// The first server will be used as the
				// coordinator of the Vuvuzela deployment.
				pki.EntryServer = fmt.Sprintf("ACS_EVAL_INSERT_%s_ADDRESS", nodeName)

			} else {

				// Each other server is a regular,
				// in-order Vuvuzela mix.
				pki.Servers[nodeName] = PKIServer{
					Address:   fmt.Sprintf("ACS_EVAL_INSERT_%s_ADDRESS", nodeName),
					PublicKey: pubKey,
				}

				pki.ServerOrder = append(pki.ServerOrder, nodeName)
			}
		}

		return nil
	})
	if err != nil {
		fmt.Printf("Walking Vuvuzela configurations path failed: %v\n", err)
		os.Exit(1)
	}

	// Marshal PKI structure to JSON.
	pkiJSON, err := json.MarshalIndent(pki, "", "  ")
	if err != nil {
		fmt.Printf("Marshaling PKI structure to JSON failed: %v\n", err)
		os.Exit(1)
	}

	// Write preliminary PKI structure to file.
	err = ioutil.WriteFile(filepath.Join(confsPath, "pki.conf"), pkiJSON, 0644)
	if err != nil {
		fmt.Printf("Writing marshaled PKI structure to file failed: %v\n", err)
		os.Exit(1)
	}

	return nil
}
