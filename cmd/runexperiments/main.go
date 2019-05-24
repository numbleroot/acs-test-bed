package main

import (
	"bufio"
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
	"sync"
	"time"
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
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_PARTNER", config.Partner)
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
	tried := 0
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Printf("Create API request failed (will try again): %v\n", err)
		tried++
	}

	for err != nil && tried < 10 {

		time.Sleep(1 * time.Second)

		resp, err = http.DefaultClient.Do(request)
		if err != nil {
			fmt.Printf("Create API request failed (will try again): %v\n", err)
		}
	}

	if tried >= 10 {
		fmt.Printf("Create API request failed permanently: %v\n", err)
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

func shutdownInstance(confChan <-chan Config, proj string, accessToken string) {

	for config := range confChan {

		endpoint := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/zones/%s/instances/%s", proj, config.Zone, config.Name)

		// Create HTTP DELETE request.
		request, err := http.NewRequest(http.MethodDelete, endpoint, nil)
		if err != nil {
			fmt.Printf("Failed creating HTTP API request: %v\n", err)
			os.Exit(1)
		}
		request.Header.Set(http.CanonicalHeaderKey("authorization"), fmt.Sprintf("Bearer %s", accessToken))

		// Send the request to GCP.
		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			fmt.Printf("Delete API request failed: %v\n", err)
			continue
		}

		// Read the response.
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Failed reading from instance delete response body: %v\n", err)
			continue
		}
		defer resp.Body.Close()
	}
}

