package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {

	// Expect command-line arguments.
	systemFlag := flag.String("system", "", "Specify which ACS to evaluate: 'zeno', 'vuvuzela', 'pung'.")
	metricsPathFlag := flag.String("metricsPath", "./results", "Specify the file system location of the directory containing the metrics files.")
	numMsgsToCalcFlag := flag.Int("numMsgsToCalc", 10, "Calculate statistics for this number of measured messages.")
	flag.Parse()

	// Enforce arguments to be set.
	if *systemFlag == "" {
		fmt.Printf("Missing arguments, please provide values for all flags: '-system'.\n")
		os.Exit(1)
	}

	system := *systemFlag
	numMsgsToCalc := *numMsgsToCalcFlag

	fmt.Printf("system=%v, numMsgsToCalc=%v\n", system, numMsgsToCalc)

	metricsPath, err := filepath.Abs(*metricsPathFlag)
	if err != nil {
		fmt.Printf("Error converting metrics path '%s' into absolute path: %v\n", *metricsPathFlag, err)
		os.Exit(1)
	}
	clientMetricsPath := filepath.Join(metricsPath, "clients")
	mixMetricsPath := filepath.Join(metricsPath, "mixes")

	clientMetrics := &ClientMetrics{
		SystemMetrics: &SystemMetrics{
			SentBytesRaw:        make(map[int64][]int64),
			RecvdBytesRaw:       make(map[int64][]int64),
			MemoryRaw:           make(map[int64][]float64),
			LoadRaw:             make(map[int64][]float64),
			TimestampLowerBound: (1 << 63) - 1,
			TimestampUpperBound: 0,
		},
		MetricsPath:   clientMetricsPath,
		NumMsgsToCalc: int64(numMsgsToCalc),
		Latencies:     make([][]*MetricLatency, 0, 10),
	}

	// Scan metrics directory of clients.
	err = filepath.Walk(clientMetricsPath, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		switch filepath.Base(path) {

		case "traffic_outgoing.evaluation":
			err = clientMetrics.AddSentBytes(path)

		case "traffic_incoming.evaluation":
			err = clientMetrics.AddRecvdBytes(path)

		case "load_unixnano.evaluation":
			err = clientMetrics.AddLoad(path)

		case "mem_unixnano.evaluation":
			err = clientMetrics.AddMem(path)

		case "send_unixnano.evaluation":
			err = clientMetrics.AddLatency(path)
		}
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		fmt.Printf("Ingesting client metrics failed: %v\n", err)
		os.Exit(1)
	}

	mixMetrics := &MixMetrics{
		SystemMetrics: &SystemMetrics{
			SentBytesRaw:        make(map[int64][]int64),
			RecvdBytesRaw:       make(map[int64][]int64),
			MemoryRaw:           make(map[int64][]float64),
			LoadRaw:             make(map[int64][]float64),
			TimestampLowerBound: (1 << 63) - 1,
			TimestampUpperBound: 0,
		},
		MetricsPath: mixMetricsPath,
		Mixes:       make([]string, 0, 10),
		MsgsPerMix:  make([][]int64, 0, 10),
	}

	// Scan metrics directory of clients.
	err = filepath.Walk(mixMetricsPath, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		switch filepath.Base(path) {

		case "traffic_outgoing.evaluation":
			err = mixMetrics.AddSentBytes(path)

		case "traffic_incoming.evaluation":
			err = mixMetrics.AddRecvdBytes(path)

		case "load_unixnano.evaluation":
			err = mixMetrics.AddLoad(path)

		case "mem_unixnano.evaluation":
			err = mixMetrics.AddMem(path)

		case "pool-sizes_round.evaluation":
			err = mixMetrics.AddMsgsPerMix(path)
		}
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		fmt.Printf("Ingesting mix metrics failed: %v\n", err)
		os.Exit(1)
	}

	// Filter and sort system metrics.
	err = clientMetrics.SystemSortByTimestamp()
	if err != nil {
		fmt.Printf("Failed to order system metrics of clients: %v\n", err)
		os.Exit(1)
	}

	err = mixMetrics.SystemSortByTimestamp()
	if err != nil {
		fmt.Printf("Failed to order system metrics of mixes: %v\n", err)
		os.Exit(1)
	}

	// Write out system metrics, ready to be
	// boxplotted with Python script.
	err = clientMetrics.SystemStoreForBoxplots(clientMetricsPath)
	if err != nil {
		fmt.Printf("Failed to write out system metrics boxplot files for clients: %v\n", err)
		os.Exit(1)
	}

	err = mixMetrics.SystemStoreForBoxplots(mixMetricsPath)
	if err != nil {
		fmt.Printf("Failed to write out system metrics boxplot files for mixes: %v\n", err)
		os.Exit(1)
	}

	// Write out message latency metrics for clients.
	// Ready to be boxplotted with Python script.
	err = clientMetrics.ClientStoreForBoxplot()
	if err != nil {
		fmt.Printf("Failed to write out message latency metrics boxplot files for clients: %v\n", err)
		os.Exit(1)
	}

	// Write out message counts per mix.
	// Ready to be boxplotted with Python script.
	err = mixMetrics.MixStoreForPlot()
	if err != nil {
		fmt.Printf("Failed to write out message number metrics boxplot files for mixes: %v\n", err)
		os.Exit(1)
	}
}
