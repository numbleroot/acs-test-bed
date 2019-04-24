package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
)

// Config describes one compute instance
// exhaustively for reproducibility.
type Config struct {
	Name               string
	Zone               string
	MachineType        string
	Image              string
	ImageProject       string
	Subnet             string
	NetworkTier        string
	Scopes             []string
	MinCPUPlatform     string
	BootDiskSize       string
	BootDiskType       string
	BootDiskDeviceName string
	MaintenancePolicy  string
	Flags              []string
}

// Header returns the field descriptions
// for a configuration line.
func (cfg *Config) Header() string {
	return fmt.Sprintf("Name ; Zone ; MachineType ; Image ; ImageProject ; Subnet ; NetworkTier ; Scopes ; MinCPUPlatform ; BootDiskSize ; BootDiskType ; BootDiskDeviceName ; MaintenancePolicy ; Flags\n")
}

// Data returns the actual string that
// contains all properly formatted data
// from one configuration line.
func (cfg *Config) Data() string {

	scopes := strings.Join(cfg.Scopes, ",")
	flags := strings.Join(cfg.Flags, " ")

	return fmt.Sprintf("%s ; %s ; %s ; %s ; %s ; %s ; %s ; %s ; %s ; %s ; %s ; %s ; %s ; %s\n", cfg.Name, cfg.Zone, cfg.MachineType, cfg.Image, cfg.ImageProject, cfg.Subnet, cfg.NetworkTier, scopes, cfg.MinCPUPlatform, cfg.BootDiskSize, cfg.BootDiskType, cfg.BootDiskDeviceName, cfg.MaintenancePolicy, flags)
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
	configs := make([]*Config, 100)

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
		configs[i] = &Config{
			Name:         fmt.Sprintf("mixnet.%04d", (i + 1)),
			MachineType:  machineType,
			Image:        "ubuntu-1804-bionic-v20190404",
			ImageProject: "ubuntu-os-cloud",
			Subnet:       "default",
			NetworkTier:  "STANDARD",
			Scopes: []string{
				"https://www.googleapis.com/auth/servicecontrol",
				"https://www.googleapis.com/auth/service.management.readonly",
				"https://www.googleapis.com/auth/logging.write",
				"https://www.googleapis.com/auth/monitoring.write",
				"https://www.googleapis.com/auth/trace.append",
				"https://www.googleapis.com/auth/devstorage.full_control",
			},
			MinCPUPlatform:     "Intel Skylake",
			BootDiskSize:       "10GB",
			BootDiskType:       "pd-ssd",
			BootDiskDeviceName: fmt.Sprintf("mixnet-%04d", (i + 1)),
			MaintenancePolicy:  "TERMINATE",
			Flags: []string{
				"--no-restart-on-failure",
			},
		}
	}

	for i := 0; i < 5; i++ {

		// Assign 5 configs to each zone.
		for j := 0; j < len(zones); j++ {
			configs[((i * len(zones)) + j)].Zone = zones[j]
		}
	}

	// Open new file in append-only mode.
	configFile, err := os.OpenFile(filepath.Join(configsPath, "gcloud-mixnet-20-40-30-10.cfg"), os.O_APPEND|os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed to create new configuration file: %v\n", err)
		os.Exit(1)
	}

	// Prepend with header.
	_, err = configFile.WriteString(configs[0].Header())
	if err != nil {
		fmt.Printf("Error while writing header to new configuration file: %v\n", err)
		os.Exit(1)
	}

	// Write created configuration slice to file.
	for i := 0; i < len(configs); i++ {

		_, err := configFile.WriteString(configs[i].Data())
		if err != nil {
			fmt.Printf("Error while writing data to new configuration file: %v\n", err)
			os.Exit(1)
		}
	}

	// Sync to stable storage.
	err = configFile.Sync()
	if err != nil {
		fmt.Printf("Error while syncing new configuration to stable storage: %v\n", err)
		os.Exit(1)
	}

	err = configFile.Close()
	if err != nil {
		fmt.Printf("Failed to close new configuration file: %v\n", err)
		os.Exit(1)
	}
}
