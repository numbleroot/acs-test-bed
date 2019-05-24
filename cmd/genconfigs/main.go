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
	Name             string `json:"Name"`
	Partner          string `json:"Partner"`
	Zone             string `json:"Zone"`
	MinCPUPlatform   string `json:"MinCPUPlatform"`
	MachineType      string `json:"MachineType"`
	TypeOfNode       string `json:"TypeOfNode"`
	EvaluationScript string `json:"EvaluationScript"`
	BinaryName       string `json:"BinaryName"`
	ParamsTC         string `json:"ParamsTC"`
	SourceImage      string `json:"SourceImage"`
	DiskType         string `json:"DiskType"`
	DiskSize         string `json:"DiskSize"`
}

// GCloudZones holds all but two GCloud zones.
var GCloudZones = [18]string{
	"asia-east1-b",
	"asia-northeast1-b",
	"asia-southeast1-b",
	"europe-west1-b",
	"europe-west4-b",
	"us-central1-b",
	"us-east1-b",
	"us-east4-b",
	"us-west1-b",
	// Low storage zones below
	"asia-east2-b",
	"asia-northeast2-b",
	"asia-south1-b",
	"australia-southeast1-b",
	"europe-north1-b",
	"europe-west2-b",
	"europe-west6-b",
	"northamerica-northeast1-b",
	"us-west2-b",
}

