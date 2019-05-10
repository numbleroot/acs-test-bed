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

// Config describes one compute instance
// exhaustively for reproducibility.
type Config struct {
	Name               string `json:"Name"`
	Zone               string `json:"Zone"`
	MachineType        string `json:"MachineType"`
	Network            string `json:"Network"`
	MinCPUPlatform     string `json:"MinCPUPlatform"`
	Scopes             string `json:"Scopes"`
	Image              string `json:"Image"`
	ImageProject       string `json:"ImageProject"`
	BootDiskSize       string `json:"BootDiskSize"`
	BootDiskType       string `json:"BootDiskType"`
	BootDiskDeviceName string `json:"BootDiskDeviceName"`
	MaintenancePolicy  string `json:"MaintenancePolicy"`
	Flags              string `json:"Flags"`
	TypeOfNode         string `json:"TypeOfNode"`
	EvaluationScript   string `json:"EvaluationScript"`
	BinaryName         string `json:"BinaryName"`
	ParamsTC           string `json:"ParamsTC"`
}

// GCloudZones contains all but one geographical
// zones GCP has to offer for compute nodes.
var GCloudZones = [19]string{
	"asia-east1-b",
	"asia-east2-b",
	"asia-northeast1-b",
	"asia-northeast2-b",
	"asia-south1-b",
	"asia-southeast1-b",
	"australia-southeast1-b",
	"europe-north1-b",
	"europe-west1-b",
	"europe-west2-b",
	"europe-west3-b",
	"europe-west4-b",
	"europe-west6-b",
	"northamerica-northeast1-b",
	"us-central1-b",
	"us-east1-b",
	"us-east4-b",
	"us-west1-b",
	"us-west2-b",
}

func shuffleZones() {

	// Shuffle zones slice.
	for i := (len(GCloudZones) - 1); i > 0; i-- {

		// Generate new CSPRNG number smaller than i.
		jBig, err := rand.Int(rand.Reader, big.NewInt(int64(i)))
		if err != nil {
			fmt.Printf("Drawing random number failed: %v\n", err)
			os.Exit(1)
		}
		j := int(jBig.Int64())

		// Swap places i and j in zones slice.
		GCloudZones[i], GCloudZones[j] = GCloudZones[j], GCloudZones[i]
	}
}

