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
	"time"
)

// Config describes one compute instance
// exhaustively for reproducibility.
type Config struct {
	Name                   string `json:"Name"`
	Partner                string `json:"Partner"`
	Zone                   string `json:"Zone"`
	MinCPUPlatform         string `json:"MinCPUPlatform"`
	MachineType            string `json:"MachineType"`
	TypeOfNode             string `json:"TypeOfNode"`
	EvaluationScript       string `json:"EvaluationScript"`
	BinaryName             string `json:"BinaryName"`
	SourceImage            string `json:"SourceImage"`
	DiskType               string `json:"DiskType"`
	DiskSize               string `json:"DiskSize"`
	NetTroublesIfApplied   string `json:"NetTroublesIfApplied"`
	ZenoMixKilledIfApplied string `json:"ZenoMixKilledIfApplied"`
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
	resultFolder string, pkiIP string, tag string, accessConfig string, tcEmulNetTroubles bool, killZenoMixesInRound int) string {

	// Customize API endpoint to send request to.
	endpoint := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/zones/%s/instances", proj, config.Zone)

	// Prepare request body.
	reqBody := strings.ReplaceAll(tmplInstanceCreate, "ACS_EVAL_INSERT_GCP_MACHINE_NAME", fmt.Sprintf("%s-%s", config.Name, resultFolder))
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_ZONE", fmt.Sprintf("projects/%s/zones/%s", proj, config.Zone))
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_MIN_CPU_PLATFORM", config.MinCPUPlatform)
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_MACHINE_TYPE", fmt.Sprintf("projects/%s/zones/%s/machineTypes/%s", proj, config.Zone,
		config.MachineType))
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_NAME", config.Name)
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_PARTNER", config.Partner)
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_TYPE_OF_NODE", config.TypeOfNode)
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_RESULT_FOLDER", resultFolder)
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_EVAL_SCRIPT_TO_PULL", config.EvaluationScript)
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_BINARY_TO_PULL", config.BinaryName)
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_PKI_IP", pkiIP)
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_STARTUP_SCRIPT", fmt.Sprintf("gs://%s/startup.sh", bucket))
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_TAG", tag)
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_SOURCE_IMAGE", fmt.Sprintf("projects/%s/global/images/%s", proj, config.SourceImage))
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_DISK_TYPE", fmt.Sprintf("projects/%s/zones/%s/diskTypes/%s", proj, config.Zone, config.DiskType))
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_DISK_SIZE", config.DiskSize)
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_SUBNETWORK", fmt.Sprintf("projects/%s/regions/%s/subnetworks/default", proj,
		strings.TrimSuffix(config.Zone, "-b")))
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_SERVICE_ACCOUNT", serviceAcc)

	if strings.Contains(config.Name, "mixnet-") {
		reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_ACCESS_CONFIG", tmplInstancePublicIP)
	} else {
		reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_ACCESS_CONFIG", accessConfig)
	}

	// If flag set to true, append the tc configuration
	// parameters to instance metadata.
	if tcEmulNetTroubles {
		reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_TC_CONFIG", config.NetTroublesIfApplied)
	} else {
		reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_TC_CONFIG", "none")
	}

	// If flag set to true, signal spawning zeno nodes
	// that in all but the first cascade the second mix
	// node is instructed to crash in the specified round.
	if killZenoMixesInRound > 0 {
		reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_KILL_ZENO_MIXES_IN_ROUND", fmt.Sprintf("%d", killZenoMixesInRound))
	} else {
		reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_KILL_ZENO_MIXES_IN_ROUND", "-1")
	}

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
	for err != nil && tried < 10 {

		tried++
		fmt.Printf("Create API request failed (will try again): %v\n", err)
		time.Sleep(1 * time.Second)

		resp, err = http.DefaultClient.Do(request)
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

func checkInstanceReady(proj string, accessToken string, name string, zone string, resultFolder string) {

	instName := fmt.Sprintf("%s-%s", name, resultFolder)

	// Customize API endpoint to send request to.
	endpoint := fmt.Sprintf("https://www.googleapis.com/compute/beta/projects/%s/zones/%s/instances/%s/getGuestAttributes?queryPath=acs-eval%%2FinitStatus&variableKey=initStatus",
		proj, zone, instName)

	// Create HTTP GET request.
	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		fmt.Printf("Failed creating HTTP API request: %v\n", err)
		os.Exit(1)
	}
	request.Header.Set(http.CanonicalHeaderKey("authorization"), fmt.Sprintf("Bearer %s", accessToken))

	// Send the request to GCP.
	resp, err := http.DefaultClient.Do(request)
	for err != nil {

		fmt.Printf("Guest attributes API request failed (will try again): %v\n", err)
		time.Sleep(1 * time.Second)

		resp, err = http.DefaultClient.Do(request)
	}

	// Read the response.
	outRaw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed reading from instance guest attributes response body: %v\n", err)
		os.Exit(1)
	}
	out := string(outRaw)

	for !strings.Contains(out, "ThisNodeIsReady") {

		time.Sleep(5 * time.Second)

		resp, err = http.DefaultClient.Do(request)
		for err != nil {

			fmt.Printf("Guest attributes API request failed (will try again): %v\n", err)
			time.Sleep(1 * time.Second)

			resp, err = http.DefaultClient.Do(request)
		}

		// Read the response.
		outRaw, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Failed reading from instance guest attributes response body: %v\n", err)
			os.Exit(1)
		}
		out = string(outRaw)
	}

	resp.Body.Close()

	fmt.Printf("Instance %s has completed initialization!\n", name)
}

