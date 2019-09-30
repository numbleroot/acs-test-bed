package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/numbleroot/acs-test-bed/cmd/operator/zenopki"
)

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
				"key": "operatorIP",
				"value": "ACS_EVAL_INSERT_META_OPERATOR_IP"
			},
			{
				"key": "expID",
				"value": "ACS_EVAL_INSERT_META_EXP_ID"
			},
			{
				"key": "nameOfNode",
				"value": "ACS_EVAL_INSERT_META_NAME_OF_NODE"
			},
			{
				"key": "evalSystem",
				"value": "ACS_EVAL_INSERT_META_EVAL_SYSTEM"
			},
			{
				"key": "numClients",
				"value": "ACS_EVAL_INSERT_META_NUM_CLIENTS"
			},
			{
				"key": "resultFolder",
				"value": "ACS_EVAL_INSERT_META_RESULT_FOLDER"
			},
			{
				"key": "typeOfNode",
				"value": "ACS_EVAL_INSERT_META_TYPE_OF_NODE"
			},
			{
				"key": "binaryToPull",
				"value": "ACS_EVAL_INSERT_META_BINARY_TO_PULL"
			},
			{
				"key": "pungServerIP",
				"value": "ACS_EVAL_INSERT_META_PUNG_SERVER_IP"
			},
			{
				"key": "tcConfig",
				"value": "ACS_EVAL_INSERT_META_TC_CONFIG"
			},
			{
				"key": "killZenoMixesInRound",
				"value": "ACS_EVAL_INSERT_META_KILL_ZENO_MIXES_IN_ROUND"
			},
			{
				"key": "client01",
				"value": "ACS_EVAL_INSERT_META_CLIENT_01_NAME"
			},
			{
				"key": "partner01",
				"value": "ACS_EVAL_INSERT_META_CLIENT_01_PARTNER"
			},
			{
				"key": "client02",
				"value": "ACS_EVAL_INSERT_META_CLIENT_02_NAME"
			},
			{
				"key": "partner02",
				"value": "ACS_EVAL_INSERT_META_CLIENT_02_PARTNER"
			},
			{
				"key": "client03",
				"value": "ACS_EVAL_INSERT_META_CLIENT_03_NAME"
			},
			{
				"key": "partner03",
				"value": "ACS_EVAL_INSERT_META_CLIENT_03_PARTNER"
			},
			{
				"key": "client04",
				"value": "ACS_EVAL_INSERT_META_CLIENT_04_NAME"
			},
			{
				"key": "partner04",
				"value": "ACS_EVAL_INSERT_META_CLIENT_04_PARTNER"
			},
			{
				"key": "client05",
				"value": "ACS_EVAL_INSERT_META_CLIENT_05_NAME"
			},
			{
				"key": "partner05",
				"value": "ACS_EVAL_INSERT_META_CLIENT_05_PARTNER"
			},
			{
				"key": "client06",
				"value": "ACS_EVAL_INSERT_META_CLIENT_06_NAME"
			},
			{
				"key": "partner06",
				"value": "ACS_EVAL_INSERT_META_CLIENT_06_PARTNER"
			},
			{
				"key": "client07",
				"value": "ACS_EVAL_INSERT_META_CLIENT_07_NAME"
			},
			{
				"key": "partner07",
				"value": "ACS_EVAL_INSERT_META_CLIENT_07_PARTNER"
			},
			{
				"key": "client08",
				"value": "ACS_EVAL_INSERT_META_CLIENT_08_NAME"
			},
			{
				"key": "partner08",
				"value": "ACS_EVAL_INSERT_META_CLIENT_08_PARTNER"
			},
			{
				"key": "client09",
				"value": "ACS_EVAL_INSERT_META_CLIENT_09_NAME"
			},
			{
				"key": "partner09",
				"value": "ACS_EVAL_INSERT_META_CLIENT_09_PARTNER"
			},
			{
				"key": "client10",
				"value": "ACS_EVAL_INSERT_META_CLIENT_10_NAME"
			},
			{
				"key": "partner10",
				"value": "ACS_EVAL_INSERT_META_CLIENT_10_PARTNER"
			},
			{
				"key": "startup-script-url",
				"value": "ACS_EVAL_INSERT_META_STARTUP_SCRIPT"
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
	"reservationAffinity": {
		"consumeReservationType": "NO_RESERVATION"
	},
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
func (op *Operator) SpawnInstance(exp *Exp, worker *Worker, publiclyReachable bool) {

	clientIDs := make(map[int]string)

	// Calculate start and end IDs for
	// this client machine to handle.
	lastClient := worker.ID * 10
	firstClient := lastClient - 10
	totalClients := len(exp.Clients) * 10

	if worker.TypeOfNode == "client" {

		for i := 1; i <= 10; i++ {
			clientIDs[i] = fmt.Sprintf("client-%05d", (firstClient + i))
		}

	} else {

		// In case of servers, only one
		// logical node will be spawned,
		// the server itself.
		clientIDs[1] = worker.Name

		for i := 2; i <= 10; i++ {
			clientIDs[i] = "irrelevant"
		}
	}

	exp.ProgressChan <- fmt.Sprintf("Spawning %s.", worker.Name)

	// Customize API endpoint to send request to.
	endpoint := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/zones/%s/instances", op.GCloudProject, worker.Zone)

	// Prepare request body.

	reqBody := strings.ReplaceAll(tmplInstanceCreate, "ACS_EVAL_INSERT_GCP_MACHINE_NAME", fmt.Sprintf("%s-%s", worker.Name, exp.ResultFolder))
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_ZONE", fmt.Sprintf("projects/%s/zones/%s", op.GCloudProject, worker.Zone))
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_MIN_CPU_PLATFORM", worker.MinCPUPlatform)
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_MACHINE_TYPE", fmt.Sprintf("projects/%s/zones/%s/machineTypes/%s", op.GCloudProject,
		worker.Zone, worker.MachineType))

	// Replace placeholders in metadata.
	// These are used by the startup script.

	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_META_OPERATOR_IP", strings.Split(op.InternalListenAddr, ":")[0])
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_META_EXP_ID", exp.ID)
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_META_NAME_OF_NODE", worker.Name)
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_META_EVAL_SYSTEM", exp.System)
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_META_NUM_CLIENTS", fmt.Sprintf("%d", totalClients))
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_META_RESULT_FOLDER", exp.ResultFolder)
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_META_TYPE_OF_NODE", worker.TypeOfNode)
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_META_BINARY_TO_PULL", worker.BinaryName)

	if (exp.System == "pung") && (worker.TypeOfNode == "client") {
		reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_META_PUNG_SERVER_IP", exp.ServersMap["server-00001"].Address)
	} else {
		reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_META_PUNG_SERVER_IP", "irrelevant")
	}

	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_META_TC_CONFIG", worker.NetTroubles)
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_META_KILL_ZENO_MIXES_IN_ROUND", fmt.Sprintf("%d", worker.ZenoMixesKilled))
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_META_CLIENT_01_NAME", clientIDs[1])
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_META_CLIENT_01_PARTNER", clientIDs[2])
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_META_CLIENT_02_NAME", clientIDs[2])
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_META_CLIENT_02_PARTNER", clientIDs[1])
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_META_CLIENT_03_NAME", clientIDs[3])
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_META_CLIENT_03_PARTNER", clientIDs[4])
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_META_CLIENT_04_NAME", clientIDs[4])
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_META_CLIENT_04_PARTNER", clientIDs[3])
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_META_CLIENT_05_NAME", clientIDs[5])
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_META_CLIENT_05_PARTNER", clientIDs[6])
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_META_CLIENT_06_NAME", clientIDs[6])
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_META_CLIENT_06_PARTNER", clientIDs[5])
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_META_CLIENT_07_NAME", clientIDs[7])
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_META_CLIENT_07_PARTNER", clientIDs[8])
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_META_CLIENT_08_NAME", clientIDs[8])
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_META_CLIENT_08_PARTNER", clientIDs[7])
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_META_CLIENT_09_NAME", clientIDs[9])
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_META_CLIENT_09_PARTNER", clientIDs[10])
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_META_CLIENT_10_NAME", clientIDs[10])
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_META_CLIENT_10_PARTNER", clientIDs[9])
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_META_STARTUP_SCRIPT", fmt.Sprintf("gs://%s/startup.sh", op.GCloudBucket))

	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_SOURCE_IMAGE", fmt.Sprintf("projects/%s/global/images/%s", op.GCloudProject,
		worker.SourceImage))
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_DISK_TYPE", fmt.Sprintf("projects/%s/zones/%s/diskTypes/%s", op.GCloudProject,
		worker.Zone, worker.DiskType))
	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_DISK_SIZE", worker.DiskSize)

	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_SUBNETWORK", fmt.Sprintf("projects/%s/regions/%s/subnetworks/default",
		op.GCloudProject, strings.TrimSuffix(worker.Zone, "-b")))

	if publiclyReachable {
		reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_ACCESS_CONFIG", tmplInstancePublicIP)
	} else {
		reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_ACCESS_CONFIG", "")
	}

	reqBody = strings.ReplaceAll(reqBody, "ACS_EVAL_INSERT_SERVICE_ACCOUNT", op.GCloudServiceAcc)

	// Create HTTP POST request.
	request, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(reqBody))
	if err != nil {
		exp.ProgressChan <- fmt.Sprintf("Failed creating HTTP API request to spawn %s: %v", worker.Name, err)
		os.Exit(1)
	}
	request.Header.Set(http.CanonicalHeaderKey("authorization"), fmt.Sprintf("Bearer %s", op.GCloudAccessToken))
	request.Header.Set(http.CanonicalHeaderKey("content-type"), "application/json")

	// Send the request to GCP.
	tried := 0
	resp, err := http.DefaultClient.Do(request)
	for err != nil && tried < 10 {

		tried++
		exp.ProgressChan <- fmt.Sprintf("Create API request for %s failed (will try again): %v", worker.Name, err)
		time.Sleep(1 * time.Second)

		resp, err = http.DefaultClient.Do(request)
	}

	if tried >= 10 {
		exp.ProgressChan <- fmt.Sprintf("Create API request for %s failed permanently: %v", worker.Name, err)
		os.Exit(1)
	}

	// Read the response.
	outRaw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		exp.ProgressChan <- fmt.Sprintf("Failed reading from instance create response body for %s: %v", worker.Name, err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	out := string(outRaw)

	// Verify successful machine creation.
	if strings.Contains(out, "RUNNING") || strings.Contains(out, "DONE") {
		exp.ProgressChan <- fmt.Sprintf("Instance %s running, waiting for initialization to finish.", worker.Name)
	} else {
		exp.ProgressChan <- fmt.Sprintf("Spawning instance %s returned failure message: '%s'", worker.Name, out)
		os.Exit(1)
	}
}

