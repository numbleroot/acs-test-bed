package main

type MixMetrics struct {
	*SystemMetrics
	MetricsPath string
	MsgsPerPool []int64
}

func (mixM *MixMetrics) AddMsgsPerPool(path string) error {
	return nil
}

func (mixM *MixMetrics) CalcAndWritePoolSizes(path string) error {
	return nil
}
