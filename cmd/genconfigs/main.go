package main

import (
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
)

// Configs bundles all instance-individual
// configurations into one slice.
type Configs struct {
	Configs []Config `json:"Configs"`
}

// Config describes one compute instance
// exhaustively for reproducibility.
type Config struct {
	Name              string   `json:"Name"`
	Zone              string   `json:"Zone"`
	MachineType       string   `json:"MachineType"`
	Subnet            string   `json:"Subnet"`
	NetworkTier       string   `json:"NetworkTier"`
	MinCPUPlatform    string   `json:"MinCPUPlatform"`
	Scopes            []string `json:"Scopes"`
	SourceSnapshot    string   `json:"SourceSnapshot"`
	BootDiskSize      string   `json:"BootDiskSize"`
	BootDiskType      string   `json:"BootDiskType"`
	Disk              string   `json:"Disk"`
	MaintenancePolicy string   `json:"MaintenancePolicy"`
	Flags             []string `json:"Flags"`
}

func main() {

	// Allow for control via command-line flags.
	configsPathFlag := flag.String("configsPath", "./gcloud-configs/", "Specify location where GCloud Compute configurations are supposed to be saved.")
	flag.Parse()

	// Extract parsed flag values.
	configsPath, err := filepath.Abs(*configsPathFlag)
	if err != nil {
		fmt.Printf("Provided configs path '%s' could not be converted to absolute path: %v\n", *configsPathFlag, err)
		os.Exit(1)
	}

	// All GCloud zones.
	zones := []string{
		"asia-east1",
		"asia-east2",
		"asia-northeast1",
		"asia-northeast2",
		"asia-south1",
		"asia-southeast1",
		"australia-southeast1",
		"europe-north1",
		"europe-west1",
		"europe-west2",
		"europe-west3",
		"europe-west4",
		"europe-west6",
		"northamerica-northeast1",
		"southamerica-east1",
		"us-central1",
		"us-east1",
		"us-east4",
		"us-west1",
		"us-west2",
	}

	// Shuffle zones slice.
	for i := (len(zones) - 1); i > 0; i-- {

		// Generate new CSPRNG number smaller than i.
		jBig, err := rand.Int(rand.Reader, big.NewInt(int64(i)))
		if err != nil {
			fmt.Printf("Drawing random number failed: %v\n", err)
			os.Exit(1)
		}
		j := int(jBig.Int64())

		// Swap places i and j in zones slice.
		zones[i], zones[j] = zones[j], zones[i]
	}

	// Prepare slice of configuration lines
	// of desired size.
	configs := Configs{
		Configs: make([]Config, 100),
	}

	for i := 0; i < 100; i++ {

		// Base machine specification on a distribution
		// approximating real-life hardware availability.
		var machineType string

		switch {
		case i >= 0 && i < 20:
			machineType = "f1-micro"

		case i >= 20 && i < 60:
			machineType = "f1-micro"

		case i >= 60 && i < 90:
			machineType = "f1-micro"

		case i >= 90 && i < 100:
			machineType = "f1-micro"
		}

		// Prefill all configuration lines.
		configs.Configs[i] = Config{
			Name:           fmt.Sprintf("mixnet.%04d", (i + 1)),
			MachineType:    machineType,
			Subnet:         "default",
			NetworkTier:    "STANDARD",
			MinCPUPlatform: "Intel Skylake",
			Scopes: []string{
				"https://www.googleapis.com/auth/servicecontrol",
				"https://www.googleapis.com/auth/service.management.readonly",
				"https://www.googleapis.com/auth/logging.write",
				"https://www.googleapis.com/auth/monitoring.write",
				"https://www.googleapis.com/auth/trace.append",
				"https://www.googleapis.com/auth/devstorage.full_control",
			},
			SourceSnapshot:    "mixnet-base",
			BootDiskSize:      "10GB",
			BootDiskType:      "pd-ssd",
			Disk:              "name=mixnet,device-name=mixnet,mode=rw,boot=yes,auto-delete=yes",
			MaintenancePolicy: "TERMINATE",
			Flags: []string{
				"--no-restart-on-failure",
			},
		}
	}

	for i := 0; i < 5; i++ {

		// Assign 5 configs to each zone.
		for j := 0; j < len(zones); j++ {
			configs.Configs[((i * len(zones)) + j)].Zone = zones[j]
		}
	}

	// Marshal slice of configs to JSON.
	configsJSON, err := json.MarshalIndent(configs, "", "\t")
	if err != nil {
		fmt.Printf("Failed to marshal configurations to JSON: %v\n", err)
		os.Exit(1)
	}

	// Write to file.
	err = ioutil.WriteFile(filepath.Join(configsPath, "gcloud-mixnet-20-40-30-10.json"), configsJSON, 0644)
	if err != nil {
		fmt.Printf("Error writing configurations in JSON format to file: %v\n", err)
		os.Exit(1)
	}
}