func main() {

	// Expect a number of command-line arguments.
	systemFlag := flag.String("system", "", "Specify which ACS to evaluate: 'zeno', 'vuvuzela', 'pung'.")
	configsPathFlag := flag.String("configsPath", "./gcloud-configs/", "Specify the file system location of the configurations folder for the compute instances.")
	resultsPathFlag := flag.String("resultsPath", "./results/", "Specify the file system location of the top-level results directory to create a new results folder under.")
	gcloudProjectFlag := flag.String("gcloudProj", "", "Supply the GCloud project identifier.")
	gcloudServiceAccFlag := flag.String("gcloudServiceAcc", "", "Supply the GCloud Service Account identifier.")
	gcloudBucketFlag := flag.String("gcloudBucket", "", "Supply the GCloud Storage Bucket to use for the experiments.")
	deleteAllFlag := flag.Bool("deleteAll", false, "Append this flag if the only purpose of this run is to delete all configured instances (caution: permanently deletes them!).")
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
	deleteAll := *deleteAllFlag

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
	var auxConfigFile string
	var auxInternalIP string
	var auxExternalIP string
	var auxConfig Config

	if system == "zeno" {

		// Prepare zeno configurations file for ingestion.
		configsFileRel := filepath.Join(*configsPathFlag, "gcloud-zeno.json")
		configsFile, err = filepath.Abs(configsFileRel)
		if err != nil {
			fmt.Printf("Unable to obtain absolute path to zeno configurations file '%s': %v\n", configsFileRel, err)
			os.Exit(1)
		}

	} else if system == "vuvuzela" {

		// Prepare vuvuzela configurations file for ingestion.
		configsFileRel := filepath.Join(*configsPathFlag, "gcloud-vuvuzela.json")
		configsFile, err = filepath.Abs(configsFileRel)
		if err != nil {
			fmt.Printf("Unable to obtain absolute path to vuvuzela configurations file '%s': %v\n", configsFileRel, err)
			os.Exit(1)
		}

	} else if system == "pung" {

		// Prepare pung configurations file for ingestion.
		configsFileRel := filepath.Join(*configsPathFlag, "gcloud-pung.json")
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

	// Read OAuth token from gcloud.
	out, err := exec.Command("/opt/google-cloud-sdk/bin/gcloud", "auth", "print-access-token").CombinedOutput()
	if err != nil {
		fmt.Printf("Could not obtain OAuth2 access token (error: '%v'):\n%s\n", err, out)
		os.Exit(1)
	}
	accessToken := strings.TrimSpace(string(out))

	if deleteAll {

		// Prepare channels to send configurations
		// to individual workers and expect responses.
		confChan := make(chan Config)

		// Spawn deletion workers.
		for i := 0; i < len(configs); i++ {
			go shutdownInstance(confChan, gcloudProject, accessToken)
		}

		fmt.Printf("WARNING: Deleting all machines...\n")

		// Shutdown and destroy disks and instances.
		for _, config := range configs {
			confChan <- config
		}
		close(confChan)

		time.Sleep(15 * time.Second)
		fmt.Printf("All machines deleted!\n\n")
		os.Exit(0)
	}

	// If zeno: create PKI with decent hardware specs.
	if system == "zeno" {

		// Retrieve file defining PKI configuration.
		pkiConfigFileRel := filepath.Join(*configsPathFlag, "gcloud-zeno-pki.json")
		auxConfigFile, err = filepath.Abs(pkiConfigFileRel)
		if err != nil {
			fmt.Printf("Unable to obtain absolute path to PKI configuration file for zeno '%s': %v\n", pkiConfigFileRel, err)
			os.Exit(1)
		}

		// Ingest GCloud configuration file.
		pkiConfigJSON, err := ioutil.ReadFile(auxConfigFile)
		if err != nil {
			fmt.Printf("Failed ingesting GCloud configuration file for zeno PKI: %v\n", err)
			os.Exit(1)
		}

		// Unmarshal JSON.
		err = json.Unmarshal(pkiConfigJSON, &auxConfig)
		if err != nil {
			fmt.Printf("Error while trying to unmarshal JSON-encoded GCloud configuration for zeno PKI: %v\n", err)
			os.Exit(1)
		}

		// Spawn PKI of zeno.
		spawnResp := spawnInstance(&auxConfig, gcloudProject, gcloudServiceAcc, accessToken, gcloudBucket, resultFolder, "irrelevant", "zeno-pki", tmplInstancePublicIP)

		// Verify successful machine creation.
		if !strings.Contains(spawnResp, "RUNNING") {
			fmt.Printf("Spawning PKI of zeno returned failure message:\n%s\n", spawnResp)
			os.Exit(1)
		}

		time.Sleep(10 * time.Second)

		request, err := http.NewRequest(http.MethodGet,
			fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/zones/%s/instances/%s", gcloudProject, auxConfig.Zone, auxConfig.Name), nil)
		if err != nil {
			fmt.Printf("Failed retrieving details of PKI of zeno: %v\n", err)
			os.Exit(1)
		}
		request.Header.Set(http.CanonicalHeaderKey("authorization"), fmt.Sprintf("Bearer %s", accessToken))

		// Send the request to GCP.
		tried := 0
		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			fmt.Printf("Details API request failed (will try again): %v\n", err)
			tried++
		}

		for err != nil && tried < 10 {

			resp, err = http.DefaultClient.Do(request)
			if err != nil {
				fmt.Printf("Details API request failed (will try again): %v\n", err)
			}
		}

		if tried >= 10 {
			fmt.Printf("Details API request failed permanently: %v\n", err)
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
				auxInternalIP = pkiIPs[i]
			} else {
				auxExternalIP = pkiIPs[i]
			}
		}

		fmt.Printf("Successfully spawned PKI of zeno with public IP %s and internal IP %s.\n", auxExternalIP, auxInternalIP)
		time.Sleep(45 * time.Second)

		// Ensure initialization has completed.
		checkInstanceReady(auxConfig.Name, auxConfig.Zone)
		time.Sleep(5 * time.Second)

	} else if system == "pung" {

		// Retrieve file defining pung's server configuration.
		serverConfigFileRel := filepath.Join(*configsPathFlag, "gcloud-pung-server.json")
		auxConfigFile, err = filepath.Abs(serverConfigFileRel)
		if err != nil {
			fmt.Printf("Unable to obtain absolute path to server configuration file for pung '%s': %v\n", serverConfigFileRel, err)
			os.Exit(1)
		}

		// Ingest GCloud configuration file.
		auxConfigJSON, err := ioutil.ReadFile(auxConfigFile)
		if err != nil {
			fmt.Printf("Failed ingesting GCloud configuration file for pung server: %v\n", err)
			os.Exit(1)
		}

		// Unmarshal JSON.
		err = json.Unmarshal(auxConfigJSON, &auxConfig)
		if err != nil {
			fmt.Printf("Error while trying to unmarshal JSON-encoded GCloud configuration for pung server: %v\n", err)
			os.Exit(1)
		}

		// Spawn pung's server.
		spawnResp := spawnInstance(&auxConfig, gcloudProject, gcloudServiceAcc, accessToken, gcloudBucket, resultFolder, "irrelevant", "pung-server", tmplInstancePublicIP)

		// Verify successful machine creation.
		if !strings.Contains(spawnResp, "RUNNING") {
			fmt.Printf("Spawning pung's server returned failure message:\n%s\n", spawnResp)
			os.Exit(1)
		}

		time.Sleep(10 * time.Second)

		request, err := http.NewRequest(http.MethodGet,
			fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/zones/%s/instances/%s", gcloudProject, auxConfig.Zone, auxConfig.Name), nil)
		if err != nil {
			fmt.Printf("Failed retrieving details of pung's server: %v\n", err)
			os.Exit(1)
		}
		request.Header.Set(http.CanonicalHeaderKey("authorization"), fmt.Sprintf("Bearer %s", accessToken))

		// Send the request to GCP.
		tried := 0
		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			fmt.Printf("Details API request failed (will try again): %v\n", err)
			tried++
		}

		for err != nil && tried < 10 {

			resp, err = http.DefaultClient.Do(request)
			if err != nil {
				fmt.Printf("Details API request failed (will try again): %v\n", err)
			}
		}

		if tried >= 10 {
			fmt.Printf("Details API request failed permanently: %v\n", err)
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
		// the assigned IP addresses of the server.
		regexpIP, err := regexp.Compile("[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}")
		if err != nil {
			fmt.Printf("Failed to match internal IP address of pung's server instance: %v\n", err)
			os.Exit(1)
		}

		pkiIPs := regexpIP.FindAllString(out, -1)
		for i := range pkiIPs {

			if strings.HasPrefix(pkiIPs[i], "10.") {
				auxInternalIP = pkiIPs[i]
			} else {
				auxExternalIP = pkiIPs[i]
			}
		}

		fmt.Printf("Successfully spawned pung's server with public IP %s and internal IP %s.\n", auxExternalIP, auxInternalIP)
		time.Sleep(45 * time.Second)

		// Ensure initialization has completed.
		checkInstanceReady(auxConfig.Name, auxConfig.Zone)
		time.Sleep(5 * time.Second)
	}

	// Spawn all machines.
	fmt.Printf("\nSpawning machines...\n")
	for i := 0; i < len(configs); i++ {
		go runInstance(&configs[i], gcloudProject, gcloudServiceAcc, accessToken, gcloudBucket, resultFolder, auxInternalIP, "client", "")
		time.Sleep(250 * time.Millisecond)
	}

	time.Sleep(10 * time.Second)
	fmt.Printf("\nWaiting for instances to initialize...\n")
	time.Sleep(80 * time.Second)

	wg := &sync.WaitGroup{}

	for i := 0; i < 10; i++ {

		wg.Add(1)

		// Run routine that checks all instances to be
		// up and running highly parallel.

		go func(wg *sync.WaitGroup, configs []Config, idx int) {

			defer wg.Done()

			for j := (idx * 200); j < ((idx + 1) * 200); j++ {

				if j < len(configs) {
					checkInstanceReady(configs[j].Name, configs[j].Zone)
				}
			}

		}(wg, configs, i)
	}

	// Catch all remaining configs.
	for i := (10 * 200); i < len(configs); i++ {
		checkInstanceReady(configs[i].Name, configs[i].Zone)
	}

	wg.Wait()

	fmt.Printf("All machines spawned!\n\n")

	// If zeno: send PKI signal to start.
	if system == "zeno" {

		time.Sleep(15 * time.Second)
		fmt.Printf("Signaling zeno's PKI to start epoch.\n")

		// Connect to control plane address used
		// only for evaluation purposes.
		ctrlConn, err := net.Dial("tcp", fmt.Sprintf("%s:26345", auxExternalIP))
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

	if system == "zeno" || system == "pung" {
		configs = append(configs, auxConfig)
	}

	// Prepare channels to send configurations
	// to individual workers and expect responses.
	confChan := make(chan Config)

	// Spawn deletion workers.
	for i := 0; i < len(configs); i++ {
		go shutdownInstance(confChan, gcloudProject, accessToken)
		time.Sleep(250 * time.Millisecond)
	}

	fmt.Printf("WARNING: Deleting all machines...\n")

	// Shutdown and destroy disks and instances.
	for _, config := range configs {
		confChan <- config
	}
	close(confChan)

	time.Sleep(15 * time.Second)
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

	if system == "zeno" || system == "pung" {

		outRaw, err = exec.Command("cp", auxConfigFile, fmt.Sprintf("%s/%s/", resultsPath, resultFolder)).CombinedOutput()
		if err != nil {
			fmt.Printf("Copying auxiliary gcloud config file to results folder failed (code: '%v'): '%s'", err, outRaw)
			os.Exit(1)
		}
	}

	fmt.Printf(" done!\n\n")
	fmt.Printf("Evaluation run completed.\n")
}