// ShutdownInstance instructs GCP to shut down and
// subsequently delete an instance.
func (op *Operator) ShutdownInstance(wg *sync.WaitGroup, exp *Exp, worker *Worker) {

	defer wg.Done()

	exp.ProgressChan <- fmt.Sprintf("Deleting %s.", worker.Name)

	instName := fmt.Sprintf("%s-%s", worker.Name, exp.ResultFolder)
	endpoint := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/zones/%s/instances/%s", op.GCloudProject, worker.Zone, instName)

	// Create HTTP DELETE request.
	request, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		exp.ProgressChan <- fmt.Sprintf("Failed creating HTTP API request for %s: %v", worker.Name, err)
		os.Exit(1)
	}
	request.Header.Set(http.CanonicalHeaderKey("authorization"), fmt.Sprintf("Bearer %s", op.GCloudAccessToken))

	// Send the request to GCP.
	resp, err := http.DefaultClient.Do(request)
	for err != nil {

		exp.ProgressChan <- fmt.Sprintf("Delete API request for %s failed (will try again): %v", worker.Name, err)
		time.Sleep(1 * time.Second)

		resp, err = http.DefaultClient.Do(request)
	}

	// Read the response.
	outRaw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		exp.ProgressChan <- fmt.Sprintf("Failed reading from delete response body for %s: %v", worker.Name, err)
		return
	}
	defer resp.Body.Close()
	out := string(outRaw)

	// Verify successful machine deletion.
	if strings.Contains(out, "RUNNING") || strings.Contains(out, "DONE") {
		exp.ProgressChan <- fmt.Sprintf("Successfully deleted %s.", worker.Name)
	} else {
		exp.ProgressChan <- fmt.Sprintf("Failed deleting %s.", worker.Name)
	}
}

