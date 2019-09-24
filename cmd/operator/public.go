package main

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"time"

	"github.com/emicklei/go-restful"
)

// ExpReq is supplied by the client talking to the
// public endpoint of the operator and specifies the
// execution details of one experiment in full.
type ExpReq struct {
	System       string    `json:"system"`
	ResultFolder string    `json:"resultFolder"`
	Servers      []*Worker `json:"servers"`
	Clients      []*Worker `json:"clients"`
}

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

	fmt.Printf("\n[PUT /experiments/new] Handling new request from %s.\n", req.Request.RemoteAddr)

	op.Lock()
	inProg := op.ExpInProgress
	op.Unlock()

	// If an experiment is already running,
	// no further one is allowed to run.
	if inProg != "" {
		resp.WriteErrorString(http.StatusConflict, "An experiment is already being conducted at the moment.")
		return
	}

	exp := &Exp{}
	expReq := &ExpReq{}

	// Extract experiment details from request.
	err := req.ReadEntity(expReq)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, nil)
		return
	}

	// Generate new random id.
	id := make([]byte, 8)
	_, err = rand.Read(id)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, nil)
		return
	}

	// Fill in complete experiment specification.
	exp.ID = fmt.Sprintf("%x", id)
	exp.Created = time.Now().Format("2006-02-03_15:04:05")
	exp.System = expReq.System
	exp.Concluded = false
	exp.ResultFolder = expReq.ResultFolder
	exp.Progress = make([]string, 0, 50)
	exp.ProgressChan = make(chan string)
	exp.Servers = make([]*Worker, len(expReq.Servers))
	exp.ServersMap = make(map[string]*Worker)
	exp.Clients = make([]*Worker, len(expReq.Clients))
	exp.ClientsMap = make(map[string]*Worker)

	for i := range expReq.Servers {
		exp.Servers[i] = expReq.Servers[i]
		exp.ServersMap[expReq.Servers[i].Name] = expReq.Servers[i]
	}

	for i := range expReq.Clients {
		exp.Clients[i] = expReq.Clients[i]
		exp.ClientsMap[expReq.Clients[i].Name] = expReq.Clients[i]
	}

	// Add experiment to map of all experiments.
	op.Exps[exp.ID] = exp

	// Signal goroutine conducting the
	// experiments availability of the new one.
	op.PublicNewChan <- exp.ID

	fmt.Printf("[PUT /experiments/new] Successfully added new experiment %s from %s.\n", exp.ID, req.Request.RemoteAddr)

	// Send experiment information up to this
	// point back to client.
	resp.WriteHeaderAndEntity(http.StatusCreated, exp)
}

// HandlerGetExpStatus returns the
// state of the specified experiment.
func (op *Operator) HandlerGetExpStatus(req *restful.Request, resp *restful.Response) {

	expID := req.PathParameter("expID")

	fmt.Printf("\n[GET /experiments/%s/status] Returning experiment status to %s.\n", expID, req.Request.RemoteAddr)

	// If experiment exists, return its status.
	exp, found := op.Exps[expID]
	if !found {
		resp.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("Experiment %s does not exist.", expID))
		return
	}

	resp.WriteHeaderAndEntity(http.StatusOK, exp)
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
