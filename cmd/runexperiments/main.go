package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Config describes one compute instance
// exhaustively for reproducibility.
//
// TODO: Adjust content after genconfigs is fixed.
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

// Exp contains all experiment information
// the operator uses to manage experiments.
type Exp struct {
	ID                     string    `json:"id"`
	Created                string    `json:"created"`
	System                 string    `json:"system"`
	Concluded              bool      `json:"concluded"`
	ResultFolder           string    `json:"resultFolder"`
	NetTroublesIfApplied   string    `json:"netTroublesIfApplied"`
	ZenoMixKilledIfApplied string    `json:"zenoMixKilledIfApplied"`
	Progress               []string  `json:"progress"`
	Servers                []*Worker `json:"servers"`
	Clients                []*Worker `json:"clients"`
}

// Worker describes one compute instance
// exhaustively for reproducibility.
type Worker struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Address        string `json:"address"`
	Status         string `json:"status"`
	Zone           string `json:"zone"`
	MinCPUPlatform string `json:"minCPUPlatform"`
	MachineType    string `json:"machineType"`
	TypeOfNode     string `json:"typeOfNode"`
	BinaryName     string `json:"binaryName"`
	SourceImage    string `json:"sourceImage"`
	DiskType       string `json:"diskType"`
	DiskSize       string `json:"diskSize"`
}

// PrettyPrint writes the experiment
// human-readable to STDOUT.
func (exp *Exp) PrettyPrint() {

	fmt.Printf("---\nExperiment for system %s with ID %s, created at %s:\n", exp.System, exp.ID, exp.Created)
	fmt.Printf("\tConcluded? %v\n", exp.Concluded)
	fmt.Printf("\tServers: %d\n", len(exp.Servers))
	fmt.Printf("\tClients: %d\n", len(exp.Clients))
	fmt.Printf("\tTC parameters (applied, if non-empty): '%s'\n", exp.NetTroublesIfApplied)
	fmt.Printf("\tRound in which zeno mixes will be terminated (applied, if non-empty): '%s'\n", exp.ZenoMixKilledIfApplied)

	fmt.Printf("\n\tPROGRESS:\n")
	for i := range exp.Progress {
		fmt.Printf("\t\t%s\n", exp.Progress[i])
	}

	fmt.Printf("---\n")
}

func init() {

	// Enable TLS 1.3.
	if os.Getenv("GODEBUG") == "" {
		os.Setenv("GODEBUG", "tls13=1")
	} else {
		os.Setenv("GODEBUG", fmt.Sprintf("%s,tls13=1", os.Getenv("GODEBUG")))
	}
}

