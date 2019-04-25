package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func spawnInstance(confChan <-chan Config, errChan chan<- error, proj string, serviceAcc string, bucket string) {

	for config := range confChan {

		// Prepare valid boot disk name.
		diskName := strings.Replace(config.Name, ".", "-", -1)

		// Prepare command to create boot disk from
		// snapshot we created.
		cmdDisk := exec.Command("/opt/google-cloud-sdk/bin/gcloud", "compute", fmt.Sprintf("--project=%s", proj), "disks", "create", diskName,
			fmt.Sprintf("--size=%s", config.BootDiskSize), fmt.Sprintf("--zone=%s", config.Zone),
			fmt.Sprintf("--source-snapshot=%s", config.SourceSnapshot), fmt.Sprintf("--type=%s", config.BootDiskType))

		// Execute command and wait for completion.
		out, err := cmdDisk.CombinedOutput()
		if err != nil {
			errChan <- fmt.Errorf("creating boot disk failed (code: '%v'): '%s'", err, out)
			return
		}

		// Verify successful boot disk creation.
		if bytes.Contains(out, []byte("Created")) && bytes.Contains(out, []byte("READY")) {
			fmt.Printf("Successfully created disk %s for instance %s\n", diskName, config.Name)
		} else {
			errChan <- fmt.Errorf("creating boot disk returned failure message: '%s'", out)
			return
		}

		// Prepare command to create a compute
		// instance configured according to Config.

		// Execute command and wait for completion.
	}

	errChan <- nil
}

func shutdownInstance(confChan <-chan Config, errChan chan<- error, proj string, serviceAcc string, bucket string) {

	for config := range confChan {

		// Prepare valid boot disk name.
		diskName := strings.Replace(config.Name, ".", "-", -1)

		// Prepare command to delete boot disk that
		// we spawned earlier.
		cmdDisk := exec.Command("/opt/google-cloud-sdk/bin/gcloud", "compute", fmt.Sprintf("--project=%s", proj),
			"disks", "delete", diskName, fmt.Sprintf("--zone=%s", config.Zone))

		// Execute command and wait for completion.
		out, err := cmdDisk.CombinedOutput()
		if err != nil {
			errChan <- fmt.Errorf("deleting boot disk failed (code: '%v'): '%s'", err, out)
			return
		}

		// Verify successful boot disk deletion.
		if bytes.Contains(out, []byte("Deleted")) {
			fmt.Printf("Successfully deleted disk %s for instance %s\n", diskName, config.Name)
		} else {
			errChan <- fmt.Errorf("deleting boot disk returned failure message: '%s'", out)
			return
		}

		// Shut down compute instance.
	}

	errChan <- nil
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

	configs = configs[:5]

	// Prepare channels to send configurations
	// to individual workers and expect responses.
	confChan := make(chan Config, len(configs))
	errChan := make(chan error)

	// Spawn 5 creation workers.
	for i := 0; i < 5; i++ {
		go spawnInstance(confChan, errChan, gcloudProject, gcloudServiceAcc, gcloudBucket)
	}

	// Iterate over configuration slice, create
	// boot disks and spawn instances.
	for _, config := range configs {
		confChan <- config
	}
	close(confChan)

	for range configs {

		// If any worker threw an error, abort.
		err := <-errChan
		if err != nil {
			fmt.Printf("Subroutine creating boot disk and spawning a machine failed: %v\n", err)
			os.Exit(1)
		}
	}
	close(errChan)

	// Wait for all instances to signal that
	// they have fetched all evaluation artifacts
	// from resources server.

	// If zeno: send PKI signal to start.

	// Spin-query GCloud bucket until we see
	// measurement files present from all nodes.

	// Download all files from GCloud bucket
	// to prepared local experiment folder.

	// Wait until enter key is pressed.
	fmt.Printf("\nPress ENTER to shutdown and delete all resources...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	// Prepare channels to send configurations
	// to individual workers and expect responses.
	confChan = make(chan Config, len(configs))
	errChan = make(chan error)

	// Spawn 5 deletion workers.
	for i := 0; i < 5; i++ {
		go shutdownInstance(confChan, errChan, gcloudProject, gcloudServiceAcc, gcloudBucket)
	}

	// Shutdown and destroy disks and instances.
	for _, config := range configs {
		confChan <- config
	}
	close(confChan)

	for range configs {

		// If any worker threw an error, abort.
		err := <-errChan
		if err != nil {
			fmt.Printf("Subroutine deleting boot disk and shutting down a machine failed: %v\n", err)
			os.Exit(1)
		}
	}
	close(errChan)
}
