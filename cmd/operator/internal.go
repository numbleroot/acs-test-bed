package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/emicklei/go-restful"
)

// RegisterReq transports the address information
// of the first logical node on a worker instance
// to the operator.
type RegisterReq struct {
	Worker  string `json:"-"`
	Address string `json:"address"`
}

// FailedReq captures all information taken
// from a failure signal from a worker.
type FailedReq struct {
	Worker string `json:"-"`
	Reason string `json:"failure"`
}

// HandlerPutRegister accepts a newly booted
// machine as a new worker for the specified
// experiment to be conducted.
func (op *Operator) HandlerPutRegister(req *restful.Request, resp *restful.Response) {

	expID := req.PathParameter("expID")
	workerName := req.PathParameter("worker")

	// Read address information from request.
	regReq := &RegisterReq{}
	err := req.ReadEntity(&regReq)
	if err != nil {
		fmt.Printf("[PUT /experiments/%s/workers/%s/register] Failed to extract payload containing address: %v.\n", expID, workerName, err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	regReq.Worker = workerName

	// Signal runner which worker intends to register.
	op.Exps[expID].RegisterChan <- regReq

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

	// Signal runner which worker is ready.
	op.Exps[expID].ReadyChan <- workerName

	// Respond to worker node.
	resp.WriteHeader(http.StatusOK)
}

// HandlerPutFinished signals the operator that
// the specified worker has completed all actions
// designated for it in the running experiment.
func (op *Operator) HandlerPutFinished(req *restful.Request, resp *restful.Response) {

	expID := req.PathParameter("expID")
	workerName := req.PathParameter("worker")

	// Signal runner which worker has finished.
	op.Exps[expID].FinishedChan <- workerName

	resp.WriteHeader(http.StatusOK)
}

// HandlerPutFailed sends a failure signal and
// message from a failed worker to the operator.
func (op *Operator) HandlerPutFailed(req *restful.Request, resp *restful.Response) {

	expID := req.PathParameter("expID")
	workerName := req.PathParameter("worker")

	// Read failure information from request.
	failedReq := &FailedReq{}
	err := req.ReadEntity(&failedReq)
	if err != nil {
		fmt.Printf("[PUT /experiments/%s/workers/%s/failed] Failed to extract payload containing error message: %v.\n", expID, workerName, err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	failedReq.Worker = workerName

	// Signal runner which worker has failed.
	op.Exps[expID].FailedChan <- failedReq

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

	op.InternalSrv.Route(op.InternalSrv.PUT("/{expID}/workers/{worker}/failed").
		To(op.HandlerPutFailed))

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
