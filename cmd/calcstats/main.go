package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
	TimestampLowest          int64
	TimestampHighest         int64
	ClientsSentBytesHighest  []int64
	ClientsRecvdBytesHighest []int64
	ClientsMemory            []*MetricsFloat64
	ClientsLoad              []*MetricsFloat64
	Latencies                [][]*MetricLatency
	ServersSentBytesHighest  []int64
	ServersRecvdBytesHighest []int64
	ServersMemory            []*MetricsFloat64
	ServersLoad              []*MetricsFloat64
	Mixes                    []string
	MsgsPerMix               [][]int64
}

type Setting struct {
	Runs                       []*Run
	ClientsBandwidthHighestAvg float64
	ClientsBandwidthHighestMed float64
	ServersBandwidthHighestAvg float64
	ServersBandwidthHighestMed float64
	ClientsLoadLowestToHighest []float64
	ServersLoadLowestToHighest []float64
	LatenciesLowestToHighest   []float64
}

type Experiment struct {
	ZenoClients0500 *Setting
	ZenoClients1000 *Setting
	ZenoClients2000 *Setting
	PungClients0500 *Setting
	PungClients1000 *Setting
	PungClients2000 *Setting
}

func (set *Setting) AppendRun(runPath string, systemUnderEval string, numMsgsToCalc int64) {

	run := &Run{
		TimestampLowest:  (1 << 63) - 1,
		TimestampHighest: 0,
	}

	clientsPath := filepath.Join(runPath, "clients")
	serversPath := filepath.Join(runPath, "servers")

	// Determine lowest and highest relevant
	// timestamp of run while ingesting message
	// latency metrics.
	err := run.AddLatency(clientsPath, systemUnderEval, numMsgsToCalc)
	if err != nil {
		fmt.Printf("Ingesting client message latency metrics failed: %v\n", err)
		os.Exit(1)
	}

	// Read into memory system metrics from clients.
	err = run.AddSentBytes(clientsPath, true)
	if err != nil {
		fmt.Printf("Ingesting client sent bytes metrics failed: %v\n", err)
		os.Exit(1)
	}

	err = run.AddRecvdBytes(clientsPath, true)
	if err != nil {
		fmt.Printf("Ingesting client received bytes metrics failed: %v\n", err)
		os.Exit(1)
	}

	err = run.AddLoad(clientsPath, true)
	if err != nil {
		fmt.Printf("Ingesting client CPU load metrics failed: %v\n", err)
		os.Exit(1)
	}

	err = run.AddMem(clientsPath, true)
	if err != nil {
		fmt.Printf("Ingesting client memory load metrics failed: %v\n", err)
		os.Exit(1)
	}

	// Read into memory system metrics from servers.
	err = run.AddSentBytes(serversPath, false)
	if err != nil {
		fmt.Printf("Ingesting server sent bytes metrics failed: %v\n", err)
		os.Exit(1)
	}

	err = run.AddRecvdBytes(serversPath, false)
	if err != nil {
		fmt.Printf("Ingesting server received bytes metrics failed: %v\n", err)
		os.Exit(1)
	}

	err = run.AddLoad(serversPath, false)
	if err != nil {
		fmt.Printf("Ingesting server CPU load metrics failed: %v\n", err)
		os.Exit(1)
	}

	err = run.AddMem(serversPath, false)
	if err != nil {
		fmt.Printf("Ingesting server memory load metrics failed: %v\n", err)
		os.Exit(1)
	}

	if systemUnderEval == "zeno" {

		// If this is zeno being evaluated, also read
		// in metrics about the number of messages in
		// each pool.
		err = run.AddMsgsPerMix(serversPath)
		if err != nil {
			fmt.Printf("Ingesting server sent bytes metrics failed: %v\n", err)
			os.Exit(1)
		}
	}

	// Append newly created run to all runs.
	set.Runs = append(set.Runs, run)
}

