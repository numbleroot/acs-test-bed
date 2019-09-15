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
	}
}
