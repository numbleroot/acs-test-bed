package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
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

func spawnInstance(confChan <-chan Config, errChan chan<- error, proj string, serviceAcc string, bucket string, resultFolder string, pkiIP string) {

	for config := range confChan {

		// Prepare command with all corresponding arguments.
		cmd := exec.Command("/opt/google-cloud-sdk/bin/gcloud", "compute", fmt.Sprintf("--project=%s", proj), "instances", "create", config.Name,
			fmt.Sprintf("--service-account=%s", serviceAcc), fmt.Sprintf("--zone=%s", config.Zone),
			fmt.Sprintf("--machine-type=%s", config.MachineType), fmt.Sprintf("--min-cpu-platform=%s", config.MinCPUPlatform),
			fmt.Sprintf("--network-interface=%s", config.Network), fmt.Sprintf("--image=%s", config.Image),
			fmt.Sprintf("--image-project=%s", config.ImageProject), fmt.Sprintf("--boot-disk-size=%s", config.BootDiskSize),
			fmt.Sprintf("--boot-disk-type=%s", config.BootDiskType), fmt.Sprintf("--boot-disk-device-name=%s", config.BootDiskDeviceName),
			fmt.Sprintf("--metadata=typeOfNode=%s,resultFolder=%s,evalScriptToPull=%s,binaryToPull=%s,tcConfig=%s,pkiIP=%s,startup-script-url=gs://%s/startup.sh",
				config.TypeOfNode, resultFolder, config.EvaluationScript, config.BinaryName, config.ParamsTC, pkiIP, bucket),
			fmt.Sprintf("--scopes=%s", config.Scopes), fmt.Sprintf("--maintenance-policy=%s", config.MaintenancePolicy), config.Flags,
		)

		// Execute command and wait for completion.
		outRaw, err := cmd.CombinedOutput()
		if err != nil {
			errChan <- fmt.Errorf("spawning compute instance failed (code: '%v'):\n'%s'", err, outRaw)
			return
		}
		out := string(outRaw)

		// Verify successful machine creation.
		if strings.Contains(out, "RUNNING") {
			fmt.Printf("Successfully spawned instance %s\n", config.Name)
		} else {
			errChan <- fmt.Errorf("spawning compute instance returned failure message:\n'%s'", out)
			return
		}
	}

	errChan <- nil
}

func shutdownInstance(confChan <-chan Config, errChan chan<- error, proj string) {

	for config := range confChan {

		// Shut down compute instance.
		cmd := exec.Command("/opt/google-cloud-sdk/bin/gcloud", "compute", fmt.Sprintf("--project=%s", proj),
			"instances", "delete", config.Name, fmt.Sprintf("--zone=%s", config.Zone))

		// Execute command and wait for completion.
		out, err := cmd.CombinedOutput()
		if err != nil {
			errChan <- fmt.Errorf("deleting compute instance failed (code: '%v'):\n'%s'", err, out)
			return
		}

		// Verify successful instance deletion.
		if bytes.Contains(out, []byte("Deleted")) {
			fmt.Printf("Successfully deleted compute instance %s\n", config.Name)
		} else {
			errChan <- fmt.Errorf("deleting compute instance returned failure message:\n'%s'", out)
			return
		}
	}

	errChan <- nil
}