func (set *Setting) MetricsToFiles(settingsPath string, systemUnderEval string) error {

	// Write bandwidth data for clients and servers.
	err := set.BandwidthToFiles(settingsPath)
	if err != nil {
		return err
	}

	// Write load data for clients and servers.
	err = set.LoadToFiles(settingsPath)
	if err != nil {
		return err
	}

	// Write message latencies for clients.
	err = set.LatenciesToFile(settingsPath)
	if err != nil {
		return err
	}

	if systemUnderEval == "zeno" {

		// Write messages-per-mix data for servers.
		err = set.MsgsPerMixToFile(settingsPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {

	// Expect command-line arguments.
	systemFlag := flag.String("system", "", "Specify which ACS to evaluate: 'zeno', 'vuvuzela', 'pung'.")
	experimentPathFlag := flag.String("experimentPath", "./results/01_tc-off_proc-off", "Specify the file system location of the directory containing the metric files for one experiment.")
	numMsgsToCalcFlag := flag.Int("numMsgsToCalc", 10, "Calculate statistics for this number of measured messages.")
	flag.Parse()

	// System flag has to be one of three values.
	if *systemFlag != "zeno" && *systemFlag != "vuvuzela" && *systemFlag != "pung" {
		fmt.Printf("Flag '-system' requires one of the three values: 'zeno', 'vuvuzela', or 'pung'.")
		os.Exit(1)
	}

	system := *systemFlag
	numMsgsToCalc := int64(*numMsgsToCalcFlag)

	experimentPath, err := filepath.Abs(*experimentPathFlag)
	if err != nil {
		fmt.Printf("Error converting metrics path '%s' into absolute path: %v\n", *experimentPathFlag, err)
		os.Exit(1)
	}

	// Prepare storage space in experiment struct.

	foldersZenoClients0500, err := ioutil.ReadDir(filepath.Join(experimentPath, "zeno", "clients-0500"))
	if err != nil {
		fmt.Printf("Failed to retrieve all folders for zeno's 500 clients setting: %v\n", err)
		os.Exit(1)
	}

	foldersZenoClients1000, err := ioutil.ReadDir(filepath.Join(experimentPath, "zeno", "clients-1000"))
	if err != nil {
		fmt.Printf("Failed to retrieve all folders for zeno's 1000 clients setting: %v\n", err)
		os.Exit(1)
	}

	foldersZenoClients2000, err := ioutil.ReadDir(filepath.Join(experimentPath, "zeno", "clients-2000"))
	if err != nil {
		fmt.Printf("Failed to retrieve all folders for zeno's 2000 clients setting: %v\n", err)
		os.Exit(1)
	}

	foldersPungClients0500, err := ioutil.ReadDir(filepath.Join(experimentPath, "pung", "clients-0500"))
	if err != nil {
		fmt.Printf("Failed to retrieve all folders for pung's 500 clients setting: %v\n", err)
		os.Exit(1)
	}

	foldersPungClients1000, err := ioutil.ReadDir(filepath.Join(experimentPath, "pung", "clients-1000"))
	if err != nil {
		fmt.Printf("Failed to retrieve all folders for pung's 1000 clients setting: %v\n", err)
		os.Exit(1)
	}

	foldersPungClients2000, err := ioutil.ReadDir(filepath.Join(experimentPath, "pung", "clients-2000"))
	if err != nil {
		fmt.Printf("Failed to retrieve all folders for pung's 2000 clients setting: %v\n", err)
		os.Exit(1)
	}

	experiment := &Experiment{
		ZenoClients0500: &Setting{Runs: make([]*Run, 0, len(foldersZenoClients0500))},
		ZenoClients1000: &Setting{Runs: make([]*Run, 0, len(foldersZenoClients1000))},
		ZenoClients2000: &Setting{Runs: make([]*Run, 0, len(foldersZenoClients2000))},
		PungClients0500: &Setting{Runs: make([]*Run, 0, len(foldersPungClients0500))},
		PungClients1000: &Setting{Runs: make([]*Run, 0, len(foldersPungClients1000))},
		PungClients2000: &Setting{Runs: make([]*Run, 0, len(foldersPungClients2000))},
	}

	// Append each run to internal structures.

	for i := range foldersZenoClients0500 {
		experiment.ZenoClients0500.AppendRun(filepath.Join(experimentPath, "zeno", "clients-0500", fmt.Sprintf("run-%02d", (i+1))), system, numMsgsToCalc)
	}

	for i := range foldersZenoClients1000 {
		experiment.ZenoClients1000.AppendRun(filepath.Join(experimentPath, "zeno", "clients-1000", fmt.Sprintf("run-%02d", (i+1))), system, numMsgsToCalc)
	}

	for i := range foldersZenoClients2000 {
		experiment.ZenoClients2000.AppendRun(filepath.Join(experimentPath, "zeno", "clients-2000", fmt.Sprintf("run-%02d", (i+1))), system, numMsgsToCalc)
	}

	for i := range foldersPungClients0500 {
		experiment.PungClients0500.AppendRun(filepath.Join(experimentPath, "pung", "clients-0500", fmt.Sprintf("run-%02d", (i+1))), system, numMsgsToCalc)
	}

	for i := range foldersPungClients1000 {
		experiment.PungClients1000.AppendRun(filepath.Join(experimentPath, "pung", "clients-1000", fmt.Sprintf("run-%02d", (i+1))), system, numMsgsToCalc)
	}

	for i := range foldersPungClients2000 {
		experiment.PungClients2000.AppendRun(filepath.Join(experimentPath, "pung", "clients-2000", fmt.Sprintf("run-%02d", (i+1))), system, numMsgsToCalc)
	}

	// Write all metrics to files ready to be
	// turned into figures by genplots.

	err = experiment.ZenoClients0500.MetricsToFiles(filepath.Join(experimentPath, "zeno", "clients-0500"), system)
	if err != nil {
		fmt.Printf("Failed to store calculated statistics for zeno's setting of 500 clients to files: %v\n", err)
		os.Exit(1)
	}

	err = experiment.ZenoClients1000.MetricsToFiles(filepath.Join(experimentPath, "zeno", "clients-1000"), system)
	if err != nil {
		fmt.Printf("Failed to store calculated statistics for zeno's setting of 1000 clients to files: %v\n", err)
		os.Exit(1)
	}

	err = experiment.ZenoClients2000.MetricsToFiles(filepath.Join(experimentPath, "zeno", "clients-2000"), system)
	if err != nil {
		fmt.Printf("Failed to store calculated statistics for zeno's setting of 2000 clients to files: %v\n", err)
		os.Exit(1)
	}

	err = experiment.PungClients0500.MetricsToFiles(filepath.Join(experimentPath, "pung", "clients-0500"), system)
	if err != nil {
		fmt.Printf("Failed to store calculated statistics for pung's setting of 500 clients to files: %v\n", err)
		os.Exit(1)
	}

	err = experiment.PungClients1000.MetricsToFiles(filepath.Join(experimentPath, "pung", "clients-1000"), system)
	if err != nil {
		fmt.Printf("Failed to store calculated statistics for pung's setting of 1000 clients to files: %v\n", err)
		os.Exit(1)
	}

	err = experiment.PungClients2000.MetricsToFiles(filepath.Join(experimentPath, "pung", "clients-2000"), system)
	if err != nil {
		fmt.Printf("Failed to store calculated statistics for pung's setting of 2000 clients to files: %v\n", err)
		os.Exit(1)
	}
}