// ProgressWriter is the only routine allowed to append
// time-stamped log lines to the progress slice of its
// experiment. New lines are sent to it via the experiment's
// progress channel.
func (exp *Exp) ProgressWriter() {

	for line := range exp.ProgressChan {
		exp.Progress = append(exp.Progress, fmt.Sprintf("[%s] %s", time.Now().Format("2006-02-03 15:04:05"), line))
		fmt.Printf(fmt.Sprintf("[%s] %s\n", time.Now().Format("2006-02-03 15:04:05"), line))
	}
}

// VuvuzelaProducePKI writes out the collected server
// addresses into the otherwise prepared pki.conf file
// that all Vuvuzela nodes use instead of an actual
// PKI node.
func (exp *Exp) VuvuzelaProducePKI() error {

	// Read preliminary PKI file into memory.
	pki, err := ioutil.ReadFile("/root/vuvuzela-confs/pki.conf")
	if err != nil {
		return err
	}

	// Replace all placeholders for a server's
	// address with its registered address in
	// this experiment.
	for i := range exp.Servers {

		if exp.Servers[i].ID == 1 {
			pki = bytes.ReplaceAll(pki, []byte(fmt.Sprintf("ACS_EVAL_INSERT_%s_ADDRESS", exp.Servers[i].Name)), []byte(fmt.Sprintf("ws://%s", exp.Servers[i].Address)))
		} else {
			pki = bytes.ReplaceAll(pki, []byte(fmt.Sprintf("ACS_EVAL_INSERT_%s_ADDRESS", exp.Servers[i].Name)), []byte(exp.Servers[i].Address))
		}
	}

	// Write out final pki.conf.
	err = ioutil.WriteFile("/root/vuvuzela-confs/pki.conf", pki, 0644)
	if err != nil {
		return err
	}

	// Upload pki.conf to GCloud bucket.
	out, err := exec.Command("/usr/bin/gsutil", "cp", "/root/vuvuzela-confs/pki.conf", "gs://acs-eval/vuvuzela-confs/pki.conf").CombinedOutput()
	if err != nil {
		return err
	}

	if !bytes.Contains(out, []byte("completed")) {
		return fmt.Errorf("uploading final pki.conf to GCloud bucket unsuccessful")
	}

	return nil
}

