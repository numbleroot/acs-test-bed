package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type MetricsInt64 struct {
	Timestamp int64
	Values    []int64
}

type MetricsFloat64 struct {
	Timestamp int64
	Values    []float64
}

type MetricLatency struct {
	MsgID            int64
	SendTimestamp    int64
	ReceiveTimestamp int64
	Latency          float64
}

type Run struct {
	TimestampLowest            int64
	TimestampHighest           int64
	ClientsSentKiBytesHighest  []float64
	ClientsRecvdKiBytesHighest []float64
	ClientsCPULoad             []float64
	ClientsMemLoad             []float64
	Latencies                  [][]*MetricLatency
	ServersSentKiBytesHighest  []float64
	ServersRecvdKiBytesHighest []float64
	ServersCPULoad             []float64
	ServersMemLoad             []float64
	Mixes                      []string
	MsgsPerMix                 [][]int64
}

type Setting struct {
	Runs []*Run
}

type Experiment struct {
	ZenoClients0500 *Setting
	ZenoClients1000 *Setting
	PungClients0500 *Setting
	PungClients1000 *Setting
}

func (set *Setting) AppendRun(runPath string, numMsgsToCalc int64) {

	run := &Run{
		TimestampLowest:            (1 << 63) - 1,
		TimestampHighest:           0,
		ClientsSentKiBytesHighest:  make([]float64, 0, 1000),
		ClientsRecvdKiBytesHighest: make([]float64, 0, 1000),
		ServersSentKiBytesHighest:  make([]float64, 0, 1000),
		ServersRecvdKiBytesHighest: make([]float64, 0, 1000),
	}

	clientsPath := filepath.Join(runPath, "clients")
	serversPath := filepath.Join(runPath, "servers")

	// Determine lowest and highest relevant
	// timestamp of run while ingesting message
	// latency metrics.
	err := run.AddLatency(clientsPath, numMsgsToCalc)
	if err != nil {
		fmt.Printf("Ingesting client message latency metrics failed: %v\n", err)
		os.Exit(1)
	}

	// Read into memory system metrics from clients.
	err = run.AddSentBytes(clientsPath, true)
	if err != nil {
		fmt.Printf("Ingesting client sent kibytes metrics failed: %v\n", err)
		os.Exit(1)
	}

	err = run.AddRecvdBytes(clientsPath, true)
	if err != nil {
		fmt.Printf("Ingesting client received kibytes metrics failed: %v\n", err)
		os.Exit(1)
	}

	/*
		err = run.AddCPULoad(clientsPath, true)
		if err != nil {
			fmt.Printf("Ingesting client CPU load metrics failed: %v\n", err)
			os.Exit(1)
		}

		err = run.AddMemLoad(clientsPath, true)
		if err != nil {
			fmt.Printf("Ingesting client memory load metrics failed: %v\n", err)
			os.Exit(1)
		}
	*/

	// Read into memory system metrics from servers.
	err = run.AddSentBytes(serversPath, false)
	if err != nil {
		fmt.Printf("Ingesting server sent kibytes metrics failed: %v\n", err)
		os.Exit(1)
	}

	err = run.AddRecvdBytes(serversPath, false)
	if err != nil {
		fmt.Printf("Ingesting server received kibytes metrics failed: %v\n", err)
		os.Exit(1)
	}

	/*
		err = run.AddCPULoad(serversPath, false)
		if err != nil {
			fmt.Printf("Ingesting server CPU load metrics failed: %v\n", err)
			os.Exit(1)
		}

		err = run.AddMemLoad(serversPath, false)
		if err != nil {
			fmt.Printf("Ingesting server memory load metrics failed: %v\n", err)
			os.Exit(1)
		}

		// If this is zeno being evaluated, also read
		// in metrics about the number of messages in
		// each pool.
		err = run.AddMsgsPerMix(serversPath)
		if err != nil {
			fmt.Printf("Ingesting server sent bytes metrics failed: %v\n", err)
			os.Exit(1)
		}
	*/

	// Append newly created run to all runs.
	set.Runs = append(set.Runs, run)
}

func (set *Setting) MetricsToFiles(settingsPath string) error {

	// Write bandwidth data for clients and servers.
	err := set.BandwidthToFiles(settingsPath)
	if err != nil {
		return err
	}

	/*
		// Write load data for clients and servers.
		err = set.LoadToFiles(settingsPath)
		if err != nil {
			return err
		}
	*/

	// Write message latencies for clients.
	err = set.LatenciesToFile(settingsPath)
	if err != nil {
		return err
	}

	/*
		// Write messages-per-mix data for servers.
		err = set.MsgsPerMixToFile(settingsPath)
		if err != nil {
			return err
		}
	*/

	return nil
}

