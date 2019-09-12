package main

import "github.com/emicklei/go-restful"

// PrepareInternalSrv initializes all API-related
// things in order to expose an internal-facing
// API endpoint for conducting experiments.
func (op *Operator) PrepareInternalSrv() {

	op.InternalSrv = new(restful.WebService)

	op.InternalSrv.Path("/workers").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	restful.Add(op.InternalSrv)
}