func runInstance(config *Config, proj string, serviceAcc string, accessToken string, bucket string,
	resultFolder string, pkiIP string, tag string, accessConfig string, tcEmulNetTroubles bool, killZenoMixesInRound int) {

	// Spawn instance and retrieve response.
	out := spawnInstance(config, proj, serviceAcc, accessToken, bucket, resultFolder, pkiIP,
		tag, accessConfig, tcEmulNetTroubles, killZenoMixesInRound)

	// Verify successful machine creation.
	if strings.Contains(out, "RUNNING") {
		fmt.Printf("Instance %s running, waiting for initialization to finish...\n", config.Name)
	} else {
		fmt.Printf("Spawning instance %s returned failure message:\n%s\n", config.Name, out)
		os.Exit(1)
	}
}

func shutdownInstance(config *Config, proj string, accessToken string, resultFolder string) {

	fmt.Printf("Deleting machine %s\n", config.Name)

	instName := fmt.Sprintf("%s-%s", config.Name, resultFolder)
	endpoint := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/zones/%s/instances/%s", proj, config.Zone, instName)

	// Create HTTP DELETE request.
	request, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		fmt.Printf("Failed creating HTTP API request: %v\n", err)
		os.Exit(1)
	}
	request.Header.Set(http.CanonicalHeaderKey("authorization"), fmt.Sprintf("Bearer %s", accessToken))

	// Send the request to GCP.
	resp, err := http.DefaultClient.Do(request)
	for err != nil {

		fmt.Printf("Delete API request failed (will try again): %v\n", err)
		time.Sleep(1 * time.Second)

		resp, err = http.DefaultClient.Do(request)
	}

	// Read the response.
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed reading from instance delete response body: %v\n", err)
		return
	}
	defer resp.Body.Close()
}