func main() {

	// Expect a number of command-line arguments.
	systemFlag := flag.String("system", "", "Specify which ACS to evaluate: 'zeno', 'vuvuzela', 'pung'.")
	configsPathFlag := flag.String("configsPath", "./gcloud-configs/", "Specify the file system location of the configurations folder for the compute instances.")
	operatorAddrFlag := flag.String("operatorAddr", "127.0.0.1:443", "Supply the address at which the TLS API of the operator is reachable.")
	certFileFlag := flag.String("certFile", "./operator-cert.pem", "Specify the file system location of the self-signed TLS certificate of the operator.")
	resultsPathFlag := flag.String("resultsPath", "./results/", "Specify the file system location of the top-level results directory to create a new results folder under.")
	gcloudBucketFlag := flag.String("gcloudBucket", "", "Supply the GCloud Storage Bucket to use for the experiments.")
	tcEmulNetTroublesFlag := flag.Bool("tcEmulNetTroubles", false, "Append this flag to emulate a network trouble in 3 out of all zones.")
	killZenoMixesInRoundFlag := flag.Int("killZenoMixesInRound", -1, "If specific mix nodes in all but one zeno cascade are supposed to crash, specify the round in which that shall happen.")
	flag.Parse()

	// Enforce arguments to be set.
	if *systemFlag == "" || *gcloudBucketFlag == "" {
		fmt.Printf("Missing arguments, please provide values for all flags: '-system' and '-gcloudBucket'.\n")
		os.Exit(1)
	}

	system := strings.ToLower(*systemFlag)
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

	// Create new empty cert pool.
	certPool := x509.NewCertPool()

	// Load supplied certificate.
	cert, err := ioutil.ReadFile(*certFileFlag)
	if err != nil {
		fmt.Printf("Could not load operator's TLS certificate: %v\n", err)
		os.Exit(1)
	}

	// Attempt to add loaded certificate to pool.
	ok := certPool.AppendCertsFromPEM(cert)
	if !ok {
		fmt.Printf("Failed to append PEM certificate to empty pool.\n")
		os.Exit(1)
	}

	// Prepare an HTTPS transport struct that uses
	// modern schemes and is set to accept the self-signed
	// server TLS certificate.
	client := &http.Client{

		Transport: &http.Transport{

			TLSClientConfig: &tls.Config{
				RootCAs:            certPool,
				InsecureSkipVerify: false,
				MinVersion:         tls.VersionTLS13,
				CurvePreferences:   []tls.CurveID{tls.X25519},
			},
		},
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
		fmt.Printf("Error while trying to unmarshal JSON-encoded GCloud configuration: %v\n", err)
		os.Exit(1)
	}

	// Prepare buffer of JSON payload to be
	// attached to the HTTPS request.
	reqBodyBuf := new(bytes.Buffer)
	err = json.NewEncoder(reqBodyBuf).Encode(configs)
	if err != nil {
		fmt.Printf("Failed to encode JSON payload for new experiment into buffer: %v\n", err)
		os.Exit(1)
	}

	// Prepare HTTP JSON request to operator
	// to launch a new experiment.
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("https://%s/public/experiments/new", *operatorAddrFlag), reqBodyBuf)
	if err != nil {
		fmt.Printf("Failed creating HTTPS API request for new experiment: %v\n", err)
		os.Exit(1)
	}
	req.Header.Set(http.CanonicalHeaderKey("Authorization"), "UniverseOfLoopholes")
	req.Header.Set(http.CanonicalHeaderKey("Content-Type"), "application/json")

	// Send experiment request.
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed sending HTTPS API request for new experiment: %v\n", err)
		os.Exit(1)
	}

	// Read the response.
	exp := &Exp{}
	err = json.NewDecoder(resp.Body).Decode(exp)
	if err != nil {
		fmt.Printf("Failed decoding response from HTTPS API request for new experiment to JSON: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	fmt.Printf("Operator responded to request for new experiment with:\n")
	exp.PrettyPrint()

	// Loop over user input. Await either status
	// request or experiment termination input.

	fmt.Printf("Type 's' for 'status' or 't' for 'terminate' and press ENTER.\n")
	fmt.Printf("This either requests the current status of the experiment or confirms shutdown and deletion of all experiment resources... ")

	input := ""
	stdIn := bufio.NewReader(os.Stdin)

	input, _ = stdIn.ReadString('\n')
	for strings.TrimSpace(input) != "t" {

		if strings.TrimSpace(input) == "s" {

			// Request current status of experiment.
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://%s/public/experiments/%s/status", *operatorAddrFlag, exp.ID), nil)
			if err != nil {
				fmt.Printf("Failed creating HTTPS API request for status of experiment: %v\n", err)
				os.Exit(1)
			}
			req.Header.Set(http.CanonicalHeaderKey("Authorization"), "UniverseOfLoopholes")
			req.Header.Set(http.CanonicalHeaderKey("Content-Type"), "application/json")

			// Send status request.
			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("Failed sending HTTPS API request for status of experiment: %v\n", err)
				os.Exit(1)
			}

			// Read the response.
			expStatus := &Exp{}
			err = json.NewDecoder(resp.Body).Decode(expStatus)
			if err != nil {
				fmt.Printf("Failed decoding response from HTTPS API request for status of experiment to JSON: %v\n", err)
				os.Exit(1)
			}
			defer resp.Body.Close()

			fmt.Printf("\nStatus of experiment %s:\n", exp.ID)
			expStatus.PrettyPrint()
		}

		input, _ = stdIn.ReadString('\n')
	}

	fmt.Printf("\nWill instruct operator to terminate experiment...")

	// Request termination of experiment.
	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("https://%s/public/experiments/%s/terminate", *operatorAddrFlag, exp.ID), nil)
	if err != nil {
		fmt.Printf("Failed creating HTTPS API request to terminate experiment: %v\n", err)
		os.Exit(1)
	}
	req.Header.Set(http.CanonicalHeaderKey("Authorization"), "UniverseOfLoopholes")
	req.Header.Set(http.CanonicalHeaderKey("Content-Type"), "application/json")

	// Send termination request.
	_, err = client.Do(req)
	if err != nil {
		fmt.Printf("Failed sending HTTPS API request to terminate experiment: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf(" done!\n")

	// Download all files from GCloud bucket
	// to prepared local experiment folder.
	fmt.Printf("\nDownloading results...")

	// Execute command to download result files.
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

	fmt.Printf(" done!\n")

	if !tcEmulNetTroubles {

		if killZenoMixesInRound == -1 {
			fmt.Printf("\nEvaluation run %s for 01_tc-off_proc-off completed\n", resultFolder)
		} else {
			fmt.Printf("\nEvaluation run %s for 03_tc-off_proc-on completed\n", resultFolder)
		}

	} else {

		if killZenoMixesInRound == -1 {
			fmt.Printf("\nEvaluation run %s for 02_tc-on_proc-off completed\n", resultFolder)
		} else {
			fmt.Printf("\nEvaluation run %s for 04_tc-on_proc-on completed\n", resultFolder)
		}
	}
}
