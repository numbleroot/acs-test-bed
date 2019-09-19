package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/numbleroot/acs-test-bed/cmd/operator/zenopki"
)

// TODO: Add all necessary metadata.
var tmplInstanceCreate = `{
	"kind": "compute#instance",
	"name": "ACS_EVAL_INSERT_GCP_MACHINE_NAME",
	"zone": "ACS_EVAL_INSERT_ZONE",
	"minCpuPlatform": "ACS_EVAL_INSERT_MIN_CPU_PLATFORM",
	"machineType": "ACS_EVAL_INSERT_MACHINE_TYPE",
	"displayDevice": {
		"enableDisplay": false
	},
	"metadata": {
		"kind": "compute#metadata",
		"items": [
			{
				"key": "nameOfNode",
				"value": "ACS_EVAL_INSERT_NAME"
			},
			{
				"key": "startup-script-url",
				"value": "ACS_EVAL_INSERT_STARTUP_SCRIPT"
			}
		]
	},
	"tags": {
		"items": []
	},
	"disks": [
		{
			"kind": "compute#attachedDisk",
			"type": "PERSISTENT",
			"boot": true,
			"mode": "READ_WRITE",
			"autoDelete": true,
			"deviceName": "ACS_EVAL_INSERT_GCP_MACHINE_NAME",
			"initializeParams": {
				"sourceImage": "ACS_EVAL_INSERT_SOURCE_IMAGE",
				"diskType": "ACS_EVAL_INSERT_DISK_TYPE",
				"diskSizeGb": "ACS_EVAL_INSERT_DISK_SIZE"
			},
			"diskEncryptionKey": {}
		}
	],
	"canIpForward": false,
	"networkInterfaces": [
		{
			"kind": "compute#networkInterface",
			"subnetwork": "ACS_EVAL_INSERT_SUBNETWORK",ACS_EVAL_INSERT_ACCESS_CONFIG
			"aliasIpRanges": []
		}
	],
	"description": "",
	"labels": {},
	"scheduling": {
		"preemptible": false,
		"onHostMaintenance": "TERMINATE",
		"automaticRestart": false,
		"nodeAffinities": []
	},
	"deletionProtection": false,
	"serviceAccounts": [
		{
			"email": "ACS_EVAL_INSERT_SERVICE_ACCOUNT",
			"scopes": [
                "https://www.googleapis.com/auth/compute",
                "https://www.googleapis.com/auth/servicecontrol",
                "https://www.googleapis.com/auth/service.management",
                "https://www.googleapis.com/auth/logging.write",
                "https://www.googleapis.com/auth/monitoring.write",
                "https://www.googleapis.com/auth/trace.append",
                "https://www.googleapis.com/auth/devstorage.full_control"
			]
		}
	]
}`

var tmplInstancePublicIP = `
			"accessConfigs": [
				{
					"kind": "compute#accessConfig",
					"name": "External NAT",
					"type": "ONE_TO_ONE_NAT",
					"networkTier": "PREMIUM"
				}
			],`

// SpawnInstance provisions a compute instance with
// the characteristics from supplied worker struct.
func (op *Operator) SpawnInstance(worker *Worker, resultFolder string, publiclyReachable bool) {

	// Customize API endpoint to send request to.
	endpoint := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/zones/%s/instances", op.GCloudProject, worker.Zone)

	// Prepare request body.
	reqBody := strings.ReplaceAll(tmplInstanceCreate, "ACS_EVAL_INSERT_GCP_MACHINE_NAME", fmt.Sprintf("%s-%s", worker.Name, resultFolder))
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_ZONE", fmt.Sprintf("projects/%s/zones/%s", op.GCloudProject, worker.Zone))
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_MIN_CPU_PLATFORM", worker.MinCPUPlatform)
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_MACHINE_TYPE", fmt.Sprintf("projects/%s/zones/%s/machineTypes/%s", op.GCloudProject,
		worker.Zone, worker.MachineType))
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_NAME", worker.Name)
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_STARTUP_SCRIPT", fmt.Sprintf("gs://%s/startup.sh", op.GCloudBucket))
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_SOURCE_IMAGE", fmt.Sprintf("projects/%s/global/images/%s", op.GCloudProject,
		worker.SourceImage))
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_DISK_TYPE", fmt.Sprintf("projects/%s/zones/%s/diskTypes/%s", op.GCloudProject,
		worker.Zone, worker.DiskType))
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_DISK_SIZE", worker.DiskSize)
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_SUBNETWORK", fmt.Sprintf("projects/%s/regions/%s/subnetworks/default", op.GCloudProject,
		strings.TrimSuffix(worker.Zone, "-b")))
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_SERVICE_ACCOUNT", op.GCloudServiceAcc)

	if publiclyReachable {
		reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_ACCESS_CONFIG", tmplInstancePublicIP)
	} else {
		reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_ACCESS_CONFIG", "")
	}

	// TODO: Replace all metadata placeholders with actual values.
	// TODO: Especially make sure to add PUNG_SERVER_IP.

	// Create HTTP POST request.
	request, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(reqBody))
	if err != nil {
		fmt.Printf("[RUNNER.SPAWNINSTANCE] Failed creating HTTP API request: %v\n", err)
		os.Exit(1)
	}
	request.Header.Set(http.CanonicalHeaderKey("content-type"), "application/json")

	// Send the request to GCP.
	tried := 0
	resp, err := http.DefaultClient.Do(request)
	for err != nil && tried < 10 {

		tried++
		fmt.Printf("[RUNNER.SPAWNINSTANCE] Create API request failed (will try again): %v\n", err)
		time.Sleep(1 * time.Second)

		resp, err = http.DefaultClient.Do(request)
	}

	if tried >= 10 {
		fmt.Printf("[RUNNER.SPAWNINSTANCE] Create API request failed permanently: %v\n", err)
		os.Exit(1)
	}

	// Read the response.
	outRaw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[RUNNER.SPAWNINSTANCE] Failed reading from instance create response body: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	out := string(outRaw)

	// Verify successful machine creation.
	if strings.Contains(out, "RUNNING") || strings.Contains(out, "DONE") {
		fmt.Printf("[RUNNER.SPAWNINSTANCE] Instance %s running, waiting for initialization to finish...\n", worker.Name)
	} else {
		fmt.Printf("[RUNNER.SPAWNINSTANCE] Spawning instance %s returned failure message:\n%s\n", worker.Name, out)
		os.Exit(1)
	}
}

