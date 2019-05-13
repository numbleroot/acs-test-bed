package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type LoadMetric struct {
	User   float64
	Nice   float64
	System float64
	IOWait float64
	Idle   float64
}

type MemoryMetric struct {
	TotalKB int64
	AvailKB int64
}

type SystemMetrics struct {
	SentBytes  map[int64][]int64
	RecvdBytes map[int64][]int64
	Load       map[int64][]*LoadMetric
	Memory     map[int64][]*MemoryMetric
}

type ClientMetrics struct {
	*SystemMetrics
	Latency map[int]int64
}

type MixMetrics struct {
	*SystemMetrics
	MsgsPerPool []int64
}

func (sysM *SystemMetrics) AddSentBytes(path string) error {

	// Ingest supplied file.
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	content = bytes.TrimSpace(content)

	// Split file contents into lines.
	lines := strings.Split(string(content), "\n")
	for i := range lines {

		// Split line at whitespace characters.
		metric := strings.Fields(lines[i])

		// Convert first element to timestamp.
		timestamp, err := strconv.ParseInt(metric[0], 10, 64)
		if err != nil {
			return err
		}

		// Convert second element to metric
		// we are interested in.
		value, err := strconv.ParseInt(metric[1], 10, 64)
		if err != nil {
			return err
		}

		// Append to corresponding slice of values.
		sysM.SentBytes[timestamp] = append(sysM.SentBytes[timestamp], value)
	}

	return nil
}

func (sysM *SystemMetrics) AddRecvdBytes(path string) error {

	// Ingest supplied file.
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	content = bytes.TrimSpace(content)

	// Split file contents into lines.
	lines := strings.Split(string(content), "\n")
	for i := range lines {

		// Split line at whitespace characters.
		metric := strings.Fields(lines[i])

		// Convert first element to timestamp.
		timestamp, err := strconv.ParseInt(metric[0], 10, 64)
		if err != nil {
			return err
		}

		// Convert second element to metric
		// we are interested in.
		value, err := strconv.ParseInt(metric[1], 10, 64)
		if err != nil {
			return err
		}

		// Append to corresponding slice of values.
		sysM.RecvdBytes[timestamp] = append(sysM.RecvdBytes[timestamp], value)
	}

	return nil
}

func (sysM *SystemMetrics) AddLoad(path string) error {

	// Ingest supplied file.
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	content = bytes.TrimSpace(content)

	// Split file contents into lines.
	lines := strings.Split(string(content), "\n")
	for i := range lines {

		// Split line at whitespace characters.
		metric := strings.Fields(lines[i])

		// Convert first element to timestamp.
		timestamp, err := strconv.ParseInt(metric[0], 10, 64)
		if err != nil {
			return err
		}

		// Convert following elements to load metrics.

		loadUser, err := strconv.ParseFloat(strings.TrimPrefix(metric[1], "usr:"), 64)
		if err != nil {
			return err
		}

		loadNice, err := strconv.ParseFloat(strings.TrimPrefix(metric[2], "nice:"), 64)
		if err != nil {
			return err
		}

		loadSys, err := strconv.ParseFloat(strings.TrimPrefix(metric[3], "sys:"), 64)
		if err != nil {
			return err
		}

		loadIOWait, err := strconv.ParseFloat(strings.TrimPrefix(metric[4], "iowait:"), 64)
		if err != nil {
			return err
		}

		loadIdle, err := strconv.ParseFloat(strings.TrimPrefix(metric[5], "idle:"), 64)
		if err != nil {
			return err
		}

		// Append to corresponding slice of values.
		sysM.Load[timestamp] = append(sysM.Load[timestamp], &LoadMetric{
			User:   loadUser,
			Nice:   loadNice,
			System: loadSys,
			IOWait: loadIOWait,
			Idle:   loadIdle,
		})
	}

	return nil
}

func (sysM *SystemMetrics) AddMem(path string) error {

	// Ingest supplied file.
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	content = bytes.TrimSpace(content)

	// Split file contents into lines.
	lines := strings.Split(string(content), "\n")
	for i := range lines {

		// Split line at whitespace characters.
		metric := strings.Fields(lines[i])

		// Convert first element to timestamp.
		timestamp, err := strconv.ParseInt(metric[0], 10, 64)
		if err != nil {
			return err
		}

		// Convert following elements to memory metrics.

		memTotal, err := strconv.ParseInt(strings.TrimPrefix(metric[1], "totalKB:"), 10, 64)
		if err != nil {
			return err
		}

		memAvail, err := strconv.ParseInt(strings.TrimPrefix(metric[2], "availKB:"), 10, 64)
		if err != nil {
			return err
		}

		// Append to corresponding slice of values.
		sysM.Memory[timestamp] = append(sysM.Memory[timestamp], &MemoryMetric{
			TotalKB: memTotal,
			AvailKB: memAvail,
		})
	}

	return nil
}

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
			SentBytes:  make(map[int64][]int64),
			RecvdBytes: make(map[int64][]int64),
			Load:       make(map[int64][]*LoadMetric),
			Memory:     make(map[int64][]*MemoryMetric),
		},
		Latency: make(map[int]int64),
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
			SentBytes:  make(map[int64][]int64),
			RecvdBytes: make(map[int64][]int64),
			Load:       make(map[int64][]*LoadMetric),
			Memory:     make(map[int64][]*MemoryMetric),
		},
		MsgsPerPool: make([]int64, 0, 7),
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
}
