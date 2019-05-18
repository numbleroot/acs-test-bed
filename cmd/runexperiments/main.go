package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
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
	Name             string `json:"Name"`
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

var tmplInstanceCreate string

var tmplInstancePublicIP = `
      "accessConfigs": [
        {
          "kind": "compute#accessConfig",
          "name": "External NAT",
          "type": "ONE_TO_ONE_NAT",
          "networkTier": "PREMIUM"
        }
      ],`

func spawnInstance(config *Config, proj string, serviceAcc string, accessToken string, bucket string,
	resultFolder string, pkiIP string, tag string, accessConfig string) string {

	// Customize API endpoint to send request to.
	endpoint := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/zones/%s/instances", proj, config.Zone)

	// Prepare request body.
	reqBody := strings.ReplaceAll(tmplInstanceCreate, "ACS_EVAL_INSERT_NAME", config.Name)
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_ZONE", fmt.Sprintf("projects/%s/zones/%s", proj, config.Zone))
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_MIN_CPU_PLATFORM", config.MinCPUPlatform)
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_MACHINE_TYPE", fmt.Sprintf("projects/%s/zones/%s/machineTypes/%s", proj, config.Zone,
		config.MachineType))
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_TYPE_OF_NODE", config.TypeOfNode)
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_RESULT_FOLDER", resultFolder)
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_EVAL_SCRIPT_TO_PULL", config.EvaluationScript)
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_BINARY_TO_PULL", config.BinaryName)
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_TC_CONFIG", config.ParamsTC)
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_PKI_IP", pkiIP)
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_STARTUP_SCRIPT", fmt.Sprintf("gs://%s/startup.sh", bucket))
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_TAG", tag)
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_SOURCE_IMAGE", fmt.Sprintf("projects/%s/global/images/%s", proj, config.SourceImage))
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_DISK_TYPE", fmt.Sprintf("projects/%s/zones/%s/diskTypes/%s", proj, config.Zone, config.DiskType))
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_DISK_SIZE", config.DiskSize)
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_SUBNETWORK", fmt.Sprintf("projects/%s/regions/%s/subnetworks/default", proj,
		strings.TrimSuffix(config.Zone, "-b")))
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_ACCESS_CONFIG", accessConfig)
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_SERVICE_ACCOUNT", serviceAcc)

	// Create HTTP POST request.
	request, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(reqBody))
	if err != nil {
		fmt.Printf("Failed creating HTTP API request: %v\n", err)
		os.Exit(1)
	}
	request.Header.Set(http.CanonicalHeaderKey("authorization"), fmt.Sprintf("Bearer %s", accessToken))
	request.Header.Set(http.CanonicalHeaderKey("content-type"), "application/json")

	// Send the request to GCP.
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Printf("Failed sending instance create API request: %v\n", err)
		os.Exit(1)
	}

	// Read the response.
	outRaw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed reading from instance create response body: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	return string(outRaw)
}

func checkInstanceReady(name string, zone string) {

	// Execute command to query guest attributes.
	outRaw, _ := exec.Command("/opt/google-cloud-sdk/bin/gcloud", "beta", "compute", "instances", "get-guest-attributes", name,
		"--query-path=acs-eval/initStatus", fmt.Sprintf("--zone=%s", zone)).CombinedOutput()
	out := string(outRaw)

	for !strings.Contains(out, "ThisNodeIsReady") {

		time.Sleep(5 * time.Second)

		outRaw, _ = exec.Command("/opt/google-cloud-sdk/bin/gcloud", "beta", "compute", "instances", "get-guest-attributes", name,
			"--query-path=acs-eval/initStatus", fmt.Sprintf("--zone=%s", zone)).CombinedOutput()
		out = string(outRaw)
	}

	fmt.Printf("Instance %s has completed initialization!\n", name)
}

