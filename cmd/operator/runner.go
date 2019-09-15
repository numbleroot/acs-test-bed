package main

// RunExperiments is the authoritative goroutine
// for provisioning machines and conducting all
// experiments queued in the Operator.
func (op *Operator) RunExperiments() {

	for {

		// Wait for new experiment to be queued.
		exp := <-op.PublicChan

		op.Lock()
		op.ExpInProgress = exp
		op.Unlock()

		// TODO: Spawn all server machines.

		// TODO: Verify enough servers ready.

		// TODO: Spawn all client machines.

		// TODO: Verify enough clients ready.

		// TODO: Conduct experiment.

		// TODO: Wait for all clients to signal completion.

		op.Exps[exp].Concluded = true

		op.Lock()
		op.ExpInProgress = ""
		op.Unlock()
	}
}
