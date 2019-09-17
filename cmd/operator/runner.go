package main

import (
	"strings"
	"sync"

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
        "key": "nameOfNode",
        "value": "ACS_EVAL_INSERT_NAME"
      },
      {
        "key": "partnerOfNode",
        "value": "ACS_EVAL_INSERT_PARTNER"
      },
      {
        "key": "typeOfNode",
        "value": "ACS_EVAL_INSERT_TYPE_OF_NODE"
      },
      {
        "key": "resultFolder",
        "value": "ACS_EVAL_INSERT_RESULT_FOLDER"
      },
      {
        "key": "evalScriptToPull",
        "value": "ACS_EVAL_INSERT_EVAL_SCRIPT_TO_PULL"
      },
      {
        "key": "binaryToPull",
        "value": "ACS_EVAL_INSERT_BINARY_TO_PULL"
      },
      {
        "key": "tcConfig",
        "value": "ACS_EVAL_INSERT_TC_CONFIG"
      },
      {
        "key": "killZenoMixesInRound",
        "value": "ACS_EVAL_INSERT_KILL_ZENO_MIXES_IN_ROUND"
      },
      {
        "key": "pkiIP",
        "value": "ACS_EVAL_INSERT_PKI_IP"
      },
      {
        "key": "startup-script-url",
        "value": "ACS_EVAL_INSERT_STARTUP_SCRIPT"
      }
    ]
  },
  "tags": {
    "items": [
      "ACS_EVAL_INSERT_TAG"
    ]
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
        "https://www.googleapis.com/auth/servicecontrol",
        "https://www.googleapis.com/auth/service.management.readonly",
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
				LisAddr:          strings.Split(op.InternalListenAddr, ":")[0],
				EvalCtrlChan:     zenoEvalCtrlChan,
				AcceptMixRegs:    0,
				AcceptClientRegs: 0,
				MuNodes:          &sync.RWMutex{},
				Nodes:            make(map[string]*zenopki.Endpoint),
			}

			// Run zeno PKI process in background.
			go op.ZenoPKI.Run("operator-cert.pem", "operator-key.pem")
		}

		// TODO: Spawn all server machines.

		// TODO: Verify enough servers ready.

		// TODO: Spawn all client machines.

		// TODO: Verify enough clients ready.

		// TODO: Conduct experiment.

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
