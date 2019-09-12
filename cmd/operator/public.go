package main

import (
	"io"

	"github.com/emicklei/go-restful"
)

// HandlerGetNew creates a new experiment, if possible.
func (op *Operator) HandlerGetNew(req *restful.Request, resp *restful.Response) {
	io.WriteString(resp, "world\n")
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

	op.PublicSrv.Path("/experiments").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	op.PublicSrv.Route(op.PublicSrv.GET("/new/{systemID}").
		To(op.HandlerGetNew).
		Doc("Trigger the start of a new experiment, if possible.").
		Param(op.PublicSrv.PathParameter("systemID", "Identifier of the ACS to evaluate.").DataType("string")).
		Writes(Exp{}))
	op.PublicSrv.Route(op.PublicSrv.GET("/{expID}/status").
		To(op.HandlerGetExpStatus).
		Doc("Return the current state of an ongoing experiment.").
		Param(op.PublicSrv.PathParameter("expID", "Identifier of the experiment.").DataType("string")).
		Writes(Exp{}))

	restful.Add(op.PublicSrv)
}
