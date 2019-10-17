package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// MetricLatency captures one message
// transmitted between two clients, with
// an end-to-end transmission latency in
// seconds.
type MetricLatency struct {
	MsgID            int64
	SendTimestamp    int64
	ReceiveTimestamp int64
	Latency          float64
}

// Run constitutes the collection of metrics
// gathered in one evaluation run. The two
// timestamp values are lower and upper bounds
// in seconds, respectively, between which all
// other metrics are caputed.
type Run struct {
	NumServers                 float64
	NumClients                 float64
	TimestampLowest            int64
	TimestampHighest           int64
	Latencies                  [][]*MetricLatency
	ClientsSentMiBytesHighest  []float64
	ClientsRecvdMiBytesHighest []float64
	ClientsCPULoad             []float64
	ClientsMemLoad             []float64
	ServersSentMiBytesHighest  []float64
	ServersRecvdMiBytesHighest []float64
	ServersCPULoad             []float64
	ServersMemLoad             []float64
	Mixes                      []string
	MsgsPerMix                 [][]int64
}

// Setting is a helper struct to allow
// defining functions on a slice of Runs.
type Setting struct {
	Runs []*Run
}

// Experiment collects all metrics for
// all ACS under evaluation across all
// repetition runs carried out.
type Experiment struct {
	ZenoClients1000     *Setting
	ZenoClients2000     *Setting
	ZenoClients3000     *Setting
	VuvuzelaClients1000 *Setting
	VuvuzelaClients2000 *Setting
	VuvuzelaClients3000 *Setting
	PungClients1000     *Setting
	PungClients2000     *Setting
	PungClients3000     *Setting
}