func main() {

	// Expect a number of command-line arguments.
	systemFlag := flag.String("system", "", "Specify which ACS to evaluate: 'zeno', 'vuvuzela', 'pung'.")
	configsPathFlag := flag.String("configsPath", "./gcloud-configs/", "Specify the file system location of the configurations folder for the compute instances.")
	resultsPathFlag := flag.String("resultsPath", "./results/", "Specify the file system location of the top-level results directory to create a new results folder under.")
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

	// Create name of results folder for this evaluation
	// run based on current time and system name.
	resultFolder := fmt.Sprintf("%s-%s", time.Now().Format("2006-01-02-15-04-05"), system)

	// Prepare local results folder.
	allResultsPath, err := filepath.Abs(*resultsPathFlag)
	if err != nil {
		fmt.Printf("Provided results path '%s' could not be converted to absolute path: %v\n", *resultsPathFlag, err)
		os.Exit(1)
	}

	// Prepare path to run-specific results folder.
	resultsPath := filepath.Join(allResultsPath, resultFolder)

	// If results folder does not exist yet, create it.
	// Also, add run-specific subfolder.
	err = os.MkdirAll(resultsPath, 0755)
	if err != nil {
		fmt.Printf("Failed to create results folder %s: %v\n", resultsPath, err)
		os.Exit(1)
	}

	var configsFile string

	if system == "zeno" {

		// Prepare zeno configurations file for ingestion.
		configsFileRel := filepath.Join(*configsPathFlag, "gcloud-mixnet-20-40-30-10_zeno.json")
		configsFile, err = filepath.Abs(configsFileRel)
		if err != nil {
			fmt.Printf("Unable to obtain absolute path to zeno configurations file '%s': %v\n", configsFileRel, err)
			os.Exit(1)
		}

	} else if system == "vuvuzela" {

		// Prepare vuvuzela configurations file for ingestion.
		configsFileRel := filepath.Join(*configsPathFlag, "gcloud-mixnet-20-40-30-10_vuvuzela.json")
		configsFile, err = filepath.Abs(configsFileRel)
		if err != nil {
			fmt.Printf("Unable to obtain absolute path to vuvuzela configurations file '%s': %v\n", configsFileRel, err)
			os.Exit(1)
		}

	} else if system == "pung" {

		// Prepare pung configurations file for ingestion.
		configsFileRel := filepath.Join(*configsPathFlag, "gcloud-mixnet-20-40-30-10_pung.json")
		configsFile, err = filepath.Abs(configsFileRel)
		if err != nil {
			fmt.Printf("Unable to obtain absolute path to pung configurations file '%s': %v\n", configsFileRel, err)
			os.Exit(1)
		}
	}

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

	var pkiConfig Config
	pkiInternalIP := ""
	pkiExternalIP := ""

	// If zeno: create PKI with decent hardware specs.
	if system == "zeno" {

		// Retrieve file defining PKI configuration.
		pkiConfigFileRel := filepath.Join(*configsPathFlag, "gcloud-mixnet-20-40-30-10_zeno-pki.json")
		pkiConfigFile, err := filepath.Abs(pkiConfigFileRel)
		if err != nil {
			fmt.Printf("Unable to obtain absolute path to PKI configuration file for zeno '%s': %v\n", pkiConfigFileRel, err)
			os.Exit(1)
		}

		// Ingest GCloud configuration file.
		pkiConfigJSON, err := ioutil.ReadFile(pkiConfigFile)
		if err != nil {
			fmt.Printf("Failed ingesting GCloud configuration file for zeno PKI: %v\n", err)
			os.Exit(1)
		}

		// Unmarshal JSON.
		err = json.Unmarshal(pkiConfigJSON, &pkiConfig)
		if err != nil {
			fmt.Printf("Error while trying to unmarshal JSON-encoded GCloud configuration for zeno PKI: %v\n", err)
			os.Exit(1)
		}

		// Prepare command to spawn zeno PKI instance.
		cmd := exec.Command("/opt/google-cloud-sdk/bin/gcloud", "compute", fmt.Sprintf("--project=%s", gcloudProject), "instances", "create", pkiConfig.Name,
			fmt.Sprintf("--service-account=%s", gcloudServiceAcc), fmt.Sprintf("--zone=%s", pkiConfig.Zone),
			fmt.Sprintf("--machine-type=%s", pkiConfig.MachineType), fmt.Sprintf("--min-cpu-platform=%s", pkiConfig.MinCPUPlatform),
			fmt.Sprintf("--network-interface=%s", pkiConfig.Network), fmt.Sprintf("--image=%s", pkiConfig.Image),
			fmt.Sprintf("--image-project=%s", pkiConfig.ImageProject), fmt.Sprintf("--boot-disk-size=%s", pkiConfig.BootDiskSize),
			fmt.Sprintf("--boot-disk-type=%s", pkiConfig.BootDiskType), fmt.Sprintf("--boot-disk-device-name=%s", pkiConfig.BootDiskDeviceName),
			fmt.Sprintf("--metadata=typeOfNode=%s,resultFolder=irrelevant,evalScriptToPull=%s,binaryToPull=%s,tcConfig=%s,pkiIP=irrelevant,startup-script-url=gs://%s/startup.sh",
				pkiConfig.TypeOfNode, pkiConfig.EvaluationScript, pkiConfig.BinaryName, pkiConfig.ParamsTC, gcloudBucket),
			fmt.Sprintf("--scopes=%s", pkiConfig.Scopes), fmt.Sprintf("--maintenance-policy=%s", pkiConfig.MaintenancePolicy), pkiConfig.Flags, "--tags=zeno-pki",
		)

		// Execute command and wait for completion.
		outRaw, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Spawning PKI for zeno mix-net failed (code: '%v'): '%s'", err, outRaw)
			os.Exit(1)
		}
		out := string(outRaw)

		// Compile regular expression matching on
		// the assigned IP addresses of the PKI.
		regexpIP, err := regexp.Compile("[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}")
		if err != nil {
			fmt.Printf("Failed to match internal IP address of zeno PKI instance: %v\n", err)
			os.Exit(1)
		}

		pkiIPs := regexpIP.FindAllString(out, -1)
		for i := range pkiIPs {

			if strings.HasPrefix(pkiIPs[i], "10.") {
				pkiInternalIP = pkiIPs[i]
			} else {
				pkiExternalIP = pkiIPs[i]
			}
		}

		// Verify successful machine creation.
		if strings.Contains(out, "RUNNING") {
			fmt.Printf("Successfully spawned PKI for zeno mix-net with public IP %s and internal IP %s\n", pkiExternalIP, pkiInternalIP)
		} else {
			fmt.Printf("Spawning PKI for zeno mix-net returned failure message:\n'%s'\n", out)
			os.Exit(1)
		}

		fmt.Printf("Waiting for zeno's PKI to finish initialization...")
		time.Sleep(75 * time.Second)
		fmt.Printf(" done.\n\n")
	}

	// Prepare channels to send configurations
	// to individual workers and expect responses.
	confChan := make(chan Config, len(configs))
	errChan := make(chan error)

	// Spawn 9 creation workers.
	for i := 0; i < 9; i++ {
		go spawnInstance(confChan, errChan, gcloudProject, gcloudServiceAcc, gcloudBucket, resultFolder, pkiInternalIP)
	}

	fmt.Printf("Spawning machines now...\n")

	// Iterate over configuration slice and spawn instances.
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

	fmt.Printf("All machines spawned!\n\n")

	// Wait for all instances to signal that
	// they have fetched all evaluation artifacts
	// from resources server.
	fmt.Printf("Waiting for machines to finish initialization...")
	time.Sleep(120 * time.Second)
	fmt.Printf(" done.\n\n")

	// If zeno: send PKI signal to start.
	if system == "zeno" {

		fmt.Printf("Signaling zeno's PKI to start epoch.\n")

		// Connect to control plane address used
		// only for evaluation purposes.
		ctrlConn, err := net.Dial("tcp", fmt.Sprintf("%s:26345", pkiExternalIP))
		if err != nil {
			fmt.Printf("Failed to connect to zeno PKI's control address to send start signal: %v\n", err)
			os.Exit(1)
		}

		// Send start signal.
		fmt.Fprintf(ctrlConn, "Please go ahead and start the epoch.\n")

		// Close connection once done.
		ctrlConn.Close()
	}

	// Spin-query GCloud bucket until we see
	// measurement files present from all nodes.

	// Download all files from GCloud bucket
	// to prepared local experiment folder.

	// Wait until enter key is pressed.
	fmt.Printf("\nPress ENTER to shutdown and delete all resources...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	if system == "zeno" {
		configs = append(configs, pkiConfig)
	}

	// Prepare channels to send configurations
	// to individual workers and expect responses.
	confChan = make(chan Config, len(configs))
	errChan = make(chan error)

	// Spawn 9 deletion workers.
	for i := 0; i < 9; i++ {
		go shutdownInstance(confChan, errChan, gcloudProject)
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
