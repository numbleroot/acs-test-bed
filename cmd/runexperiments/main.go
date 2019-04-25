package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

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

	// Expect a number of command-line arguments.
	systemFlag := flag.String("system", "", "Specify which ACS to evaluate: 'zeno', 'vuvuzela', 'pung'.")
	configsFileFlag := flag.String("configsFile", "./gcloud-configs/gcloud-mixnet-20-40-30-10.json", "Specify the file system location of the configuration file for the compute instances.")
	resultsPathFlag := flag.String("resultsFile", "./results/", "Specify the file system location of the top-level results directory to create a new results folder under.")
	gcloudProjectFlag := flag.String("gcloudProj", "", "Supply the GCloud project identifier.")
	gcloudServiceAccFlag := flag.String("gcloudServiceAcc", "", "Supply the GCloud Service Account identifier.")
	gcloudBucketFlag := flag.String("gcloudBucket", "", "Supply the GCloud Storage Bucket to use for the experiments.")
	flag.Parse()

	// Enforce arguments to be set.
	if *systemFlag == "" || *gcloudProjectFlag == "" || *gcloudServiceAccFlag == "" || *gcloudBucketFlag == "" {
		fmt.Printf("Missing arguments, please provide values for all flags: '-system', '-gcloudProj', '-gcloudServiceAcc', and '-gcloudBucket'.\n")
		os.Exit(1)
	}

	system := *systemFlag
	gcloudProject := *gcloudProjectFlag
	gcloudServiceAcc := *gcloudServiceAccFlag
	gcloudBucket := *gcloudBucketFlag

	// Prepare configurations file for ingestion.
	configsFile, err := filepath.Abs(*configsFileFlag)
	if err != nil {
		fmt.Printf("Provided configs file '%s' could not be converted to absolute path: %v\n", *configsFileFlag, err)
		os.Exit(1)
	}

	// Prepare local results folder.
	allResultsPath, err := filepath.Abs(*resultsPathFlag)
	if err != nil {
		fmt.Printf("Provided results path '%s' could not be converted to absolute path: %v\n", *resultsPathFlag, err)
		os.Exit(1)
	}

	// Prepare path to run-specific results folder.
	resultsPath := filepath.Join(allResultsPath, fmt.Sprintf("%s-%s", time.Now().Format("2006-01-02-15-04-05"), system))

	// If results folder does not exist yet, create it.
	// Also, add run-specific subfolder.
	_, err = os.Stat(allResultsPath)
	if os.IsNotExist(err) {

		err := os.MkdirAll(resultsPath, 0755)
		if err != nil {
			fmt.Printf("Failed to create results folder %s: %v\n", resultsPath, err)
			os.Exit(1)
		}
	}

	// If zeno: create PKI with PREMIUM network
	// tier and decent hardware specs.

	// Verify it is running.

	// Ingest GCloud configuration file.
	configsJSON, err := ioutil.ReadFile(configsFile)
	if err != nil {
		fmt.Printf("Failed ingesting GCloud configuration file: %v\n", err)
		os.Exit(1)
	}

	// Unmarshal JSON.
	var configs []Config
	err = json.Unmarshal(configsJSON, &configs)
	if err != nil {
		fmt.Printf("Error while trying to unmarshal JSON-encoded GCloud configurations: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("gcloudProject: '%s', gcloudServiceAcc: '%s', gcloudBucket: '%s'\n\n", gcloudProject, gcloudServiceAcc, gcloudBucket)

	// Iterate over configuration slice and execute
	// two gcloud util commands for each line.
	for _, config := range configs {
		fmt.Printf("\n--- CONFIG: ---\n%v\n", config)
	}

	// First: create boot drives from our snapshot
	// for each instance.

	// Second: start compute instances.

	// Verify all instances are running.

	// Wait for all instances to signal that
	// they have fetched all evaluation artifacts
	// from resources server.

	// If zeno: send PKI signal to start.

	// Spin-query GCloud bucket until we see
	// measurement files present from all nodes.

	// Download all files from GCloud bucket
	// to prepared local experiment folder.
}
