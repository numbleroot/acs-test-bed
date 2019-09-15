package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/emicklei/go-restful"
)

// Operator describes the node in the
// public cloud that orchestrates the
// execution of experiments from start
// to finish.
type Operator struct {
	sync.Mutex

	GCloudBucket string

	InternalListenAddr string
	InternalSrv        *restful.WebService
	InternalChan       chan string

	PublicListenAddr string
	PublicSrv        *restful.WebService
	PublicChan       chan string

	ExpInProgress string
	Exps          map[string]*Exp
}

// Exp contains all information relevant
// to monitoring an experiment.
type Exp struct {
	ID           string    `json:"id"`
	InitTime     time.Time `json:"initTime"`
	System       string    `json:"system"`
	Concluded    bool      `json:"concluded"`
	ResultFolder string    `json:"resultFolder"`
	Progress     []string  `json:"progress"`

	ServersSpawned map[string]*Worker `json:"serversSpawned"`
	ServersUsed    map[string]*Worker `json:"serversUsed"`

	ClientsSpawned map[string]*Worker `json:"clientsSpawned"`
	ClientsUsed    map[string]*Worker `json:"clientsUsed"`
}

// Enable TLS 1.3.
func init() {
	os.Setenv("GODEBUG", fmt.Sprintf("%s,tls13=1", os.Getenv("GODEBUG")))
}

// GetTLSMaterial downloads the TLS certificate
// and key to be used by the operator.
func (op *Operator) GetTLSMaterial() error {

	outRaw, err := exec.Command("/opt/google-cloud-sdk/bin/gsutil", "cp",
		fmt.Sprintf("gs://%s/operator-cert.pem", op.GCloudBucket), "operator-cert.pem").CombinedOutput()
	if err != nil {
		return fmt.Errorf("downloading 'operator-cert.pem' from GCloud bucket failed (code: '%v'): '%s'", err, outRaw)
	}

	outRaw, err = exec.Command("/opt/google-cloud-sdk/bin/gsutil", "cp",
		fmt.Sprintf("gs://%s/operator-key.pem", op.GCloudBucket), "operator-key.pem").CombinedOutput()
	if err != nil {
		return fmt.Errorf("downloading 'operator-key.pem' from GCloud bucket failed (code: '%v'): '%s'", err, outRaw)
	}

	// Ensure appropriate permission on
	// sensistive files.
	err = os.Chmod("operator-cert.pem", 0644)
	if err != nil {
		return fmt.Errorf("failed to set appropriate permissions: %v", err)
	}

	err = os.Chmod("operator-key.pem", 0600)
	if err != nil {
		return fmt.Errorf("failed to set appropriate permissions: %v", err)
	}

	return nil
}

func main() {

	// Command-line options.
	publicListenAddrFlag := flag.String("publicAddr", "0.0.0.0:26345", "Specify HTTPS address for receiving experiment instructions.")
	internalListenAddrFlag := flag.String("internalAddr", "0.0.0.0:33000", "Specify HTTPS address for administrating experiments.")
	gcloudBucketFlag := flag.String("gcloudBucket", "", "Supply the GCloud Storage Bucket to use for the experiments.")

	flag.Parse()

	if *gcloudBucketFlag == "" {
		fmt.Printf("Missing argument, please provide a value for flag: '-gcloudBucket'.\n")
		os.Exit(1)
	}

	op := &Operator{
		GCloudBucket: *gcloudBucketFlag,

		InternalListenAddr: *internalListenAddrFlag,
		InternalChan:       make(chan string),

		PublicListenAddr: *publicListenAddrFlag,
		PublicChan:       make(chan string),

		ExpInProgress: "",
		Exps:          make(map[string]*Exp),
	}

	// Download TLS certificate and key from
	// supplied GCP storage bucket.
	err := op.GetTLSMaterial()
	if err != nil {
		fmt.Printf("Preparing TLS material failed: %v\n", err)
		os.Exit(1)
	}

	// Create goroutine that completely
	// handles experiment procedure.
	go op.RunExperiments()

	// Prepare and listen for API calls on the
	// internal network endpoint (worker nodes).
	op.PrepareInternalSrv()
	go op.RunInternalSrv()

	// Prepare and listen for API calls on the
	// Internet-facing endpoint (start experiments).
	op.PreparePublicSrv()

	fmt.Printf("[PUBLIC] Listening on https://%s/public/experiments for API calls regarding experiments...\n", op.PublicListenAddr)

	err = http.ListenAndServeTLS(op.PublicListenAddr, "operator-cert.pem", "operator-key.pem", nil)
	if err != nil {
		fmt.Printf("Failed handling public experiment requests: %v\n", err)
		os.Exit(1)
	}
}