// GCloudZonesLowStorage defines limits in terms
// of number of instances with a standard disk
// size of 10GB that can be spawned in a GCloud
// zone without risking overstepping the quota.
// This might change if a GCP project is assigned
// more lenient quotas.
var GCloudZonesLowStorage = map[string]int{
	"asia-east2-b":              45,
	"asia-northeast2-b":         45,
	"asia-south1-b":             45,
	"australia-southeast1-b":    45,
	"europe-north1-b":           45,
	"europe-west2-b":            45,
	"europe-west6-b":            45,
	"northamerica-northeast1-b": 45,
	"us-west2-b":                45,
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
	numClientsToGenFlag := flag.Int("numClientsToGen", 1000, "Specify the number of client nodes to generate. Has to be an even number.")
	numVuvuzelaMixesToGenFlag := flag.Int("numVuvuzelaMixesToGen", 7, "Specify the number of vuvuzela mix nodes to generate (number of zeno mixes is twice this number minus 1).")
	numZenoCascadesFlag := flag.Int("numZenoCascades", 1, "Specify the number of cascades in zeno to generate.")
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

	if ((*numClientsToGenFlag % 2) != 0) || (*numClientsToGenFlag < 2) {
		fmt.Printf("Number of clients to generate has to be an even number > 1.\n")
		os.Exit(1)
	}

	numClientsToGen := *numClientsToGenFlag
	numVuvuzelaMixesToGen := *numVuvuzelaMixesToGenFlag
	numZenoCascades := *numZenoCascadesFlag

	// Create local map that tracks how many
	// instances have already been assigned
	// to each of the low storage zones.
	gcloudZonesLowStorage := map[string]int{
		"asia-east2-b":              0,
		"asia-northeast2-b":         0,
		"asia-south1-b":             0,
		"australia-southeast1-b":    0,
		"europe-north1-b":           0,
		"europe-west2-b":            0,
		"europe-west6-b":            0,
		"northamerica-northeast1-b": 0,
		"us-west2-b":                0,
	}

	// Prepare slices for respective client
	// compute node configurations.
	zenoConfigs := make([]Config, 0, (numClientsToGen + (numZenoCascades * ((2 * numVuvuzelaMixesToGen) - 1))))
	vuvuzelaConfigs := make([]Config, 0, (numClientsToGen + numVuvuzelaMixesToGen))
	pungConfigs := make([]Config, 0, numClientsToGen)

	// Shuffle zones array.
	shuffleZones()
	zoneIdx := 0

	for i := 0; i < numClientsToGen; i++ {

		name := fmt.Sprintf("client-%05d", (i + 1))

		// Determine partner client of this client.
		var partner string
		if (i % 2) == 0 {
			partner = fmt.Sprintf("client-%05d", (i + 2))
		} else {
			partner = fmt.Sprintf("client-%05d", i)
		}

		// Pick next zone from randomized zones array.
		zone := GCloudZones[zoneIdx]

		// In case the storage quota for this zone has
		// been reached, pick the next zone.
		for gcloudZonesLowStorage[zone] > GCloudZonesLowStorage[zone] {

			// Increment counter. If we traversed zones array
			// once, shuffle it again and reset index.
			zoneIdx++
			if zoneIdx == len(GCloudZones) {
				shuffleZones()
				zoneIdx = 0
			}

			zone = GCloudZones[zoneIdx]
		}

		// Increment local map counter if this is a zone
		// constrained with persistent disk space.
		if GCloudZonesLowStorage[zone] == 45 {
			gcloudZonesLowStorage[zone]++
		}

		zoneIdx++
		if zoneIdx == len(GCloudZones) {
			shuffleZones()
			zoneIdx = 0
		}

		// Specify machine type.
		machineType := "f1-micro"

		// Prefill all configurations.
		zenoConfigs = append(zenoConfigs, Config{
			Name:             name,
			Partner:          partner,
			Zone:             zone,
			MinCPUPlatform:   "Intel Skylake",
			MachineType:      machineType,
			TypeOfNode:       "client",
			EvaluationScript: "zeno_client_eval.sh",
			BinaryName:       "zeno",
			ParamsTC:         "none so far",
			SourceImage:      "acs",
			DiskType:         "pd-ssd",
			DiskSize:         "10",
		})

		vuvuzelaConfigs = append(vuvuzelaConfigs, Config{
			Name:             name,
			Partner:          partner,
			Zone:             zone,
			MinCPUPlatform:   "Intel Skylake",
			MachineType:      machineType,
			TypeOfNode:       "client",
			EvaluationScript: "vuvuzela-client_eval.sh",
			BinaryName:       "vuvuzela-client",
			ParamsTC:         "none so far",
			SourceImage:      "acs",
			DiskType:         "pd-ssd",
			DiskSize:         "10",
		})

		pungConfigs = append(pungConfigs, Config{
			Name:             name,
			Partner:          partner,
			Zone:             zone,
			MinCPUPlatform:   "Intel Skylake",
			MachineType:      machineType,
			TypeOfNode:       "client",
			EvaluationScript: "pung_client_eval.sh",
			BinaryName:       "pung-client",
			ParamsTC:         "none so far",
			SourceImage:      "acs",
			DiskType:         "pd-ssd",
			DiskSize:         "10",
		})
	}

	// Reset zones array and counter.
	shuffleZones()
	zoneIdx = 0

	// Also generate the specified number
	// of mix or server nodes.
	for i := numClientsToGen; i < (numClientsToGen + (numZenoCascades * ((2 * numVuvuzelaMixesToGen) - 1))); i++ {

		// Pick next zone from randomized zones array.
		zone := GCloudZones[zoneIdx]

		// In case the storage quota for this zone has
		// been reached, pick the next zone.
		for gcloudZonesLowStorage[zone] > GCloudZonesLowStorage[zone] {

			// Increment counter. If we traversed zones array
			// once, shuffle it again and reset index.
			zoneIdx++
			if zoneIdx == len(GCloudZones) {
				shuffleZones()
				zoneIdx = 0
			}

			zone = GCloudZones[zoneIdx]
		}

		// Increment local map counter if this is a zone
		// constrained with persistent disk space.
		if GCloudZonesLowStorage[zone] == 45 {
			gcloudZonesLowStorage[zone]++
		}

		zoneIdx++
		if zoneIdx == len(GCloudZones) {
			shuffleZones()
			zoneIdx = 0
		}

		zenoConfigs = append(zenoConfigs, Config{
			Name:             fmt.Sprintf("mixnet-%05d", (i + 1)),
			Partner:          "irrelevant",
			Zone:             zone,
			MinCPUPlatform:   "Intel Skylake",
			MachineType:      "n1-standard-4",
			TypeOfNode:       "mix",
			EvaluationScript: "zeno_mix_eval.sh",
			BinaryName:       "zeno",
			ParamsTC:         "none so far",
			SourceImage:      "acs",
			DiskType:         "pd-ssd",
			DiskSize:         "10",
		})
	}

	zoneIdx = 0

	for i := numClientsToGen; i < (numClientsToGen + numVuvuzelaMixesToGen); i++ {

		zone := GCloudZones[zoneIdx]

		for gcloudZonesLowStorage[zone] > GCloudZonesLowStorage[zone] {

			zoneIdx++
			if zoneIdx == len(GCloudZones) {
				shuffleZones()
				zoneIdx = 0
			}

			zone = GCloudZones[zoneIdx]
		}

		if GCloudZonesLowStorage[zone] == 45 {
			gcloudZonesLowStorage[zone]++
		}

		zoneIdx++
		if zoneIdx == len(GCloudZones) {
			shuffleZones()
			zoneIdx = 0
		}

		vuvuzelaConfigs = append(vuvuzelaConfigs, Config{
			Name:             fmt.Sprintf("mixnet-%05d", (i + 1)),
			Partner:          "irrelevant",
			Zone:             zone,
			MinCPUPlatform:   "Intel Skylake",
			MachineType:      "n1-standard-4",
			TypeOfNode:       "vuvuzela-mixer",
			EvaluationScript: "vuvuzela-mixer_eval.sh",
			BinaryName:       "vuvuzela-mixer",
			ParamsTC:         "none so far",
			SourceImage:      "acs",
			DiskType:         "pd-ssd",
			DiskSize:         "10",
		})

		if i == numClientsToGen {
			vuvuzelaConfigs[i].TypeOfNode = "vuvuzela-coordinator"
			vuvuzelaConfigs[i].EvaluationScript = "vuvuzela-coordinator_eval.sh"
			vuvuzelaConfigs[i].BinaryName = "vuvuzela-coordinator_eval.sh"
		}
	}

	// Marshal slice of zeno configs to JSON.
	zenoConfigsJSON, err := json.MarshalIndent(zenoConfigs, "", "\t")
	if err != nil {
		fmt.Printf("Failed to marshal zeno configurations to JSON: %v\n", err)
		os.Exit(1)
	}

	// Write zeno configs to file.
	err = ioutil.WriteFile(filepath.Join(configsPath, "gcloud-zeno.json"), zenoConfigsJSON, 0644)
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
	err = ioutil.WriteFile(filepath.Join(configsPath, "gcloud-vuvuzela.json"), vuvuzelaConfigsJSON, 0644)
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
	err = ioutil.WriteFile(filepath.Join(configsPath, "gcloud-pung.json"), pungConfigsJSON, 0644)
	if err != nil {
		fmt.Printf("Error writing pung configurations in JSON format to file: %v\n", err)
		os.Exit(1)
	}

	// Additionally, create configuration for zeno's PKI node.
	zenoPKIConfigsJSON, err := json.MarshalIndent(Config{
		Name:             "zeno-pki",
		Partner:          "irrelevant",
		Zone:             "europe-west3-b",
		MinCPUPlatform:   "Intel Skylake",
		MachineType:      "n1-standard-8",
		TypeOfNode:       "zeno-pki",
		EvaluationScript: "zeno-pki_eval.sh",
		BinaryName:       "zeno-pki",
		ParamsTC:         "irrelevant",
		SourceImage:      "acs",
		DiskType:         "pd-ssd",
		DiskSize:         "10",
	}, "", "\t")
	if err != nil {
		fmt.Printf("Failed to marshal configuration for zeno's PKI to JSON: %v\n", err)
		os.Exit(1)
	}

	// Write PKI configuration for zeno to file.
	err = ioutil.WriteFile(filepath.Join(configsPath, "gcloud-zeno-pki.json"), zenoPKIConfigsJSON, 0644)
	if err != nil {
		fmt.Printf("Error writing configuration for zeno's PKI in JSON format to file: %v\n", err)
		os.Exit(1)
	}

	// Prepare server configuration for pung.
	pungServerJSON, err := json.MarshalIndent(Config{
		Name:             fmt.Sprintf("mixnet-%05d", (numClientsToGen + 1)),
		Partner:          "irrelevant",
		Zone:             GCloudZones[0],
		MinCPUPlatform:   "Intel Skylake",
		MachineType:      "n1-highmem-16",
		TypeOfNode:       "server",
		EvaluationScript: "pung_server_eval.sh",
		BinaryName:       "pung-server",
		ParamsTC:         "none so far",
		SourceImage:      "acs",
		DiskType:         "pd-ssd",
		DiskSize:         "10",
	}, "", "\t")
	if err != nil {
		fmt.Printf("Failed to marshal configuration for pung's server to JSON: %v\n", err)
		os.Exit(1)
	}

	// Write pung's server configuration to file.
	err = ioutil.WriteFile(filepath.Join(configsPath, "gcloud-pung-server.json"), pungServerJSON, 0644)
	if err != nil {
		fmt.Printf("Error writing configuration for pung's server in JSON format to file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("All done!\n")
}
