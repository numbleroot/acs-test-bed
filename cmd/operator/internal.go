package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/emicklei/go-restful"
)

// Worker describes one compute instance
// exhaustively for reproducibility.
type Worker struct {
	Name                   string `json:"name"`
	Status                 string `json:"status"`
	Partner                string `json:"partner"`
	Zone                   string `json:"zone"`
	MinCPUPlatform         string `json:"minCPUPlatform"`
	MachineType            string `json:"machineType"`
	TypeOfNode             string `json:"typeOfNode"`
	EvaluationScript       string `json:"evaluationScript"`
	BinaryName             string `json:"binaryName"`
	SourceImage            string `json:"sourceImage"`
	DiskType               string `json:"diskType"`
	DiskSize               string `json:"diskSize"`
	NetTroublesIfApplied   string `json:"netTroublesIfApplied"`
	ZenoMixKilledIfApplied string `json:"zenoMixKilledIfApplied"`
}

// RegisterResp contains all experiment setup
// instructions relevant for a machine that
// has booted and attempts to complete setup.
// Sent as response to register call.
type RegisterResp struct {
	Partner              string `json:"partner"`
	TypeOfNode           string `json:"typeOfNode"`
	ResultFolder         string `json:"resultFolder"`
	PKI                  string `json:"pki"`
	EvaluationScript     string `json:"evaluationScript"`
	BinaryName           string `json:"binaryName"`
	TCConfig             string `json:"tcConfig"`
	KillZenoMixesInRound string `json:"killZenoMixesInRound"`
}

// HandlerPutRegister accepts a newly booted
// machine as a new worker for the specified
// experiment to be conducted.
func (op *Operator) HandlerPutRegister(req *restful.Request, resp *restful.Response) {

	exp := req.PathParameter("expID")
	workerName := req.PathParameter("worker")

	fmt.Printf("\n[PUT /experiments/%s/workers/%s/register] Handling registration intent.\n", exp, workerName)

	regResp := &RegisterResp{}

	// Figure out whether calling node is a server
	// or a client in the current experiment.
	_, foundAsServer := op.Exps[exp].ServersSpawned[workerName]
	if foundAsServer {

		op.Exps[exp].ServersSpawned[workerName].Status = "registered"

		regResp.Partner = op.Exps[exp].ServersSpawned[workerName].Partner
		regResp.TypeOfNode = op.Exps[exp].ServersSpawned[workerName].TypeOfNode
		regResp.ResultFolder = op.Exps[exp].ResultFolder
		regResp.EvaluationScript = op.Exps[exp].ServersSpawned[workerName].EvaluationScript
		regResp.BinaryName = op.Exps[exp].ServersSpawned[workerName].BinaryName

		regResp.PKI = "TODO"
		regResp.TCConfig = "TODO"
		regResp.KillZenoMixesInRound = "TODO"
	}

	_, foundAsClient := op.Exps[exp].ClientsSpawned[workerName]
	if foundAsClient {

		op.Exps[exp].ClientsSpawned[workerName].Status = "registered"

		regResp.Partner = op.Exps[exp].ClientsSpawned[workerName].Partner
		regResp.TypeOfNode = op.Exps[exp].ClientsSpawned[workerName].TypeOfNode
		regResp.ResultFolder = op.Exps[exp].ResultFolder
		regResp.EvaluationScript = op.Exps[exp].ClientsSpawned[workerName].EvaluationScript
		regResp.BinaryName = op.Exps[exp].ClientsSpawned[workerName].BinaryName

		regResp.PKI = "TODO"
		regResp.TCConfig = "TODO"
		regResp.KillZenoMixesInRound = "TODO"
	}

	if !foundAsServer && !foundAsClient {
		fmt.Printf("\n[PUT /experiments/%s/workers/%s/register] Neither server nor client: %s.\n", exp, workerName, workerName)
		resp.WriteErrorString(http.StatusInternalServerError, "Neither server nor client.")
		return
	}

	fmt.Printf("\n[PUT /experiments/%s/workers/%s/register] Registration successful.\n", exp, workerName)

	// Respond to worker node.
	resp.WriteHeaderAndEntity(http.StatusOK, regResp)
}

// HandlerGetReady marks an previously registered
// worker as prepared for experiment execution.
// The worker has finished all initialization
// steps before calling this endpoint.
func (op *Operator) HandlerGetReady(req *restful.Request, resp *restful.Response) {}

// HandlerPutFinished signals the operator that
// the specified worker has completed all actions
// designated for it in the running experiment.
func (op *Operator) HandlerPutFinished(req *restful.Request, resp *restful.Response) {}

// PrepareInternalSrv initializes all API-related
// things in order to expose an internal-facing
// API endpoint for conducting experiments.
func (op *Operator) PrepareInternalSrv() {

	op.InternalSrv = new(restful.WebService)

	op.InternalSrv.Path("/internal/experiments").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	op.InternalSrv.Route(op.InternalSrv.PUT("/{expID}/workers/{worker}/register").
		To(op.HandlerPutRegister))

	op.InternalSrv.Route(op.InternalSrv.GET("/{expID}/workers/{worker}/ready").
		To(op.HandlerGetReady))

	op.InternalSrv.Route(op.InternalSrv.PUT("/{expID}/workers/{worker}/finished").
		To(op.HandlerPutFinished))

	restful.Add(op.InternalSrv)
}

// RunInternalSrv starts the TLS endpoint for
// internal experiment requests.
func (op *Operator) RunInternalSrv() {

	fmt.Printf("[INTERNAL] Listening on https://%s/internal/experiments for API calls regarding experiments...\n", op.InternalListenAddr)

	err := http.ListenAndServeTLS(op.InternalListenAddr, "operator-cert.pem", "operator-key.pem", nil)
	if err != nil {
		fmt.Printf("Failed handling internal experiment requests: %v\n", err)
		os.Exit(1)
	}
}
