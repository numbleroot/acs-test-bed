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
	Name               string   `json:"Name"`
	Zone               string   `json:"Zone"`
	MachineType        string   `json:"MachineType"`
	Subnet             string   `json:"Subnet"`
	NetworkTier        string   `json:"NetworkTier"`
	MinCPUPlatform     string   `json:"MinCPUPlatform"`
	Scopes             []string `json:"Scopes"`
	Image              string   `json:"Image"`
	ImageProject       string   `json:"ImageProject"`
	BootDiskSize       string   `json:"BootDiskSize"`
	BootDiskType       string   `json:"BootDiskType"`
	BootDiskDeviceName string   `json:"BootDiskDeviceName"`
	MaintenancePolicy  string   `json:"MaintenancePolicy"`
	Flags              []string `json:"Flags"`
	TypeOfNode         string
	EvaluationScript   string
	BinaryName         string
	ParamsTC           string
}

func spawnInstance(confChan <-chan Config, errChan chan<- error, proj string, serviceAcc string, bucket string) {

	for config := range confChan {

		scopes := strings.Join(config.Scopes, ",")
		flags := strings.Join(config.Flags, " ")

		// Prepare command to create a compute
		// instance configured according to Config.
		cmd := exec.Command("/opt/google-cloud-sdk/bin/gcloud", "compute", fmt.Sprintf("--project=%s", proj), "instances", "create", config.Name,
			fmt.Sprintf("--service-account=%s", serviceAcc), fmt.Sprintf("--zone=%s", config.Zone),
			fmt.Sprintf("--machine-type=%s", config.MachineType), fmt.Sprintf("--min-cpu-platform=%s", config.MinCPUPlatform),
			fmt.Sprintf("--subnet=%s", config.Subnet), fmt.Sprintf("--network-tier=%s", config.NetworkTier),
			fmt.Sprintf("--image=%s", config.Image), fmt.Sprintf("--image-project=%s", config.ImageProject), fmt.Sprintf("--boot-disk-size=%s", config.BootDiskSize),
			fmt.Sprintf("--boot-disk-type=%s", config.BootDiskType), fmt.Sprintf("--boot-disk-device-name=%s", config.BootDiskDeviceName),
			fmt.Sprintf("--metadata=typeOfNode=%s,evalScriptToPull=%s,binaryToPull=%s,tcConfig=%s,startup-script-url=gs://%s/startup.sh", config.TypeOfNode, config.EvaluationScript, config.BinaryName, config.ParamsTC, bucket),
			fmt.Sprintf("--scopes=%s", scopes), fmt.Sprintf("--maintenance-policy=%s", config.MaintenancePolicy), flags)

		// Execute command and wait for completion.
		out, err := cmd.CombinedOutput()
		if err != nil {
			errChan <- fmt.Errorf("spawning compute instance failed (code: '%v'): '%s'", err, out)
			return
		}

		// Verify successful machine creation.
		if bytes.Contains(out, []byte("Created")) && bytes.Contains(out, []byte("RUNNING")) {
			fmt.Printf("Successfully spawned instance %s\n", config.Name)
		} else {
			errChan <- fmt.Errorf("spawning compute instance returned failure message: '%s'", out)
			return
		}
	}

	errChan <- nil
}

func shutdownInstance(confChan <-chan Config, errChan chan<- error, proj string, serviceAcc string, bucket string) {

	for config := range confChan {

		// Shut down compute instance.
		cmd := exec.Command("/opt/google-cloud-sdk/bin/gcloud", "compute", fmt.Sprintf("--project=%s", proj),
			"instances", "delete", config.Name, fmt.Sprintf("--zone=%s", config.Zone))

		// Execute command and wait for completion.
		out, err := cmd.CombinedOutput()
		if err != nil {
			errChan <- fmt.Errorf("deleting compute instance failed (code: '%v'): '%s'", err, out)
			return
		}

		// Verify successful instance deletion.
		if bytes.Contains(out, []byte("Deleted")) {
			fmt.Printf("Successfully deleted compute instance %s\n", config.Name)
		} else {
			errChan <- fmt.Errorf("deleting compute instance returned failure message: '%s'", out)
			return
		}
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

	system := strings.ToLower(*systemFlag)
	gcloudProject := *gcloudProjectFlag
	gcloudServiceAcc := *gcloudServiceAccFlag
	gcloudBucket := *gcloudBucketFlag

	// System flag has to be one of three values.
	if system != "zeno" && system != "vuvuzela" && system != "pung" {
		fmt.Printf("Flag '-system' requires one of the three values: 'zeno', 'vuvuzela', or 'pung'.")
		os.Exit(1)
	}

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
	_, err = os.Stat(resultsPath)
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

	// Create a new results subfolder within the
	// supplied GCloud Storage bucket.

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

	configs = configs[:1]

	// Prepare channels to send configurations
	// to individual workers and expect responses.
	confChan := make(chan Config, len(configs))
	errChan := make(chan error)

	// Spawn 1 creation workers.
	for i := 0; i < 1; i++ {
		go spawnInstance(confChan, errChan, gcloudProject, gcloudServiceAcc, gcloudBucket)
	}

	fmt.Printf("Spawning machines now...\n")

	// Iterate over configuration slice and spawn instances.
	for _, config := range configs {

		config.TypeOfNode = "client"
		config.EvaluationScript = "zeno-pki_eval.sh"
		config.BinaryName = "zeno-pki"
		config.ParamsTC = "none so far"

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

	fmt.Printf("All machines spawned!\n\n")

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

	// Spawn 1 deletion workers.
	for i := 0; i < 1; i++ {
		go shutdownInstance(confChan, errChan, gcloudProject, gcloudServiceAcc, gcloudBucket)
	}

	fmt.Printf("Deleting boot disks and machines now...\n")

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

	fmt.Printf("All disks and machines deleted!\n\n")
}
