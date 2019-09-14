package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/emicklei/go-restful"
	uuid "github.com/satori/go.uuid"
)

// PublicAuth augments all routes by requiring an
// Authorization header in order for continuation.
func (op *Operator) PublicAuth(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {

	authHeader := req.HeaderParameter("Authorization")

	// Expect correct (static) password.
	if authHeader != "UniverseOfLoopholes" {
		resp.WriteError(http.StatusUnauthorized, nil)
		return
	}

	// Possibly move to next filter.
	chain.ProcessFilter(req, resp)
}

// HandlerPutNew creates a new experiment, if possible.
func (op *Operator) HandlerPutNew(req *restful.Request, resp *restful.Response) {

	op.Lock()
	defer op.Unlock()

	fmt.Printf("\n[PUT /experiments/new] Handling new request from %s\n", req.Request.RemoteAddr)

	// If an experiment is already running,
	// no further one is allowed to run.
	if op.ExpInProgress != "" {
		resp.WriteErrorString(http.StatusConflict, "An experiment is already being conducted at the moment.")
		return
	}

	exp := &Exp{}

	// Read values regarding evaluation system
	// and result folder from received request.
	err := req.ReadEntity(exp)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, nil)
		return
	}

	// Generate new random UUID.
	id, err := uuid.NewV4()
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, nil)
		return
	}

	// Fill in or overwrite remaining
	// fields of experiment struct.
	exp.ID = id.String()
	exp.InitTime = time.Now()
	exp.Concluded = false
	exp.Progress = make([]string, 0, 50)

	// Add experiment to map of all experiments.
	op.Exps[exp.ID] = exp

	// Signal goroutine conducting the
	// experiments availability of the new one.
	// TODO: add.

	fmt.Printf("[PUT /experiments/new] Successfully added new experiment from %s\n", req.Request.RemoteAddr)

	// Send experiment information up to this
	// point back to client.
	resp.WriteHeaderAndEntity(http.StatusCreated, exp)
}

// HandlerGetExpStatus returns the current
// state of an ongoing experiment.
func (op *Operator) HandlerGetExpStatus(req *restful.Request, resp *restful.Response) {
	io.WriteString(resp, "test\n")
}

// PreparePublicSrv initializes all API-related
// things in order to expose an Internet-facing
// API endpoint for conducting experiments.
func (op *Operator) PreparePublicSrv() {

	op.PublicSrv = new(restful.WebService)

	op.PublicSrv.Path("/public/experiments").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	op.PublicSrv.Route(op.PublicSrv.PUT("/new").
		Filter(op.PublicAuth).
		To(op.HandlerPutNew))

	op.PublicSrv.Route(op.PublicSrv.GET("/{expID}/status").
		Filter(op.PublicAuth).
		To(op.HandlerGetExpStatus))

	restful.Add(op.PublicSrv)
}
