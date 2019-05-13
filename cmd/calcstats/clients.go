package main

type ClientMetrics struct {
	*SystemMetrics
	Latency map[int]int64
}

func (clM *ClientMetrics) AddLatency(path string) error {
	return nil
}

func (clM *ClientMetrics) CalcAndWriteLatency(path string) error {
	return nil
}