// RunExperiments is the authoritative goroutine
// for provisioning machines and conducting all
// experiments queued in the Operator.
func (op *Operator) RunExperiments() {

	for {

		// Wait for new experiment to be queued.
		expID := <-op.PublicChan

		op.Lock()

		// Mark this experiment as in progress.
		op.ExpInProgress = expID

		// Stash all experiment data.
		exp := op.Exps[expID]

		op.Unlock()

		// Prepare zeno evaluation control channel in
		// case it is needed later on.
		zenoEvalCtrlChan := make(chan struct{})

		if exp.System == "zeno" {

			// If zeno is being evaluated, initialize
			// a PKI struct and have it listen in background.
			op.ZenoPKI = &zenopki.PKI{
				LisAddr:          fmt.Sprintf("%s:44001", strings.Split(op.InternalListenAddr, ":")[0]),
				EvalCtrlChan:     zenoEvalCtrlChan,
				AcceptMixRegs:    0,
				AcceptClientRegs: 0,
				MuNodes:          &sync.RWMutex{},
				Nodes:            make(map[string]*zenopki.Endpoint),
			}

			// Run zeno PKI process in background.
			go op.ZenoPKI.Run(op.TLSCertPath, op.TLSKeyPath)
		}

		// Spawn all server machines.
		for i := range exp.ServersSpawned {
			go op.SpawnInstance(exp.ServersSpawned[i], exp.ResultFolder, true)
		}

		// Handle incoming registration requests.
		for range exp.ServersSpawned {

			workerName := <-op.InternalRegisterChan

			_, found := exp.ServersSpawned[workerName]
			if found {
				exp.ServersSpawned[workerName].Status = "registered"
			}
		}

		// Handle incoming ready or failed requests.
		for range exp.ServersSpawned {

			select {

			case workerName := <-op.InternalReadyChan:

				_, found := exp.ServersSpawned[workerName]
				if found {
					exp.ServersSpawned[workerName].Status = "ready"
				}

			case failedReq := <-op.InternalFailedChan:

				_, found := exp.ServersSpawned[failedReq.Worker]
				if found {
					exp.ServersSpawned[failedReq.Worker].Status = fmt.Sprintf("failed with: '%s'", failedReq.Reason)
				}
			}
		}

		// TODO: Verify enough servers ready.

		// TODO: Spawn all client machines.

		// Handle incoming client registration requests.
		for range exp.ClientsSpawned {

			workerName := <-op.InternalRegisterChan

			_, foundAsClient := exp.ClientsSpawned[workerName]
			if foundAsClient {
				exp.ClientsSpawned[workerName].Status = "registered"
			}
		}

		// TODO: Verify enough clients ready.

		// TODO: Conduct experiment.

		if exp.System == "zeno" {

			// If ACS under evaluation is zeno, signal
			// PKI routine to start broadcasting.
			zenoEvalCtrlChan <- struct{}{}
		}

		// TODO: Wait for all clients to signal completion.

		exp.Concluded = true

		op.Lock()

		// Reset in-progress indicator.
		op.ExpInProgress = ""

		// Overwrite experiment values in Operator
		// struct in case they did not update.
		op.Exps[expID] = exp

		op.Unlock()
	}
}