func main() {

	// Expect a number of command-line arguments.
	systemFlag := flag.String("system", "", "Specify which ACS to evaluate: 'zeno', 'vuvuzela', 'pung'.")
	configsPathFlag := flag.String("configsPath", "./gcloud-configs/", "Specify the file system location of the configurations folder for the compute instances.")
	resultsPathFlag := flag.String("resultsPath", "./results/", "Specify the file system location of the top-level results directory to create a new results folder under.")
	gcloudProjectFlag := flag.String("gcloudProj", "", "Supply the GCloud project identifier.")
	gcloudServiceAccFlag := flag.String("gcloudServiceAcc", "", "Supply the GCloud Service Account identifier.")
	gcloudBucketFlag := flag.String("gcloudBucket", "", "Supply the GCloud Storage Bucket to use for the experiments.")
	tcEmulNetTroublesFlag := flag.Bool("tcEmulNetTroubles", false, "Append this flag to emulate a network trouble in 3 out of all zones.")
	killZenoMixesInRoundFlag := flag.Int("killZenoMixesInRound", -1, "If specific mix nodes in all but one zeno cascade are supposed to crash, specify the round in which that shall happen.")
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
	tcEmulNetTroubles := *tcEmulNetTroublesFlag
	killZenoMixesInRound := *killZenoMixesInRoundFlag

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
		spawnResp := spawnInstance(&auxConfig, gcloudProject, gcloudServiceAcc, accessToken, gcloudBucket,
			resultFolder, "irrelevant", "zeno-pki", tmplInstancePublicIP, false, 0)

		// Verify successful machine creation.
		if !strings.Contains(spawnResp, "RUNNING") {
			fmt.Printf("Spawning PKI of zeno returned failure message:\n%s\n", spawnResp)
			os.Exit(1)
		}

		time.Sleep(10 * time.Second)

		request, err := http.NewRequest(http.MethodGet,
			fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/zones/%s/instances/%s",
				gcloudProject, auxConfig.Zone, fmt.Sprintf("%s-%s", auxConfig.Name, resultFolder)), nil)
		if err != nil {
			fmt.Printf("Failed retrieving details of PKI of zeno: %v\n", err)
			os.Exit(1)
		}
		request.Header.Set(http.CanonicalHeaderKey("authorization"), fmt.Sprintf("Bearer %s", accessToken))

		// Send the request to GCP.
		resp, err := http.DefaultClient.Do(request)
		for err != nil {

			fmt.Printf("Details API request failed (will try again): %v\n", err)
			time.Sleep(1 * time.Second)

			resp, err = http.DefaultClient.Do(request)
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
		time.Sleep(30 * time.Second)

		// Ensure initialization has completed.
		checkInstanceReady(gcloudProject, accessToken, auxConfig.Name, auxConfig.Zone, resultFolder)

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
		spawnResp := spawnInstance(&auxConfig, gcloudProject, gcloudServiceAcc, accessToken, gcloudBucket, resultFolder,
			"irrelevant", "pung-server", tmplInstancePublicIP, false, 0)

		// Verify successful machine creation.
		if !strings.Contains(spawnResp, "RUNNING") {
			fmt.Printf("Spawning pung's server returned failure message:\n%s\n", spawnResp)
			os.Exit(1)
		}

		time.Sleep(10 * time.Second)

		request, err := http.NewRequest(http.MethodGet,
			fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/zones/%s/instances/%s",
				gcloudProject, auxConfig.Zone, fmt.Sprintf("%s-%s", auxConfig.Name, resultFolder)), nil)
		if err != nil {
			fmt.Printf("Failed retrieving details of pung's server: %v\n", err)
			os.Exit(1)
		}
		request.Header.Set(http.CanonicalHeaderKey("authorization"), fmt.Sprintf("Bearer %s", accessToken))

		// Send the request to GCP.
		resp, err := http.DefaultClient.Do(request)
		for err != nil {

			fmt.Printf("Details API request failed (will try again): %v\n", err)
			time.Sleep(1 * time.Second)

			resp, err = http.DefaultClient.Do(request)
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
		time.Sleep(30 * time.Second)

		// Ensure initialization has completed.
		checkInstanceReady(gcloudProject, accessToken, auxConfig.Name, auxConfig.Zone, resultFolder)
	}

	// Spawn all machines.
	fmt.Printf("\nSpawning machines...\n")
	for i := 0; i < (len(configs) - 10); i++ {
		go runInstance(&configs[i], gcloudProject, gcloudServiceAcc, accessToken, gcloudBucket, resultFolder,
			auxInternalIP, "client", "", tcEmulNetTroubles, killZenoMixesInRound)
		time.Sleep(50 * time.Millisecond)
	}

	// Ensure last ten instances are spawned sequentially.
	runInstance(&configs[(len(configs)-10)], gcloudProject, gcloudServiceAcc, accessToken, gcloudBucket, resultFolder,
		auxInternalIP, "client", "", tcEmulNetTroubles, killZenoMixesInRound)
	runInstance(&configs[(len(configs)-9)], gcloudProject, gcloudServiceAcc, accessToken, gcloudBucket, resultFolder,
		auxInternalIP, "client", "", tcEmulNetTroubles, killZenoMixesInRound)
	runInstance(&configs[(len(configs)-8)], gcloudProject, gcloudServiceAcc, accessToken, gcloudBucket, resultFolder,
		auxInternalIP, "client", "", tcEmulNetTroubles, killZenoMixesInRound)
	runInstance(&configs[(len(configs)-7)], gcloudProject, gcloudServiceAcc, accessToken, gcloudBucket, resultFolder,
		auxInternalIP, "client", "", tcEmulNetTroubles, killZenoMixesInRound)
	runInstance(&configs[(len(configs)-6)], gcloudProject, gcloudServiceAcc, accessToken, gcloudBucket, resultFolder,
		auxInternalIP, "client", "", tcEmulNetTroubles, killZenoMixesInRound)
	runInstance(&configs[(len(configs)-5)], gcloudProject, gcloudServiceAcc, accessToken, gcloudBucket, resultFolder,
		auxInternalIP, "client", "", tcEmulNetTroubles, killZenoMixesInRound)
	runInstance(&configs[(len(configs)-4)], gcloudProject, gcloudServiceAcc, accessToken, gcloudBucket, resultFolder,
		auxInternalIP, "client", "", tcEmulNetTroubles, killZenoMixesInRound)
	runInstance(&configs[(len(configs)-3)], gcloudProject, gcloudServiceAcc, accessToken, gcloudBucket, resultFolder,
		auxInternalIP, "client", "", tcEmulNetTroubles, killZenoMixesInRound)
	runInstance(&configs[(len(configs)-2)], gcloudProject, gcloudServiceAcc, accessToken, gcloudBucket, resultFolder,
		auxInternalIP, "client", "", tcEmulNetTroubles, killZenoMixesInRound)
	runInstance(&configs[(len(configs)-1)], gcloudProject, gcloudServiceAcc, accessToken, gcloudBucket, resultFolder,
		auxInternalIP, "client", "", tcEmulNetTroubles, killZenoMixesInRound)

	time.Sleep(5 * time.Second)
	fmt.Printf("\nWaiting for instances to initialize...\n")
	time.Sleep(30 * time.Second)

	for i := 0; i < (len(configs) - 10); i++ {
		go checkInstanceReady(gcloudProject, accessToken, configs[i].Name, configs[i].Zone, resultFolder)
		time.Sleep(50 * time.Millisecond)
	}

	// Ensure last ten instances are checked sequentially.
	checkInstanceReady(gcloudProject, accessToken, configs[(len(configs)-10)].Name, configs[(len(configs)-10)].Zone, resultFolder)
	checkInstanceReady(gcloudProject, accessToken, configs[(len(configs)-9)].Name, configs[(len(configs)-9)].Zone, resultFolder)
	checkInstanceReady(gcloudProject, accessToken, configs[(len(configs)-8)].Name, configs[(len(configs)-8)].Zone, resultFolder)
	checkInstanceReady(gcloudProject, accessToken, configs[(len(configs)-7)].Name, configs[(len(configs)-7)].Zone, resultFolder)
	checkInstanceReady(gcloudProject, accessToken, configs[(len(configs)-6)].Name, configs[(len(configs)-6)].Zone, resultFolder)
	checkInstanceReady(gcloudProject, accessToken, configs[(len(configs)-5)].Name, configs[(len(configs)-5)].Zone, resultFolder)
	checkInstanceReady(gcloudProject, accessToken, configs[(len(configs)-4)].Name, configs[(len(configs)-4)].Zone, resultFolder)
	checkInstanceReady(gcloudProject, accessToken, configs[(len(configs)-3)].Name, configs[(len(configs)-3)].Zone, resultFolder)
	checkInstanceReady(gcloudProject, accessToken, configs[(len(configs)-2)].Name, configs[(len(configs)-2)].Zone, resultFolder)
	checkInstanceReady(gcloudProject, accessToken, configs[(len(configs)-1)].Name, configs[(len(configs)-1)].Zone, resultFolder)

	fmt.Printf("All machines spawned!\n\n")

	// If zeno: send PKI signal to start.
	if system == "zeno" {

		time.Sleep(1 * time.Second)
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

	fmt.Printf("\nExperiment run for %s\n", resultFolder)
	fmt.Printf("Type 'yes' and press ENTER to shutdown and delete all resources... ")

	shutdown := ""
	stdIn := bufio.NewReader(os.Stdin)

	shutdown, _ = stdIn.ReadString('\n')
	for strings.TrimSpace(shutdown) != "yes" {
		shutdown, _ = stdIn.ReadString('\n')
	}

	if system == "zeno" || system == "pung" {
		configs = append(configs, auxConfig)
	}

	fmt.Printf("Going to delete all %d machines\n", len(configs))

	// Spawn deletion workers.
	for i := 0; i < (len(configs) - 10); i++ {
		go shutdownInstance(&configs[i], gcloudProject, accessToken, resultFolder)
		time.Sleep(50 * time.Millisecond)
	}

	// Ensure last ten instances are deleted sequentially.
	shutdownInstance(&configs[(len(configs)-10)], gcloudProject, accessToken, resultFolder)
	shutdownInstance(&configs[(len(configs)-9)], gcloudProject, accessToken, resultFolder)
	shutdownInstance(&configs[(len(configs)-8)], gcloudProject, accessToken, resultFolder)
	shutdownInstance(&configs[(len(configs)-7)], gcloudProject, accessToken, resultFolder)
	shutdownInstance(&configs[(len(configs)-6)], gcloudProject, accessToken, resultFolder)
	shutdownInstance(&configs[(len(configs)-5)], gcloudProject, accessToken, resultFolder)
	shutdownInstance(&configs[(len(configs)-4)], gcloudProject, accessToken, resultFolder)
	shutdownInstance(&configs[(len(configs)-3)], gcloudProject, accessToken, resultFolder)
	shutdownInstance(&configs[(len(configs)-2)], gcloudProject, accessToken, resultFolder)
	shutdownInstance(&configs[(len(configs)-1)], gcloudProject, accessToken, resultFolder)

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