func main() {

	// Expect command-line arguments.
	experimentPathFlag := flag.String("experimentPath", "", "Specify the file system location of the directory containing the metric files for one experiment.")
	numMsgsToCalcFlag := flag.Int("numMsgsToCalc", 25, "Calculate statistics for this number of measured messages.")
	flag.Parse()

	numMsgsToCalc := int64(*numMsgsToCalcFlag)

	experimentPath, err := filepath.Abs(*experimentPathFlag)
	if err != nil {
		fmt.Printf("Error converting metrics path '%s' into absolute path: %v\n", *experimentPathFlag, err)
		os.Exit(1)
	}

	// Prepare storage space in experiment struct.

	foldersZenoClients0500, err := ioutil.ReadDir(filepath.Join(experimentPath, "zeno", "clients-0500"))
	if err != nil {

		if !strings.Contains(err.Error(), "no such file or directory") {
			fmt.Printf("Failed to retrieve all folders for zeno's 500 clients setting: %v\n", err)
			os.Exit(1)
		}
	}

	foldersZenoClients1000, err := ioutil.ReadDir(filepath.Join(experimentPath, "zeno", "clients-1000"))
	if err != nil {

		if !strings.Contains(err.Error(), "no such file or directory") {
			fmt.Printf("Failed to retrieve all folders for zeno's 1000 clients setting: %v\n", err)
			os.Exit(1)
		}
	}

	foldersPungClients0500, err := ioutil.ReadDir(filepath.Join(experimentPath, "pung", "clients-0500"))
	if err != nil {

		if !strings.Contains(err.Error(), "no such file or directory") {
			fmt.Printf("Failed to retrieve all folders for pung's 500 clients setting: %v\n", err)
			os.Exit(1)
		}
	}

	foldersPungClients1000, err := ioutil.ReadDir(filepath.Join(experimentPath, "pung", "clients-1000"))
	if err != nil {

		if !strings.Contains(err.Error(), "no such file or directory") {
			fmt.Printf("Failed to retrieve all folders for pung's 1000 clients setting: %v\n", err)
			os.Exit(1)
		}
	}

	experiment := &Experiment{
		ZenoClients0500: &Setting{Runs: make([]*Run, 0, len(foldersZenoClients0500))},
		ZenoClients1000: &Setting{Runs: make([]*Run, 0, len(foldersZenoClients1000))},
		PungClients0500: &Setting{Runs: make([]*Run, 0, len(foldersPungClients0500))},
		PungClients1000: &Setting{Runs: make([]*Run, 0, len(foldersPungClients1000))},
	}

	// Append each run to internal structures.

	for i := range foldersZenoClients0500 {

		if foldersZenoClients0500[i].IsDir() {
			experiment.ZenoClients0500.AppendRun(filepath.Join(experimentPath, "zeno", "clients-0500", foldersZenoClients0500[i].Name()), numMsgsToCalc)
		}
	}

	for i := range foldersZenoClients1000 {

		if foldersZenoClients1000[i].IsDir() {
			experiment.ZenoClients1000.AppendRun(filepath.Join(experimentPath, "zeno", "clients-1000", foldersZenoClients1000[i].Name()), numMsgsToCalc)
		}
	}

	for i := range foldersPungClients0500 {

		if foldersPungClients0500[i].IsDir() {
			experiment.PungClients0500.AppendRun(filepath.Join(experimentPath, "pung", "clients-0500", foldersPungClients0500[i].Name()), numMsgsToCalc)
		}
	}

	for i := range foldersPungClients1000 {

		if foldersPungClients1000[i].IsDir() {
			experiment.PungClients1000.AppendRun(filepath.Join(experimentPath, "pung", "clients-1000", foldersPungClients1000[i].Name()), numMsgsToCalc)
		}
	}

	// Write all metrics to files ready to be
	// turned into figures by genplots.

	if len(experiment.ZenoClients0500.Runs) > 0 {

		err = experiment.ZenoClients0500.MetricsToFiles(filepath.Join(experimentPath, "zeno", "clients-0500"))
		if err != nil {
			fmt.Printf("Failed to store calculated statistics for zeno's setting of 500 clients to files: %v\n", err)
			os.Exit(1)
		}
	}

	if len(experiment.ZenoClients1000.Runs) > 0 {

		err = experiment.ZenoClients1000.MetricsToFiles(filepath.Join(experimentPath, "zeno", "clients-1000"))
		if err != nil {
			fmt.Printf("Failed to store calculated statistics for zeno's setting of 1000 clients to files: %v\n", err)
			os.Exit(1)
		}
	}

	if len(experiment.PungClients0500.Runs) > 0 {

		err = experiment.PungClients0500.MetricsToFiles(filepath.Join(experimentPath, "pung", "clients-0500"))
		if err != nil {
			fmt.Printf("Failed to store calculated statistics for pung's setting of 500 clients to files: %v\n", err)
			os.Exit(1)
		}
	}

	if len(experiment.PungClients1000.Runs) > 0 {

		err = experiment.PungClients1000.MetricsToFiles(filepath.Join(experimentPath, "pung", "clients-1000"))
		if err != nil {
			fmt.Printf("Failed to store calculated statistics for pung's setting of 1000 clients to files: %v\n", err)
			os.Exit(1)
		}
	}
}