// AppendRun reads in all metric files for one
// complete run of a system and appends the data
// as a new run to the internal state.
func (set *Setting) AppendRun(runPath string, numServers float64, numClients float64, numMsgsToCalc int64) {

	// Prepare space for a Run with the anticipated
	// maximum possible number of metrics per category.
	run := &Run{
		NumServers:                 numServers,
		NumClients:                 numClients,
		TimestampLowest:            (1 << 63) - 1,
		TimestampHighest:           0,
		Latencies:                  make([][]*MetricLatency, 0, 3000),
		ClientsSentMiBytesHighest:  make([]float64, 0, 3000),
		ClientsRecvdMiBytesHighest: make([]float64, 0, 3000),
		ClientsCPULoad:             make([]float64, 0, 75000),
		ClientsMemLoad:             make([]float64, 0, 75000),
		ServersSentMiBytesHighest:  make([]float64, 0, 25),
		ServersRecvdMiBytesHighest: make([]float64, 0, 25),
		ServersCPULoad:             make([]float64, 0, 1000),
		ServersMemLoad:             make([]float64, 0, 1000),
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

	if len(run.Latencies) != int(run.NumClients) {
		fmt.Printf("Run '%s' produced possibly wrong number of latency measurements (want: %d, saw: %d)\n",
			runPath, int(run.NumClients), len(run.Latencies))
	}

	fmt.Printf("Done adding clients latency for '%s'\n", runPath)

	// Read in highest value for number of outgoing
	// bytes on each client.
	err = run.AddHighestSentBytes(clientsPath, true)
	if err != nil {
		fmt.Printf("Ingesting clients sent mebibytes metrics failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Done adding clients sent bytes for '%s'\n", runPath)

	// Read in highest value for number of incoming
	// bytes on each client.
	err = run.AddHighestRecvdBytes(clientsPath, true)
	if err != nil {
		fmt.Printf("Ingesting client received mebibytes metrics failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Done adding clients received bytes for '%s'\n", runPath)

	err = run.AddCPULoad(clientsPath, true)
	if err != nil {
		fmt.Printf("Ingesting clients CPU load metrics failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Done adding clients CPU load for '%s'\n", runPath)

	err = run.AddMemLoad(clientsPath, true)
	if err != nil {
		fmt.Printf("Ingesting clients memory load metrics failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Done adding clients mem load for '%s'\n", runPath)

	err = run.AddHighestSentBytes(serversPath, false)
	if err != nil {
		fmt.Printf("Ingesting servers sent mebibytes metrics failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Done adding servers sent bytes for '%s'\n", runPath)

	err = run.AddHighestRecvdBytes(serversPath, false)
	if err != nil {
		fmt.Printf("Ingesting servers received mebibytes metrics failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Done adding servers received bytes for '%s'\n", runPath)

	err = run.AddCPULoad(serversPath, false)
	if err != nil {
		fmt.Printf("Ingesting servers CPU load metrics failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Done adding servers CPU load for '%s'\n", runPath)

	err = run.AddMemLoad(serversPath, false)
	if err != nil {
		fmt.Printf("Ingesting servers memory load metrics failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Done adding servers mem load for '%s'\n", runPath)

	// If this is zeno being evaluated, also read
	// in metrics about the number of messages in
	// each pool.
	err = run.AddMsgsPerMix(serversPath)
	if err != nil {
		fmt.Printf("Ingesting server sent bytes metrics failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Done adding servers messages per mix for '%s'\n\n", runPath)

	// Append newly created run to all runs.
	set.Runs = append(set.Runs, run)
}

// MetricsToFiles calls the appropriate functions
// for each metric to write out the gathered
// measurements into all respective files for
// subsequent visualization.
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

	fmt.Printf("Done writing messages per mix on zeno to file for %s\n\n", settingsPath)

	return nil
}

func main() {

	// Expect command-line arguments.
	experimentPathFlag := flag.String("experimentPath", "", "Specify the file system location of the directory containing the metric files for one experiment.")
	numMsgsToCalcFlag := flag.Int("numMsgsToCalc", 20, "Calculate statistics for this number of measured messages.")
	flag.Parse()

	numMsgsToCalc := int64(*numMsgsToCalcFlag)

	experimentPath, err := filepath.Abs(*experimentPathFlag)
	if err != nil {
		fmt.Printf("Error converting metrics path '%s' into absolute path: %v\n", *experimentPathFlag, err)
		os.Exit(1)
	}

	// Prepare storage space in experiment struct.

	foldersZenoClients1000, err := ioutil.ReadDir(filepath.Join(experimentPath, "zeno", "clients-1000"))
	if err != nil {

		if !strings.Contains(err.Error(), "no such file or directory") {
			fmt.Printf("Failed to retrieve all folders for zeno's 1000 clients setting: %v\n", err)
			os.Exit(1)
		}
	}

	foldersZenoClients2000, err := ioutil.ReadDir(filepath.Join(experimentPath, "zeno", "clients-2000"))
	if err != nil {

		if !strings.Contains(err.Error(), "no such file or directory") {
			fmt.Printf("Failed to retrieve all folders for zeno's 2000 clients setting: %v\n", err)
			os.Exit(1)
		}
	}

	foldersZenoClients3000, err := ioutil.ReadDir(filepath.Join(experimentPath, "zeno", "clients-3000"))
	if err != nil {

		if !strings.Contains(err.Error(), "no such file or directory") {
			fmt.Printf("Failed to retrieve all folders for zeno's 3000 clients setting: %v\n", err)
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

	foldersVuvuzelaClients2000, err := ioutil.ReadDir(filepath.Join(experimentPath, "vuvuzela", "clients-2000"))
	if err != nil {

		if !strings.Contains(err.Error(), "no such file or directory") {
			fmt.Printf("Failed to retrieve all folders for Vuvuzela's 2000 clients setting: %v\n", err)
			os.Exit(1)
		}
	}

	foldersVuvuzelaClients3000, err := ioutil.ReadDir(filepath.Join(experimentPath, "vuvuzela", "clients-3000"))
	if err != nil {

		if !strings.Contains(err.Error(), "no such file or directory") {
			fmt.Printf("Failed to retrieve all folders for Vuvuzela's 3000 clients setting: %v\n", err)
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

	foldersPungClients2000, err := ioutil.ReadDir(filepath.Join(experimentPath, "pung", "clients-2000"))
	if err != nil {

		if !strings.Contains(err.Error(), "no such file or directory") {
			fmt.Printf("Failed to retrieve all folders for Pung's 2000 clients setting: %v\n", err)
			os.Exit(1)
		}
	}

	foldersPungClients3000, err := ioutil.ReadDir(filepath.Join(experimentPath, "pung", "clients-3000"))
	if err != nil {

		if !strings.Contains(err.Error(), "no such file or directory") {
			fmt.Printf("Failed to retrieve all folders for Pung's 3000 clients setting: %v\n", err)
			os.Exit(1)
		}
	}

	experiment := &Experiment{
		ZenoClients1000:     &Setting{Runs: make([]*Run, 0, len(foldersZenoClients1000))},
		ZenoClients2000:     &Setting{Runs: make([]*Run, 0, len(foldersZenoClients2000))},
		ZenoClients3000:     &Setting{Runs: make([]*Run, 0, len(foldersZenoClients3000))},
		VuvuzelaClients1000: &Setting{Runs: make([]*Run, 0, len(foldersVuvuzelaClients1000))},
		VuvuzelaClients2000: &Setting{Runs: make([]*Run, 0, len(foldersVuvuzelaClients2000))},
		VuvuzelaClients3000: &Setting{Runs: make([]*Run, 0, len(foldersVuvuzelaClients3000))},
		PungClients1000:     &Setting{Runs: make([]*Run, 0, len(foldersPungClients1000))},
		PungClients2000:     &Setting{Runs: make([]*Run, 0, len(foldersPungClients2000))},
		PungClients3000:     &Setting{Runs: make([]*Run, 0, len(foldersPungClients3000))},
	}

	// Append each run to internal structures.
	// Subsequently, write out all metrics files.

	for i := range foldersZenoClients1000 {

		if foldersZenoClients1000[i].IsDir() {
			experiment.ZenoClients1000.AppendRun(filepath.Join(experimentPath, "zeno", "clients-1000",
				foldersZenoClients1000[i].Name()), 21, 1000, numMsgsToCalc)
		}
	}

	if len(experiment.ZenoClients1000.Runs) > 0 {

		err = experiment.ZenoClients1000.MetricsToFiles(filepath.Join(experimentPath, "zeno", "clients-1000"))
		if err != nil {
			fmt.Printf("Failed to store calculated statistics for zeno's setting of 1000 clients to files: %v\n", err)
			os.Exit(1)
		}
	}

	for i := range foldersZenoClients2000 {

		if foldersZenoClients2000[i].IsDir() {
			experiment.ZenoClients2000.AppendRun(filepath.Join(experimentPath, "zeno", "clients-2000",
				foldersZenoClients2000[i].Name()), 21, 2000, numMsgsToCalc)
		}
	}

	if len(experiment.ZenoClients2000.Runs) > 0 {

		err = experiment.ZenoClients2000.MetricsToFiles(filepath.Join(experimentPath, "zeno", "clients-2000"))
		if err != nil {
			fmt.Printf("Failed to store calculated statistics for zeno's setting of 2000 clients to files: %v\n", err)
			os.Exit(1)
		}
	}

	for i := range foldersZenoClients3000 {

		if foldersZenoClients3000[i].IsDir() {
			experiment.ZenoClients3000.AppendRun(filepath.Join(experimentPath, "zeno", "clients-3000",
				foldersZenoClients3000[i].Name()), 21, 3000, numMsgsToCalc)
		}
	}

	if len(experiment.ZenoClients3000.Runs) > 0 {

		err = experiment.ZenoClients3000.MetricsToFiles(filepath.Join(experimentPath, "zeno", "clients-3000"))
		if err != nil {
			fmt.Printf("Failed to store calculated statistics for zeno's setting of 3000 clients to files: %v\n", err)
			os.Exit(1)
		}
	}

	for i := range foldersVuvuzelaClients1000 {

		if foldersVuvuzelaClients1000[i].IsDir() {
			experiment.VuvuzelaClients1000.AppendRun(filepath.Join(experimentPath, "vuvuzela", "clients-1000",
				foldersVuvuzelaClients1000[i].Name()), 4, 1000, numMsgsToCalc)
		}
	}

	if len(experiment.VuvuzelaClients1000.Runs) > 0 {

		err = experiment.VuvuzelaClients1000.MetricsToFiles(filepath.Join(experimentPath, "vuvuzela", "clients-1000"))
		if err != nil {
			fmt.Printf("Failed to store calculated statistics for Vuvuzela's setting of 1000 clients to files: %v\n", err)
			os.Exit(1)
		}
	}

	for i := range foldersVuvuzelaClients2000 {

		if foldersVuvuzelaClients2000[i].IsDir() {
			experiment.VuvuzelaClients2000.AppendRun(filepath.Join(experimentPath, "vuvuzela", "clients-2000",
				foldersVuvuzelaClients2000[i].Name()), 4, 2000, numMsgsToCalc)
		}
	}

	if len(experiment.VuvuzelaClients2000.Runs) > 0 {

		err = experiment.VuvuzelaClients2000.MetricsToFiles(filepath.Join(experimentPath, "vuvuzela", "clients-2000"))
		if err != nil {
			fmt.Printf("Failed to store calculated statistics for Vuvuzela's setting of 2000 clients to files: %v\n", err)
			os.Exit(1)
		}
	}

	for i := range foldersVuvuzelaClients3000 {

		if foldersVuvuzelaClients3000[i].IsDir() {
			experiment.VuvuzelaClients3000.AppendRun(filepath.Join(experimentPath, "vuvuzela", "clients-3000",
				foldersVuvuzelaClients3000[i].Name()), 4, 3000, numMsgsToCalc)
		}
	}

	if len(experiment.VuvuzelaClients3000.Runs) > 0 {

		err = experiment.VuvuzelaClients3000.MetricsToFiles(filepath.Join(experimentPath, "vuvuzela", "clients-3000"))
		if err != nil {
			fmt.Printf("Failed to store calculated statistics for Vuvuzela's setting of 3000 clients to files: %v\n", err)
			os.Exit(1)
		}
	}

	for i := range foldersPungClients1000 {

		if foldersPungClients1000[i].IsDir() {
			experiment.PungClients1000.AppendRun(filepath.Join(experimentPath, "pung", "clients-1000",
				foldersPungClients1000[i].Name()), 1, 1000, numMsgsToCalc)
		}
	}

	if len(experiment.PungClients1000.Runs) > 0 {

		err = experiment.PungClients1000.MetricsToFiles(filepath.Join(experimentPath, "pung", "clients-1000"))
		if err != nil {
			fmt.Printf("Failed to store calculated statistics for Pung's setting of 1000 clients to files: %v\n", err)
			os.Exit(1)
		}
	}

	for i := range foldersPungClients2000 {

		if foldersPungClients2000[i].IsDir() {
			experiment.PungClients2000.AppendRun(filepath.Join(experimentPath, "pung", "clients-2000",
				foldersPungClients2000[i].Name()), 1, 2000, numMsgsToCalc)
		}
	}

	if len(experiment.PungClients2000.Runs) > 0 {

		err = experiment.PungClients2000.MetricsToFiles(filepath.Join(experimentPath, "pung", "clients-2000"))
		if err != nil {
			fmt.Printf("Failed to store calculated statistics for Pung's setting of 2000 clients to files: %v\n", err)
			os.Exit(1)
		}
	}

	for i := range foldersPungClients3000 {

		if foldersPungClients3000[i].IsDir() {
			experiment.PungClients3000.AppendRun(filepath.Join(experimentPath, "pung", "clients-3000",
				foldersPungClients3000[i].Name()), 1, 3000, numMsgsToCalc)
		}
	}

	if len(experiment.PungClients3000.Runs) > 0 {

		err = experiment.PungClients3000.MetricsToFiles(filepath.Join(experimentPath, "pung", "clients-3000"))
		if err != nil {
			fmt.Printf("Failed to store calculated statistics for Pung's setting of 3000 clients to files: %v\n", err)
			os.Exit(1)
		}
	}
}
