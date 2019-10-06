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
	ClientsSentMiBytesHighest  []float64
	ClientsRecvdMiBytesHighest []float64
	ClientsCPULoad             []float64
	ClientsMemLoad             []float64
	Latencies                  [][]*MetricLatency
	ServersSentMiBytesHighest  []float64
	ServersRecvdMiBytesHighest []float64
	ServersCPULoad             []float64
	ServersMemLoad             []float64
	Mixes                      []string
	MsgsPerMix                 [][]int64
}

type Setting struct {
	Runs []*Run
}

type Experiment struct {
	ZenoClients0500     *Setting
	ZenoClients1000     *Setting
	VuvuzelaClients0500 *Setting
	VuvuzelaClients1000 *Setting
	PungClients0500     *Setting
	PungClients1000     *Setting
}

func (set *Setting) AppendRun(runPath string, numMsgsToCalc int64) {

	run := &Run{
		TimestampLowest:            (1 << 63) - 1,
		TimestampHighest:           0,
		ClientsSentMiBytesHighest:  make([]float64, 0, 1000),
		ClientsRecvdMiBytesHighest: make([]float64, 0, 1000),
		ClientsCPULoad:             make([]float64, 0, 50000),
		ClientsMemLoad:             make([]float64, 0, 50000),
		ServersSentMiBytesHighest:  make([]float64, 0, 50),
		ServersRecvdMiBytesHighest: make([]float64, 0, 50),
		ServersCPULoad:             make([]float64, 0, 50000),
		ServersMemLoad:             make([]float64, 0, 50000),
	}

	clientsPath := filepath.Join(runPath, "clients")
	serversPath := filepath.Join(runPath, "servers")

	// Determine lowest and highest relevant
	// timestamp of run while ingesting message
	// latency metrics.
	err := run.AddLatency(clientsPath, numMsgsToCalc)
	if err != nil {
		fmt.Printf("Ingesting clients message latency metrics failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Done adding clients latency for %s\n", runPath)

	// Read in memory system metrics from clients.
	err = run.AddSentBytes(clientsPath, true)
	if err != nil {
		fmt.Printf("Ingesting clients sent mebibytes metrics failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Done adding clients sent bytes for %s\n", runPath)

	err = run.AddRecvdBytes(clientsPath, true)
	if err != nil {
		fmt.Printf("Ingesting client received mebibytes metrics failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Done adding clients received bytes for %s\n", runPath)

	err = run.AddCPULoad(clientsPath, true)
	if err != nil {
		fmt.Printf("Ingesting clients CPU load metrics failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Done adding clients CPU load for %s\n", runPath)

	err = run.AddMemLoad(clientsPath, true)
	if err != nil {
		fmt.Printf("Ingesting clients memory load metrics failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Done adding clients mem load for %s\n", runPath)

	// Read in memory system metrics from servers.
	err = run.AddSentBytes(serversPath, false)
	if err != nil {
		fmt.Printf("Ingesting servers sent mebibytes metrics failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Done adding servers sent bytes for %s\n", runPath)

	err = run.AddRecvdBytes(serversPath, false)
	if err != nil {
		fmt.Printf("Ingesting servers received mebibytes metrics failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Done adding servers received bytes for %s\n", runPath)

	err = run.AddCPULoad(serversPath, false)
	if err != nil {
		fmt.Printf("Ingesting servers CPU load metrics failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Done adding servers CPU load for %s\n", runPath)

	err = run.AddMemLoad(serversPath, false)
	if err != nil {
		fmt.Printf("Ingesting servers memory load metrics failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Done adding servers mem load for %s\n", runPath)

	// If this is zeno being evaluated, also read
	// in metrics about the number of messages in
	// each pool.
	err = run.AddMsgsPerMix(serversPath)
	if err != nil {
		fmt.Printf("Ingesting server sent bytes metrics failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Done adding servers messages per mix for %s\n\n", runPath)

	// Append newly created run to all runs.
	set.Runs = append(set.Runs, run)
}

func (set *Setting) MetricsToFiles(settingsPath string) error {

	// Write traffic data for clients and servers.
	err := set.TrafficToFiles(settingsPath)
	if err != nil {
		return err
	}

	fmt.Printf("Done writing bandwidth to file for %s\n", settingsPath)

	// Write load data for clients and servers.
	err = set.LoadToFiles(settingsPath)
	if err != nil {
		return err
	}

	fmt.Printf("Done writing load to file for %s\n", settingsPath)

	// Write message latencies for clients.
	err = set.LatenciesToFile(settingsPath)
	if err != nil {
		return err
	}

	fmt.Printf("Done writing latencies to file for %s\n", settingsPath)

	// Write total experiment times for clients.
	err = set.TotalExpTimesToFile(settingsPath)
	if err != nil {
		return err
	}

	fmt.Printf("Done writing total experiment times to file for %s\n", settingsPath)

	// Write messages-per-mix data for mix nodes
	// in a zeno evaluation.
	err = set.MsgsPerMixToFile(settingsPath)
	if err != nil {
		return err
	}

	fmt.Printf("Done writing messages per mix on zeno to file for %s\n", settingsPath)

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

	foldersVuvuzelaClients0500, err := ioutil.ReadDir(filepath.Join(experimentPath, "vuvuzela", "clients-0500"))
	if err != nil {

		if !strings.Contains(err.Error(), "no such file or directory") {
			fmt.Printf("Failed to retrieve all folders for Vuvuzela's 500 clients setting: %v\n", err)
			os.Exit(1)
		}
	}

	foldersVuvuzelaClients1000, err := ioutil.ReadDir(filepath.Join(experimentPath, "vuvuzela", "clients-1000"))
	if err != nil {

		if !strings.Contains(err.Error(), "no such file or directory") {
			fmt.Printf("Failed to retrieve all folders for Vuvuzela's 1000 clients setting: %v\n", err)
			os.Exit(1)
		}
	}

	foldersPungClients0500, err := ioutil.ReadDir(filepath.Join(experimentPath, "pung", "clients-0500"))
	if err != nil {

		if !strings.Contains(err.Error(), "no such file or directory") {
			fmt.Printf("Failed to retrieve all folders for Pung's 500 clients setting: %v\n", err)
			os.Exit(1)
		}
	}

	foldersPungClients1000, err := ioutil.ReadDir(filepath.Join(experimentPath, "pung", "clients-1000"))
	if err != nil {

		if !strings.Contains(err.Error(), "no such file or directory") {
			fmt.Printf("Failed to retrieve all folders for Pung's 1000 clients setting: %v\n", err)
			os.Exit(1)
		}
	}

	experiment := &Experiment{
		ZenoClients0500:     &Setting{Runs: make([]*Run, 0, len(foldersZenoClients0500))},
		ZenoClients1000:     &Setting{Runs: make([]*Run, 0, len(foldersZenoClients1000))},
		VuvuzelaClients0500: &Setting{Runs: make([]*Run, 0, len(foldersVuvuzelaClients0500))},
		VuvuzelaClients1000: &Setting{Runs: make([]*Run, 0, len(foldersVuvuzelaClients1000))},
		PungClients0500:     &Setting{Runs: make([]*Run, 0, len(foldersPungClients0500))},
		PungClients1000:     &Setting{Runs: make([]*Run, 0, len(foldersPungClients1000))},
	}

	// Append each run to internal structures.

	for i := range foldersZenoClients0500 {

		if foldersZenoClients0500[i].IsDir() {
			experiment.ZenoClients0500.AppendRun(filepath.Join(experimentPath, "zeno", "clients-0500",
				foldersZenoClients0500[i].Name()), numMsgsToCalc)
		}
	}

	for i := range foldersZenoClients1000 {

		if foldersZenoClients1000[i].IsDir() {
			experiment.ZenoClients1000.AppendRun(filepath.Join(experimentPath, "zeno", "clients-1000",
				foldersZenoClients1000[i].Name()), numMsgsToCalc)
		}
	}

	for i := range foldersVuvuzelaClients0500 {

		if foldersVuvuzelaClients0500[i].IsDir() {
			experiment.VuvuzelaClients0500.AppendRun(filepath.Join(experimentPath, "vuvuzela", "clients-0500",
				foldersVuvuzelaClients0500[i].Name()), numMsgsToCalc)
		}
	}

	for i := range foldersVuvuzelaClients1000 {

		if foldersVuvuzelaClients1000[i].IsDir() {
			experiment.VuvuzelaClients1000.AppendRun(filepath.Join(experimentPath, "vuvuzela", "clients-1000",
				foldersVuvuzelaClients1000[i].Name()), numMsgsToCalc)
		}
	}

	for i := range foldersPungClients0500 {

		if foldersPungClients0500[i].IsDir() {
			experiment.PungClients0500.AppendRun(filepath.Join(experimentPath, "pung", "clients-0500",
				foldersPungClients0500[i].Name()), numMsgsToCalc)
		}
	}

	for i := range foldersPungClients1000 {

		if foldersPungClients1000[i].IsDir() {
			experiment.PungClients1000.AppendRun(filepath.Join(experimentPath, "pung", "clients-1000",
				foldersPungClients1000[i].Name()), numMsgsToCalc)
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

	if len(experiment.VuvuzelaClients0500.Runs) > 0 {

		err = experiment.VuvuzelaClients0500.MetricsToFiles(filepath.Join(experimentPath, "vuvuzela", "clients-0500"))
		if err != nil {
			fmt.Printf("Failed to store calculated statistics for Vuvuzela's setting of 500 clients to files: %v\n", err)
			os.Exit(1)
		}
	}

	if len(experiment.VuvuzelaClients1000.Runs) > 0 {

		err = experiment.VuvuzelaClients1000.MetricsToFiles(filepath.Join(experimentPath, "vuvuzela", "clients-1000"))
		if err != nil {
			fmt.Printf("Failed to store calculated statistics for Vuvuzela's setting of 1000 clients to files: %v\n", err)
			os.Exit(1)
		}
	}

	if len(experiment.PungClients0500.Runs) > 0 {

		err = experiment.PungClients0500.MetricsToFiles(filepath.Join(experimentPath, "pung", "clients-0500"))
		if err != nil {
			fmt.Printf("Failed to store calculated statistics for Pung's setting of 500 clients to files: %v\n", err)
			os.Exit(1)
		}
	}

	if len(experiment.PungClients1000.Runs) > 0 {

		err = experiment.PungClients1000.MetricsToFiles(filepath.Join(experimentPath, "pung", "clients-1000"))
		if err != nil {
			fmt.Printf("Failed to store calculated statistics for Pung's setting of 1000 clients to files: %v\n", err)
			os.Exit(1)
		}
	}
}