// RunExperiments is the authoritative goroutine
// for provisioning machines and conducting all
// experiments queued in the Operator.
func (op *Operator) RunExperiments() {

	for {

		// Wait for new experiment to be queued.
		expID := <-op.PublicNewChan

		op.Lock()

		// Mark this experiment as in progress.
		op.ExpInProgress = expID

		// Retrieve experiment data.
		exp := op.Exps[expID]

		op.Unlock()

		// Launch progress writing routine.
		go exp.ProgressWriter()

		exp.ProgressChan <- fmt.Sprintf("Commencing experiment %s for system %s with %d server and %d client instances.",
			expID, exp.System, len(exp.Servers), len(exp.Clients))

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

			exp.ProgressChan <- fmt.Sprintf("Launching zeno PKI process at %s.", op.ZenoPKI.LisAddr)

			// Run zeno PKI process in background.
			go op.ZenoPKI.Run(op.TLSCertPath, op.TLSKeyPath)
		}

		// Spawn all server machines in reverse order.
		// This hopefully enables Vuvuzela mixes to connect
		// to their successor on first try.
		for i := (len(exp.Servers) - 1); i >= 0; i-- {
			op.SpawnInstance(exp, exp.Servers[i], true)
			time.Sleep(1 * time.Second)
		}

		exp.ProgressChan <- fmt.Sprintf("All %d servers instructed to spawn, waiting for registration requests.", len(exp.Servers))

		// Handle incoming registration requests.
		for range exp.Servers {

			workerReg := <-op.InternalRegisterChan

			_, found := exp.ServersMap[workerReg.Worker]
			if found {
				exp.ServersMap[workerReg.Worker].Address = workerReg.Address
				exp.ServersMap[workerReg.Worker].Status = "registered"
				exp.ProgressChan <- fmt.Sprintf("Server %s at %s marked as registered.", workerReg.Worker, workerReg.Address)
			}

			status := "Servers status:\n"
			for i := range exp.Servers {
				status = fmt.Sprintf("%s\t%s@%s: %s\n", status, exp.Servers[i].Name, exp.Servers[i].Address, exp.Servers[i].Status)
			}

			exp.ProgressChan <- status
		}

		if exp.System == "vuvuzela" {

			// If Vuvuzela is being evaluated, we need to
			// quickly produce an appropriate pki.conf file.
			err := exp.VuvuzelaProducePKI()
			if err != nil {
				exp.ProgressChan <- fmt.Sprintf("Failed to produce final pki.conf file for Vuvuzela: %v", err)
				os.Exit(1)
			}

			exp.ProgressChan <- "pki.conf for Vuvuzela created and uploaded."
		}

		// Handle incoming ready or failed requests.
		for range exp.Servers {

			select {

			case workerName := <-op.InternalReadyChan:

				_, found := exp.ServersMap[workerName]
				if found {
					exp.ServersMap[workerName].Status = "ready"
					exp.ProgressChan <- fmt.Sprintf("Server %s marked as ready.", workerName)
				}

			case failedReq := <-op.InternalFailedChan:

				_, found := exp.ServersMap[failedReq.Worker]
				if found {
					exp.ServersMap[failedReq.Worker].Status = "failed"
					exp.ProgressChan <- fmt.Sprintf("Server %s failed with: %s", failedReq.Worker, failedReq.Reason)
				}
			}

			status := "Servers status:\n"
			for i := range exp.Servers {
				status = fmt.Sprintf("%s\t%s@%s: %s\n", status, exp.Servers[i].Name, exp.Servers[i].Address, exp.Servers[i].Status)
			}

			exp.ProgressChan <- status
		}

		// Verify all servers ready.
		for i := range exp.Servers {

			if exp.Servers[i].Status != "ready" {
				exp.ProgressChan <- fmt.Sprintf("At least one server (%s) failed to initialize, ending experiment %s.",
					exp.Servers[i].Name, expID)
				goto END
			}
		}

		exp.ProgressChan <- fmt.Sprintf("All %d servers spawned and ready, launching clients.", len(exp.Servers))

		// Spawn all client machines.
		for i := range exp.Clients {
			go op.SpawnInstance(exp, exp.Clients[i], false)
		}

		exp.ProgressChan <- fmt.Sprintf("All %d clients instructed to spawn, waiting for registration requests.", len(exp.Clients))

		// Handle incoming client registration requests.
		for range exp.Clients {

			workerReg := <-op.InternalRegisterChan

			_, found := exp.ClientsMap[workerReg.Worker]
			if found {
				exp.ClientsMap[workerReg.Worker].Status = "registered"
				exp.ProgressChan <- fmt.Sprintf("Client %s marked as registered.", workerReg.Worker)
			}

			status := "Clients status:\n"
			for i := range exp.Clients {
				status = fmt.Sprintf("%s\t%s@%s: %s\n", status, exp.Clients[i].Name, exp.Clients[i].Address, exp.Clients[i].Status)
			}

			exp.ProgressChan <- status
		}

		// Handle incoming ready or failed requests.
		for range exp.Clients {

			select {

			case workerName := <-op.InternalReadyChan:

				_, found := exp.ClientsMap[workerName]
				if found {
					exp.ClientsMap[workerName].Status = "ready"
					exp.ProgressChan <- fmt.Sprintf("Client %s marked as ready.", workerName)
				}

			case failedReq := <-op.InternalFailedChan:

				_, found := exp.ClientsMap[failedReq.Worker]
				if found {
					exp.ClientsMap[failedReq.Worker].Status = "failed"
					exp.ProgressChan <- fmt.Sprintf("Client %s failed with: %s", failedReq.Worker, failedReq.Reason)
				}
			}

			status := "Clients status:\n"
			for i := range exp.Clients {
				status = fmt.Sprintf("%s\t%s@%s: %s\n", status, exp.Clients[i].Name, exp.Clients[i].Address, exp.Clients[i].Status)
			}

			exp.ProgressChan <- status
		}

		// Verify all clients ready.
		for i := range exp.Clients {

			if exp.Clients[i].Status != "ready" {
				exp.ProgressChan <- fmt.Sprintf("At least one client (%s) failed to initialize, ending experiment %s.",
					exp.Clients[i].Name, expID)
				goto END
			}
		}

		if exp.System == "zeno" {

			// If ACS under evaluation is zeno, signal
			// PKI routine to start broadcasting.
			zenoEvalCtrlChan <- struct{}{}
		}

		// Wait for all clients to signal completion.
		for range exp.Clients {

			workerName := <-op.InternalFinishedChan

			_, found := exp.ClientsMap[workerName]
			if found {
				exp.ClientsMap[workerName].Status = "finished"
				exp.ProgressChan <- fmt.Sprintf("Client %s marked as finished.", workerName)
			}

			status := "Clients status:\n"
			for i := range exp.Clients {
				status = fmt.Sprintf("%s\t%s@%s: %s\n", status, exp.Clients[i].Name, exp.Clients[i].Address, exp.Clients[i].Status)
			}

			exp.ProgressChan <- status
		}

		// Verify all clients completed.
		for i := range exp.Clients {

			if exp.Clients[i].Status != "finished" {
				exp.ProgressChan <- fmt.Sprintf("At least one client (%s) did not finish in experiment %s.",
					exp.Clients[i].Name, expID)
			}
		}

	END:
		exp.ProgressChan <- fmt.Sprintf("Experiment %s reached end, awaiting shutdown confirmation.", expID)

		// Wait for explicit shutdown confirmation.
		<-op.PublicTerminateChan

		exp.ProgressChan <- fmt.Sprintf("Shutdown confirmation for experiment %s received.", expID)

		// Shut down all client machines.
		wg := &sync.WaitGroup{}

		for i := 0; i < len(exp.Clients); i++ {
			wg.Add(1)
			go op.ShutdownInstance(wg, exp, exp.Clients[i])
		}

		wg.Wait()

		exp.ProgressChan <- fmt.Sprintf("All %d clients successfully shut down.", len(exp.Clients))

		// Shut down all server machines.
		wg = &sync.WaitGroup{}

		for i := 0; i < len(exp.Servers); i++ {
			wg.Add(1)
			go op.ShutdownInstance(wg, exp, exp.Servers[i])
		}

		wg.Wait()

		exp.ProgressChan <- fmt.Sprintf("All %d servers successfully shut down.", len(exp.Servers))

		// Mark experiment as done.
		exp.Concluded = true
		close(exp.ProgressChan)

		op.Lock()

		// Reset in-progress indicator.
		op.ExpInProgress = ""

		// Overwrite experiment values in Operator
		// struct in case they did not update.
		op.Exps[expID] = exp

		op.Unlock()
	}
}
