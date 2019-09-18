package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

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
