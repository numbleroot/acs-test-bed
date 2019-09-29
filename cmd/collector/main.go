package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Collector comprises all flags and values
// required for the metrics collector of our
// ACS evaluation to work correctly.
type Collector struct {
	shutdownChan chan struct{}
	System       string
	TypeOfNode   string
	MetricsPath  string
}

func init() {

	// Enable TLS 1.3.
	if os.Getenv("GODEBUG") == "" {
		os.Setenv("GODEBUG", "tls13=1")
	} else {
		os.Setenv("GODEBUG", fmt.Sprintf("%s,tls13=1", os.Getenv("GODEBUG")))
	}
}

func (col *Collector) collectSystemMetrics() {

	// Attempt to create file for sent bytes metric.
	sentBytesFile, err := os.OpenFile(filepath.Join(col.MetricsPath, "traffic_outgoing.evaluation"),
		(os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		fmt.Printf("Unable to open or create 'traffic_outgoing.evaluation': %v\n", err)
		os.Exit(1)
	}

	// Attempt to create file for received bytes metric.
	recvdBytesFile, err := os.OpenFile(filepath.Join(col.MetricsPath, "traffic_incoming.evaluation"),
		(os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		fmt.Printf("Unable to open or create 'traffic_incoming.evaluation': %v\n", err)
		os.Exit(1)
	}

	// Attempt to create file for load metrics.
	loadFile, err := os.OpenFile(filepath.Join(col.MetricsPath, "load_unixnano.evaluation"),
		(os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		fmt.Printf("Unable to open or create 'load_unixnano.evaluation': %v\n", err)
		os.Exit(1)
	}

	// Attempt to create file for memory usage.
	memFile, err := os.OpenFile(filepath.Join(col.MetricsPath, "mem_unixnano.evaluation"),
		(os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		fmt.Printf("Unable to open or create 'mem_unixnano.evaluation': %v\n", err)
		os.Exit(1)
	}

	// Receive tick every second.
	secTicker := time.NewTicker(time.Second)

	for {

		select {

		// Leave function once a close
		// signal is received.
		case <-col.shutdownChan:

			sentBytesFile.Close()
			recvdBytesFile.Close()
			loadFile.Close()
			memFile.Close()

			return

		case <-secTicker.C:

			sentBytes := "n/a\n"
			recvdBytes := "n/a\n"
			load := "n/a\n"
			mem := "n/a\n"

			// Prepare various system metric collection commands.
			cmdSent := exec.Command("iptables", "-t", "filter", "-nvx", "-L", "OUTPUT")
			cmdRecvd := exec.Command("iptables", "-t", "filter", "-nvx", "-L", "INPUT")
			cmdLoad := exec.Command("mpstat")
			cmdMem := exec.Command("head", "-3", "/proc/meminfo")

			// TODO: If zeno, add search for port 44001.
			searchSent := "spt:33001"
			searchRecvd := "dpt:33001"

			// Obtain current timestamp.
			now := time.Now().Unix()

			// Execute command to find sent bytes value.
			outSentRaw, err := cmdSent.CombinedOutput()
			if err != nil {
				fmt.Printf("Collecting sent bytes metric failed (error: %v):\n%s\n", err, outSentRaw)
			}
			outSent := string(outSentRaw)

			// Execute command to find received bytes value.
			outRecvdRaw, err := cmdRecvd.CombinedOutput()
			if err != nil {
				fmt.Printf("Collecting received bytes metric failed (error: %v):\n%s\n", err, outRecvdRaw)
			}
			outRecvd := string(outRecvdRaw)

			// Execute command to find current load.
			outLoadRaw, err := cmdLoad.CombinedOutput()
			if err != nil {
				fmt.Printf("Collecting load metric failed (error: %v):\n%s\n", err, outLoadRaw)
			}
			outLoad := string(outLoadRaw)

			// Execute command to find memory usage.
			outMemRaw, err := cmdMem.CombinedOutput()
			if err != nil {
				fmt.Printf("Collecting memory usage failed (error: %v):\n%s\n", err, outMemRaw)
			}
			outMem := string(outMemRaw)

			outSentLines := strings.Split(strings.TrimSpace(outSent), "\n")
			for i := range outSentLines {

				if strings.Contains(outSentLines[i], searchSent) {

					// Split at one or more whitespace characters.
					// The bytes value is the second one.
					outSentParts := strings.Fields(outSentLines[i])
					sentBytes = fmt.Sprintf("%d %s\n", now, outSentParts[1])
				}
			}

			outRecvdLines := strings.Split(strings.TrimSpace(outRecvd), "\n")
			for i := range outRecvdLines {

				// If this is a pung client, we cannot look for
				// incoming traffic to a specific port, but only
				// for packets sent via connections established
				// by this client.
				if col.System == "pung" && col.TypeOfNode == "client" {
					searchRecvd = "ESTABLISHED"
				}

				if strings.Contains(outRecvdLines[i], searchRecvd) {

					// Split at one or more whitespace characters.
					// The bytes value is the second one.
					outRecvdParts := strings.Fields(outRecvdLines[i])
					recvdBytes = fmt.Sprintf("%d %s\n", now, outRecvdParts[1])
				}
			}

			// Extract the interesting load metrics.
			outLoadLines := strings.Split(strings.TrimSpace(outLoad), "\n")
			outLoadParts := strings.Fields(outLoadLines[(len(outLoadLines) - 1)])
			load = fmt.Sprintf("%d usr:%s nice:%s sys:%s iowait:%s idle:%s\n", now,
				outLoadParts[2], outLoadParts[3], outLoadParts[4], outLoadParts[5], outLoadParts[11])

			// Extract memory usage values.
			outMemLines := strings.Split(strings.TrimSpace(outMem), "\n")
			memTotal := strings.Fields(outMemLines[0])
			memAvail := strings.Fields(outMemLines[2])
			mem = fmt.Sprintf("%d totalKB:%s availKB:%s\n", now, memTotal[1], memAvail[1])

			// Write all values to their respective
			// metrics files on disk.
			fmt.Fprint(sentBytesFile, sentBytes)
			_ = sentBytesFile.Sync()

			fmt.Fprint(recvdBytesFile, recvdBytes)
			_ = recvdBytesFile.Sync()

			fmt.Fprint(loadFile, load)
			_ = loadFile.Sync()

			fmt.Fprint(memFile, mem)
			_ = memFile.Sync()
		}
	}
}

func (col *Collector) collectTimingMetrics(wg *sync.WaitGroup, pipePath string, clientPath string) {

	defer wg.Done()

	// Attempt to create file for send times metric.
	sendTimeFile, err := os.OpenFile(filepath.Join(clientPath, "send_unixnano.evaluation"),
		(os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		fmt.Printf("Unable to open or create 'send_unixnano.evaluation' in '%s': %v\n", clientPath, err)
		os.Exit(1)
	}

	// Attempt to create file for receive times metric.
	recvTimeFile, err := os.OpenFile(filepath.Join(clientPath, "recv_unixnano.evaluation"),
		(os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		fmt.Printf("Unable to open or create 'recv_unixnano.evaluation' in '%s': %v\n", clientPath, err)
		os.Exit(1)
	}

	// Prepare channel to send metrics over.
	metricsChan := make(chan string, 100)
	go col.writeTimingMetrics(metricsChan, sendTimeFile, recvTimeFile)

	// Open named pipe for receiving metrics from
	// system under evaluation.
	pipe, err := os.OpenFile(pipePath, os.O_RDONLY, 0600)
	if err != nil {
		fmt.Printf("Unable to open named pipe '%s' for passing metrics: %v\n", pipePath, err)
		os.Exit(1)
	}
	pipeReader := bufio.NewReader(pipe)

	// Read next metric line from named pipe.
	metric, err := pipeReader.ReadString('\n')
	if err != nil {
		fmt.Printf("Failed reading from named pipe '%s': %v\n", pipePath, err)
		return
	}

	for metric != "done\n" {

		// Off-load metric line to file writer.
		metricsChan <- metric

		// Read next metric line from named pipe.
		metric, err = pipeReader.ReadString('\n')
		if err != nil {
			fmt.Printf("Failed reading from named pipe '%s': %v\n", pipePath, err)
			return
		}
	}

	// Clean up.
	close(metricsChan)
	sendTimeFile.Close()
	recvTimeFile.Close()
}

func (col *Collector) writeTimingMetrics(metricsChan chan string, sendTimeFile *os.File, recvTimeFile *os.File) {

	// Receive timing values come in two parts,
	// first the time for the subsequently
	// transmitted message ID. Hold buffer.
	recvTime := ""

	for metric := range metricsChan {

		// Split at semicolon.
		metricParts := strings.Split(metric, ";")

		if metricParts[0] == "send" {

			// Write to file and sync to stable storage.
			fmt.Fprint(sendTimeFile, metricParts[1])
			_ = sendTimeFile.Sync()

		} else if metricParts[0] == "recv" {

			if recvTime == "" {

				// Stash receive time until message ID
				// associated with it is sent next.
				recvTime = strings.TrimSpace(metricParts[1])

			} else {

				// Write to file and sync to stable storage.
				fmt.Fprintf(recvTimeFile, "%s %s\n", recvTime, strings.TrimSpace(metricParts[1]))
				_ = recvTimeFile.Sync()

				// Reset buffer for receive timestamp.
				recvTime = ""
			}
		}
	}
}

func (col *Collector) collectPoolSizesMetrics(wg *sync.WaitGroup, pipePath string, clientPath string) {

	defer wg.Done()

	// Attempt to create file for pool sizes metrics.
	poolSizesFile, err := os.OpenFile(filepath.Join(clientPath, "pool-sizes_round.evaluation"),
		(os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		fmt.Printf("Unable to open or create 'pool-sizes_round.evaluation' in '%s': %v\n", clientPath, err)
		os.Exit(1)
	}

	// Prepare channel to send metrics over.
	metricsChan := make(chan string, 100)
	go col.writePoolSizesMetrics(metricsChan, poolSizesFile)

	// Open named pipe for receiving metrics from
	// system under evaluation.
	pipe, err := os.OpenFile(pipePath, os.O_RDONLY, 0600)
	if err != nil {
		fmt.Printf("Unable to open named pipe '%s' for passing metrics: %v\n", pipePath, err)
		os.Exit(1)
	}
	pipeReader := bufio.NewReader(pipe)

	// Read next metric line from named pipe.
	metric, err := pipeReader.ReadString('\n')
	if err != nil {
		fmt.Printf("Failed reading from named pipe '%s': %v\n", pipePath, err)
		return
	}

	for metric != "done\n" {

		// Off-load metric line to file writer.
		metricsChan <- metric

		// Read next metric line from named pipe.
		metric, err = pipeReader.ReadString('\n')
		if err != nil {
			fmt.Printf("Failed reading from named pipe '%s': %v\n", pipePath, err)
			return
		}
	}

	// Clean up.
	close(metricsChan)
	poolSizesFile.Close()
}

func (col *Collector) writePoolSizesMetrics(metricsChan chan string, poolSizesFile *os.File) {

	for metric := range metricsChan {

		// Write to file and sync to stable storage.
		fmt.Fprint(poolSizesFile, metric)
		_ = poolSizesFile.Sync()
	}
}

func main() {

	// Allow some command-line arguments.
	systemFlag := flag.String("system", "", "Specify system that is being evaluated ('zeno', 'vuvuzela', or 'pung').")
	typeOfNodeFlag := flag.String("typeOfNode", "", "Specify the type of node being evaluated ('client', 'server', 'coordinator').")
	metricsPathFlag := flag.String("metricsPath", "./", "Specify the file system folder where the various metric files generated here should be placed.")
	client01PathFlag := flag.String("client01Path", "./client-00001", "Specify file system folder in which to store metrics for client 01.")
	pipe01Flag := flag.String("pipe01", "/tmp/collect01", "Specify named pipe 01 to use for metrics IPC.")
	client02PathFlag := flag.String("client02Path", "./client-00002", "Specify file system folder in which to store metrics for client 02.")
	pipe02Flag := flag.String("pipe02", "/tmp/collect02", "Specify named pipe 02 to use for metrics IPC.")
	client03PathFlag := flag.String("client03Path", "./client-00003", "Specify file system folder in which to store metrics for client 03.")
	pipe03Flag := flag.String("pipe03", "/tmp/collect03", "Specify named pipe 03 to use for metrics IPC.")
	client04PathFlag := flag.String("client04Path", "./client-00004", "Specify file system folder in which to store metrics for client 04.")
	pipe04Flag := flag.String("pipe04", "/tmp/collect04", "Specify named pipe 04 to use for metrics IPC.")
	client05PathFlag := flag.String("client05Path", "./client-00005", "Specify file system folder in which to store metrics for client 05.")
	pipe05Flag := flag.String("pipe05", "/tmp/collect05", "Specify named pipe 05 to use for metrics IPC.")
	client06PathFlag := flag.String("client06Path", "./client-00006", "Specify file system folder in which to store metrics for client 06.")
	pipe06Flag := flag.String("pipe06", "/tmp/collect06", "Specify named pipe 06 to use for metrics IPC.")
	client07PathFlag := flag.String("client07Path", "./client-00007", "Specify file system folder in which to store metrics for client 07.")
	pipe07Flag := flag.String("pipe07", "/tmp/collect07", "Specify named pipe 07 to use for metrics IPC.")
	client08PathFlag := flag.String("client08Path", "./client-00008", "Specify file system folder in which to store metrics for client 08.")
	pipe08Flag := flag.String("pipe08", "/tmp/collect08", "Specify named pipe 08 to use for metrics IPC.")
	client09PathFlag := flag.String("client09Path", "./client-00009", "Specify file system folder in which to store metrics for client 09.")
	pipe09Flag := flag.String("pipe09", "/tmp/collect09", "Specify named pipe 09 to use for metrics IPC.")
	client10PathFlag := flag.String("client10Path", "./client-00010", "Specify file system folder in which to store metrics for client 10.")
	pipe10Flag := flag.String("pipe10", "/tmp/collect10", "Specify named pipe 10 to use for metrics IPC.")
	flag.Parse()

	if *systemFlag != "zeno" && *systemFlag != "vuvuzela" && *systemFlag != "pung" {
		fmt.Printf("Flag '-system' requires one of the three values: 'zeno', 'vuvuzela', or 'pung'.")
		os.Exit(1)
	}

	if *typeOfNodeFlag == "" {
		fmt.Printf("Please either specify '-client', '-server', or '-coordinator'.\n")
		os.Exit(1)
	}

	metricsPath, err := filepath.Abs(*metricsPathFlag)
	if err != nil {
		fmt.Printf("Error converting metrics path '%s' into absolute path: %v\n", *metricsPathFlag, err)
		os.Exit(1)
	}

	// Initialize collector struct.
	col := &Collector{
		shutdownChan: make(chan struct{}),
		System:       strings.ToLower(*systemFlag),
		TypeOfNode:   *typeOfNodeFlag,
		MetricsPath:  metricsPath,
	}

	// Spawn background process writing sent and
	// received bytes values to file every second.
	go col.collectSystemMetrics()

	wg := &sync.WaitGroup{}

	if col.TypeOfNode == "client" {

		wg.Add(10)

		// Spawn background processes writing timing
		// values into metrics files.
		go col.collectTimingMetrics(wg, *pipe01Flag, *client01PathFlag)
		go col.collectTimingMetrics(wg, *pipe02Flag, *client02PathFlag)
		go col.collectTimingMetrics(wg, *pipe03Flag, *client03PathFlag)
		go col.collectTimingMetrics(wg, *pipe04Flag, *client04PathFlag)
		go col.collectTimingMetrics(wg, *pipe05Flag, *client05PathFlag)
		go col.collectTimingMetrics(wg, *pipe06Flag, *client06PathFlag)
		go col.collectTimingMetrics(wg, *pipe07Flag, *client07PathFlag)
		go col.collectTimingMetrics(wg, *pipe08Flag, *client08PathFlag)
		go col.collectTimingMetrics(wg, *pipe09Flag, *client09PathFlag)
		go col.collectTimingMetrics(wg, *pipe10Flag, *client10PathFlag)

	} else {

		wg.Add(1)

		// Spawn background process writing message
		// pool sizes to into metrics file.
		go col.collectPoolSizesMetrics(wg, *pipe01Flag, *client01PathFlag)
	}

	// Wait for all metric collections
	// to complete, then clean up.
	wg.Wait()
	col.shutdownChan <- struct{}{}

	fmt.Printf("Terminating collector, goodbye.\n")
}