func runInstance(config *Config, proj string, serviceAcc string, accessToken string, bucket string,
	resultFolder string, pkiIP string, tag string, accessConfig string) {

	// Spawn instance and retrieve response.
	out := spawnInstance(config, proj, serviceAcc, accessToken, bucket, resultFolder, pkiIP, tag, accessConfig)

	// Verify successful machine creation.
	if strings.Contains(out, "RUNNING") {
		fmt.Printf("Instance %s running, waiting for initialization to finish...\n", config.Name)
	} else {
		fmt.Printf("Spawning instance %s returned failure message:\n%s\n", config.Name, out)
		os.Exit(1)
	}
}

func shutdownInstance(confChan <-chan Config, errChan chan<- error, proj string) {

	for config := range confChan {

		// Execute command to shut down compute instance.
		out, err := exec.Command("/opt/google-cloud-sdk/bin/gcloud", "compute", fmt.Sprintf("--project=%s", proj),
			"instances", "delete", config.Name, fmt.Sprintf("--zone=%s", config.Zone)).CombinedOutput()
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
	resultsPath, err := filepath.Abs(*resultsPathFlag)
	if err != nil {
		fmt.Printf("Provided results path '%s' could not be converted to absolute path: %v\n", *resultsPathFlag, err)
		os.Exit(1)
	}

	var configsFile string
	var pkiConfigFile string

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

	// Prepare string containing template
	// for body of instance creation call.
	tmplRaw, err := ioutil.ReadFile("./cmd/runexperiments/gcp-instance-create.tmpl")
	if err != nil {
		fmt.Printf("Failed ingesting file containing instance creation template: %v\n", err)
		os.Exit(1)
	}
	tmplInstanceCreate = string(tmplRaw)

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

	// TODO: Read OAuth token from gcloud.
	out, err := exec.Command("/opt/google-cloud-sdk/bin/gcloud", "auth", "print-access-token").CombinedOutput()
	if err != nil {
		fmt.Printf("Could not obtain OAuth2 access token (error: '%v'):\n%s\n", err, out)
		os.Exit(1)
	}
	accessToken := strings.TrimSpace(string(out))

	var pkiConfig Config
	pkiInternalIP := ""
	pkiExternalIP := ""

	// If zeno: create PKI with decent hardware specs.
	if system == "zeno" {

		// Retrieve file defining PKI configuration.
		pkiConfigFileRel := filepath.Join(*configsPathFlag, "gcloud-mixnet-20-40-30-10_zeno-pki.json")
		pkiConfigFile, err = filepath.Abs(pkiConfigFileRel)
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

		// Spawn PKI of zeno.
		spawnResp := spawnInstance(&pkiConfig, gcloudProject, gcloudServiceAcc, accessToken, gcloudBucket, resultFolder, "irrelevant", "zeno-pki", tmplInstancePublicIP)

		// Verify successful machine creation.
		if !strings.Contains(spawnResp, "RUNNING") {
			fmt.Printf("Spawning PKI of zeno returned failure message:\n%s\n", spawnResp)
			os.Exit(1)
		}

		time.Sleep(10 * time.Second)

		request, err := http.NewRequest(http.MethodGet,
			fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/zones/%s/instances/%s", gcloudProject, pkiConfig.Zone, pkiConfig.Name), nil)
		if err != nil {
			fmt.Printf("Failed retrieving details of PKI of zeno: %v\n", err)
			os.Exit(1)
		}
		request.Header.Set(http.CanonicalHeaderKey("authorization"), fmt.Sprintf("Bearer %s", accessToken))

		// Send the request to GCP.
		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			fmt.Printf("Failed sending instance details API request: %v\n", err)
			os.Exit(1)
		}

		// Read the response.
		outRaw, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Failed reading from instance details response body: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()
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

		fmt.Printf("Successfully spawned PKI of zeno with public IP %s and internal IP %s.\n", pkiExternalIP, pkiInternalIP)
		time.Sleep(45 * time.Second)

		// Ensure initialization has completed.
		checkInstanceReady(pkiConfig.Name, pkiConfig.Zone)
		time.Sleep(5 * time.Second)
	}

	// Spawn all machines.
	fmt.Printf("\nSpawning machines...\n")
	for i := 0; i < len(configs); i++ {
		go runInstance(&configs[i], gcloudProject, gcloudServiceAcc, accessToken, gcloudBucket, resultFolder, pkiInternalIP, "zeno-node", "")
	}

	fmt.Printf("\nWaiting for instances to initialize...\n")
	time.Sleep(90 * time.Second)

	for i := 0; i < 10; i++ {

		go func(configs []Config, i int) {

			start := i * (len(configs) / 10)
			end := (i + 1) * (len(configs) / 10)

			for j := start; j < end; j++ {
				checkInstanceReady(configs[j].Name, configs[j].Zone)
			}

		}(configs, i)
	}

	// Ensure initialization of all
	// machines has completed.
	for i := 0; i < len(configs); i++ {
	}

	fmt.Printf("All machines spawned!\n\n")
	time.Sleep(5 * time.Second)

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

	fmt.Printf("\nType 'yes' and press ENTER to shutdown and delete all resources...")

	shutdown := ""
	stdIn := bufio.NewReader(os.Stdin)

	shutdown, _ = stdIn.ReadString('\n')
	for strings.TrimSpace(shutdown) != "yes" {
		shutdown, _ = stdIn.ReadString('\n')
	}

	if system == "zeno" {
		configs = append(configs, pkiConfig)
	}

	// Prepare channels to send configurations
	// to individual workers and expect responses.
	confChan := make(chan Config, len(configs))
	errChan := make(chan error, len(configs))

	// Spawn deletion workers.
	for i := 0; i < len(configs); i++ {
		go shutdownInstance(confChan, errChan, gcloudProject)
	}

	fmt.Printf("Deleting machines...\n")

	// Shutdown and destroy disks and instances.
	for _, config := range configs {
		confChan <- config
		time.Sleep(1 * time.Second)
	}
	close(confChan)

	for i := 0; i < len(configs); i++ {

		// If any worker threw an error, abort.
		err := <-errChan
		if err != nil {
			fmt.Printf("Subroutine deleting boot disk and shutting down a machine failed: %v\n", err)
			os.Exit(1)
		}
	}
	close(errChan)

	fmt.Printf("All machines deleted!\n\n")

	// Download all files from GCloud bucket
	// to prepared local experiment folder.
	fmt.Printf("Downloading results...")

	// Execute command download result files.
	outRaw, err := exec.Command("/opt/google-cloud-sdk/bin/gsutil", "-m", "cp", "-r",
		fmt.Sprintf("gs://%s/%s/", gcloudBucket, resultFolder), resultsPath).CombinedOutput()
	if err != nil {
		fmt.Printf("Downloading results from GCloud bucket failed (code: '%v'): '%s'", err, outRaw)
		os.Exit(1)
	}

	// Also copy machine configuration files
	// into created results folder.
	outRaw, err = exec.Command("cp", configsFile, fmt.Sprintf("%s/%s/", resultsPath, resultFolder)).CombinedOutput()
	if err != nil {
		fmt.Printf("Copying gcloud config file to results folder failed (code: '%v'): '%s'", err, outRaw)
		os.Exit(1)
	}

	if system == "zeno" {

		outRaw, err = exec.Command("cp", pkiConfigFile, fmt.Sprintf("%s/%s/", resultsPath, resultFolder)).CombinedOutput()
		if err != nil {
			fmt.Printf("Copying gcloud config file for PKI of zeno to results folder failed (code: '%v'): '%s'", err, outRaw)
			os.Exit(1)
		}
	}

	fmt.Printf(" done!\n\n")
	fmt.Printf("Evaluation run completed.\n")
}
