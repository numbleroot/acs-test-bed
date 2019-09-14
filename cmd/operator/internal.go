package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/emicklei/go-restful"
)

// HandlerPutRegister accepts a newly booted
// machine as a new worker for the specified
// experiment to be conducted.
func (op *Operator) HandlerPutRegister(req *restful.Request, resp *restful.Response) {}

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