func main() {

	// Allow for control via command-line flags.
	configsPathFlag := flag.String("configsPath", "./gcloud-configs/", "Specify file system location where GCloud Compute configurations are supposed to be saved.")
	numClientsToGenFlag := flag.Int("numClientsToGen", 1000, "Specify the number of client nodes to generate according to the 20%%-40%%-30%%-10%% machine power classification. Should be a multiple of 100.")
	numVuvuzelaMixesToGenFlag := flag.Int("numVuvuzelaMixesToGen", 7, "Specify the number of vuvuzela mix nodes to generate (number of zeno mixes is twice this number minus 1).")
	flag.Parse()

	// Extract parsed flag values.
	configsPath, err := filepath.Abs(*configsPathFlag)
	if err != nil {
		fmt.Printf("Provided path to configuration files '%s' could not be converted to absolute path: %v\n", *configsPathFlag, err)
		os.Exit(1)
	}

	// Create configuration files folder
	// if it does not exist.
	err = os.MkdirAll(configsPath, 0755)
	if err != nil {
		fmt.Printf("Failed to create configurations folder %s: %v\n", configsPath, err)
		os.Exit(1)
	}

	numClientsToGen := *numClientsToGenFlag
	numVuvuzelaMixesToGen := *numVuvuzelaMixesToGenFlag

	numClientsFactor := numClientsToGen / 100.0
	if numClientsFactor < 1 {
		numClientsFactor = 1.0
	}

	// Prepare slices for respective client
	// compute node configurations.
	zenoConfigs := make([]Config, 0, (numClientsToGen + ((2 * numVuvuzelaMixesToGen) - 1)))
	vuvuzelaConfigs := make([]Config, 0, (numClientsToGen + numVuvuzelaMixesToGen))
	pungConfigs := make([]Config, 0, (numClientsToGen + 1))

	// Shuffle zones array.
	shuffleZones()
	zoneIdx := 0

	for i := 0; i < numClientsToGen; i++ {

		// Base machine specification on a distribution
		// approximating real-life hardware availability.
		machineType := "f1-micro"

		switch {
		case (i >= (numClientsFactor * 20)) && (i < (numClientsFactor * 60)):
			machineType = "f1-micro"

		case (i >= (numClientsFactor * 60)) && (i < (numClientsFactor * 90)):
			machineType = "f1-micro"

		case (i >= (numClientsFactor * 90)) && (i < (numClientsFactor * 100)):
			machineType = "f1-micro"
		}

		// Pick next zone from randomized zones array.
		zone := GCloudZones[zoneIdx]

		// Increment counter. If we traversed zones array
		// once, shuffle it again and reset index.
		zoneIdx++
		if zoneIdx == len(GCloudZones) {
			shuffleZones()
			zoneIdx = 0
		}

		// Prefill all configurations.
		zenoConfigs = append(zenoConfigs, Config{
			Name:               fmt.Sprintf("mixnet-%05d", (i + 1)),
			Zone:               zone,
			MachineType:        machineType,
			Network:            "subnet=default,network-tier=PREMIUM", //,no-address",
			MinCPUPlatform:     "Intel Skylake",
			Scopes:             "https://www.googleapis.com/auth/servicecontrol,https://www.googleapis.com/auth/service.management.readonly,https://www.googleapis.com/auth/logging.write,https://www.googleapis.com/auth/monitoring.write,https://www.googleapis.com/auth/trace.append,https://www.googleapis.com/auth/devstorage.full_control",
			Image:              "ubuntu-1804-bionic-v20190429",
			ImageProject:       "ubuntu-os-cloud",
			BootDiskSize:       "10GB",
			BootDiskType:       "pd-ssd",
			BootDiskDeviceName: "mixnet",
			MaintenancePolicy:  "TERMINATE",
			Flags:              "--no-restart-on-failure",
			TypeOfNode:         "client",
			EvaluationScript:   "zeno_client_eval.sh",
			BinaryName:         "zeno",
			ParamsTC:           "none so far",
		})

		vuvuzelaConfigs = append(vuvuzelaConfigs, Config{
			Name:               fmt.Sprintf("mixnet-%05d", (i + 1)),
			Zone:               zone,
			MachineType:        machineType,
			Network:            "subnet=default,network-tier=PREMIUM", //,no-address",
			MinCPUPlatform:     "Intel Skylake",
			Scopes:             "https://www.googleapis.com/auth/servicecontrol,https://www.googleapis.com/auth/service.management.readonly,https://www.googleapis.com/auth/logging.write,https://www.googleapis.com/auth/monitoring.write,https://www.googleapis.com/auth/trace.append,https://www.googleapis.com/auth/devstorage.full_control",
			Image:              "ubuntu-1804-bionic-v20190429",
			ImageProject:       "ubuntu-os-cloud",
			BootDiskSize:       "10GB",
			BootDiskType:       "pd-ssd",
			BootDiskDeviceName: "mixnet",
			MaintenancePolicy:  "TERMINATE",
			Flags:              "--no-restart-on-failure",
			TypeOfNode:         "client",
			EvaluationScript:   "vuvuzela-client_eval.sh",
			BinaryName:         "vuvuzela-client",
			ParamsTC:           "none so far",
		})

		pungConfigs = append(pungConfigs, Config{
			Name:               fmt.Sprintf("mixnet-%05d", (i + 1)),
			Zone:               zone,
			MachineType:        machineType,
			Network:            "subnet=default,network-tier=PREMIUM", //,no-address",
			MinCPUPlatform:     "Intel Skylake",
			Scopes:             "https://www.googleapis.com/auth/servicecontrol,https://www.googleapis.com/auth/service.management.readonly,https://www.googleapis.com/auth/logging.write,https://www.googleapis.com/auth/monitoring.write,https://www.googleapis.com/auth/trace.append,https://www.googleapis.com/auth/devstorage.full_control",
			Image:              "ubuntu-1804-bionic-v20190429",
			ImageProject:       "ubuntu-os-cloud",
			BootDiskSize:       "10GB",
			BootDiskType:       "pd-ssd",
			BootDiskDeviceName: "mixnet",
			MaintenancePolicy:  "TERMINATE",
			Flags:              "--no-restart-on-failure",
			TypeOfNode:         "client",
			EvaluationScript:   "pung_client_eval.sh",
			BinaryName:         "pung",
			ParamsTC:           "none so far",
		})
	}

	// Reset zones array and counter.
	shuffleZones()
	zoneIdx = 0

	// Also generate the specified number
	// of mix or server nodes.
	for i := numClientsToGen; i < (numClientsToGen + ((2 * numVuvuzelaMixesToGen) - 1)); i++ {

		// Pick next zone from randomized zones array.
		zone := GCloudZones[zoneIdx]

		// Increment counter. If we traversed zones array
		// once, shuffle it again and reset index.
		zoneIdx++
		if zoneIdx == len(GCloudZones) {
			shuffleZones()
			zoneIdx = 0
		}

		zenoConfigs = append(zenoConfigs, Config{
			Name:               fmt.Sprintf("mixnet-%05d", (i + 1)),
			Zone:               zone,
			MachineType:        "n1-standard-4",
			Network:            "subnet=default,network-tier=PREMIUM", //,no-address",
			MinCPUPlatform:     "Intel Skylake",
			Scopes:             "https://www.googleapis.com/auth/servicecontrol,https://www.googleapis.com/auth/service.management.readonly,https://www.googleapis.com/auth/logging.write,https://www.googleapis.com/auth/monitoring.write,https://www.googleapis.com/auth/trace.append,https://www.googleapis.com/auth/devstorage.full_control",
			Image:              "ubuntu-1804-bionic-v20190429",
			ImageProject:       "ubuntu-os-cloud",
			BootDiskSize:       "10GB",
			BootDiskType:       "pd-ssd",
			BootDiskDeviceName: "mixnet",
			MaintenancePolicy:  "TERMINATE",
			Flags:              "--no-restart-on-failure",
			TypeOfNode:         "mix",
			EvaluationScript:   "zeno_mix_eval.sh",
			BinaryName:         "zeno",
			ParamsTC:           "none so far",
		})
	}

	// Reset zones counter.
	zoneIdx = 0

	for i := numClientsToGen; i < (numClientsToGen + numVuvuzelaMixesToGen); i++ {

		zone := GCloudZones[zoneIdx]

		zoneIdx++
		if zoneIdx == len(GCloudZones) {
			shuffleZones()
			zoneIdx = 0
		}

		vuvuzelaConfigs = append(vuvuzelaConfigs, Config{
			Name:               fmt.Sprintf("mixnet-%05d", (i + 1)),
			Zone:               zone,
			MachineType:        "n1-standard-4",
			Network:            "subnet=default,network-tier=PREMIUM", //,no-address",
			MinCPUPlatform:     "Intel Skylake",
			Scopes:             "https://www.googleapis.com/auth/servicecontrol,https://www.googleapis.com/auth/service.management.readonly,https://www.googleapis.com/auth/logging.write,https://www.googleapis.com/auth/monitoring.write,https://www.googleapis.com/auth/trace.append,https://www.googleapis.com/auth/devstorage.full_control",
			Image:              "ubuntu-1804-bionic-v20190429",
			ImageProject:       "ubuntu-os-cloud",
			BootDiskSize:       "10GB",
			BootDiskType:       "pd-ssd",
			BootDiskDeviceName: "mixnet",
			MaintenancePolicy:  "TERMINATE",
			Flags:              "--no-restart-on-failure",
			TypeOfNode:         "vuvuzela-mixer",
			EvaluationScript:   "vuvuzela-mixer_eval.sh",
			BinaryName:         "vuvuzela-mixer",
			ParamsTC:           "none so far",
		})

		if i == numClientsToGen {
			vuvuzelaConfigs[i].TypeOfNode = "vuvuzela-coordinator"
			vuvuzelaConfigs[i].EvaluationScript = "vuvuzela-coordinator_eval.sh"
			vuvuzelaConfigs[i].BinaryName = "vuvuzela-coordinator_eval.sh"
		}
	}

	pungConfigs = append(pungConfigs, Config{
		Name:               fmt.Sprintf("mixnet-%05d", (numClientsToGen + 1)),
		Zone:               GCloudZones[0],
		MachineType:        "n1-standard-4",
		Network:            "subnet=default,network-tier=PREMIUM", //,no-address",
		MinCPUPlatform:     "Intel Skylake",
		Scopes:             "https://www.googleapis.com/auth/servicecontrol,https://www.googleapis.com/auth/service.management.readonly,https://www.googleapis.com/auth/logging.write,https://www.googleapis.com/auth/monitoring.write,https://www.googleapis.com/auth/trace.append,https://www.googleapis.com/auth/devstorage.full_control",
		Image:              "ubuntu-1804-bionic-v20190429",
		ImageProject:       "ubuntu-os-cloud",
		BootDiskSize:       "10GB",
		BootDiskType:       "pd-ssd",
		BootDiskDeviceName: "mixnet",
		MaintenancePolicy:  "TERMINATE",
		Flags:              "--no-restart-on-failure",
		TypeOfNode:         "server",
		EvaluationScript:   "pung_server_eval.sh",
		BinaryName:         "pung",
		ParamsTC:           "none so far",
	})

	// Marshal slice of zeno configs to JSON.
	zenoConfigsJSON, err := json.MarshalIndent(zenoConfigs, "", "\t")
	if err != nil {
		fmt.Printf("Failed to marshal zeno configurations to JSON: %v\n", err)
		os.Exit(1)
	}

	// Write zeno configs to file.
	err = ioutil.WriteFile(filepath.Join(configsPath, "gcloud-mixnet-20-40-30-10_zeno.json"), zenoConfigsJSON, 0644)
	if err != nil {
		fmt.Printf("Error writing zeno configurations in JSON format to file: %v\n", err)
		os.Exit(1)
	}

	// Marshal slice of Vuvuzela configs to JSON.
	vuvuzelaConfigsJSON, err := json.MarshalIndent(vuvuzelaConfigs, "", "\t")
	if err != nil {
		fmt.Printf("Failed to marshal vuvuzela configurations to JSON: %v\n", err)
		os.Exit(1)
	}

	// Write vuvuzela configs to file.
	err = ioutil.WriteFile(filepath.Join(configsPath, "gcloud-mixnet-20-40-30-10_vuvuzela.json"), vuvuzelaConfigsJSON, 0644)
	if err != nil {
		fmt.Printf("Error writing vuvuzela configurations in JSON format to file: %v\n", err)
		os.Exit(1)
	}

	// Marshal slice of pung configs to JSON.
	pungConfigsJSON, err := json.MarshalIndent(pungConfigs, "", "\t")
	if err != nil {
		fmt.Printf("Failed to marshal pung configurations to JSON: %v\n", err)
		os.Exit(1)
	}

	// Write pung configs to file.
	err = ioutil.WriteFile(filepath.Join(configsPath, "gcloud-mixnet-20-40-30-10_pung.json"), pungConfigsJSON, 0644)
	if err != nil {
		fmt.Printf("Error writing pung configurations in JSON format to file: %v\n", err)
		os.Exit(1)
	}

	// Additionally, create configuration for zeno's PKI node.
	zenoPKIConfigsJSON, err := json.MarshalIndent(Config{
		Name:               "zeno-pki",
		Zone:               "europe-west3-a",
		MachineType:        "n1-standard-4",
		Network:            "subnet=default,network-tier=PREMIUM",
		MinCPUPlatform:     "Intel Skylake",
		Scopes:             "https://www.googleapis.com/auth/servicecontrol,https://www.googleapis.com/auth/service.management.readonly,https://www.googleapis.com/auth/logging.write,https://www.googleapis.com/auth/monitoring.write,https://www.googleapis.com/auth/trace.append,https://www.googleapis.com/auth/devstorage.full_control",
		Image:              "ubuntu-1804-bionic-v20190429",
		ImageProject:       "ubuntu-os-cloud",
		BootDiskSize:       "10GB",
		BootDiskType:       "pd-ssd",
		BootDiskDeviceName: "mixnet",
		MaintenancePolicy:  "TERMINATE",
		Flags:              "--no-restart-on-failure",
		TypeOfNode:         "zeno-pki",
		EvaluationScript:   "zeno-pki_eval.sh",
		BinaryName:         "zeno-pki",
		ParamsTC:           "irrelevant",
	}, "", "\t")
	if err != nil {
		fmt.Printf("Failed to marshal configuration for zeno's PKI to JSON: %v\n", err)
		os.Exit(1)
	}

	// Write PKI configuration for zeno to file.
	err = ioutil.WriteFile(filepath.Join(configsPath, "gcloud-mixnet-20-40-30-10_zeno-pki.json"), zenoPKIConfigsJSON, 0644)
	if err != nil {
		fmt.Printf("Error writing configuration for zeno's PKI in JSON format to file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("All done!\n")
}
