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

// HandlerPutRegister accepts a newly booted
// machine as a new worker for the specified
// experiment to be conducted.
func (op *Operator) HandlerPutRegister(req *restful.Request, resp *restful.Response) {

	expID := req.PathParameter("expID")
	workerName := req.PathParameter("worker")

	fmt.Printf("\n[PUT /experiments/%s/workers/%s/register] Handling registration intent.\n", expID, workerName)

	// Signal runner which worker intends to register.
	op.InternalRegisterChan <- workerName

	fmt.Printf("[PUT /experiments/%s/workers/%s/register] Registration successful.\n", expID, workerName)

	// Respond to worker node.
	resp.WriteHeader(http.StatusOK)
}

// HandlerPutReady marks an previously registered
// worker as prepared for experiment execution.
// The worker has finished all initialization
// steps before calling this endpoint.
func (op *Operator) HandlerPutReady(req *restful.Request, resp *restful.Response) {

	expID := req.PathParameter("expID")
	workerName := req.PathParameter("worker")

	fmt.Printf("\n[PUT /experiments/%s/workers/%s/ready] Handling ready signal.\n", expID, workerName)

	// Signal runner which worker is ready.
	op.InternalReadyChan <- workerName

	fmt.Printf("[PUT /experiments/%s/workers/%s/ready] Worker marked as ready.\n", expID, workerName)

	// Respond to worker node.
	resp.WriteHeader(http.StatusOK)
}

// HandlerPutFinished signals the operator that
// the specified worker has completed all actions
// designated for it in the running experiment.
func (op *Operator) HandlerPutFinished(req *restful.Request, resp *restful.Response) {

	expID := req.PathParameter("expID")
	workerName := req.PathParameter("worker")

	fmt.Printf("\n[PUT /experiments/%s/workers/%s/finished] Handling finished signal.\n", expID, workerName)

	// Signal runner which worker has finished.
	op.InternalFinishedChan <- workerName

	fmt.Printf("[PUT /experiments/%s/workers/%s/finished] Worker marked as finished.\n", expID, workerName)

	resp.WriteHeader(http.StatusOK)
}

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

	op.InternalSrv.Route(op.InternalSrv.PUT("/{expID}/workers/{worker}/ready").
		To(op.HandlerPutReady))

	op.InternalSrv.Route(op.InternalSrv.PUT("/{expID}/workers/{worker}/finished").
		To(op.HandlerPutFinished))

	restful.Add(op.InternalSrv)
}

// RunInternalSrv starts the TLS endpoint for
// internal experiment requests.
func (op *Operator) RunInternalSrv() {

	fmt.Printf("[INTERNAL] Listening on https://%s/internal/experiments for API calls regarding experiments...\n", op.InternalListenAddr)

	err := http.ListenAndServeTLS(op.InternalListenAddr, op.TLSCertPath, op.TLSKeyPath, nil)
	if err != nil {
		fmt.Printf("Failed handling internal experiment requests: %v\n", err)
		os.Exit(1)
	}
}
