package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Worker comprises all flags and values
// required for the metrics collector of our
// ACS evaluation to work correctly.
type Worker struct {
	shutdownChan   chan struct{}
	System         string
	IsClient       bool
	IsServer       bool
	MetricsChan    chan string
	SentBytesFile  *os.File
	RecvdBytesFile *os.File
	LoadFile       *os.File
	MemFile        *os.File
	SendTimeFile   *os.File
	RecvTimeFile   *os.File
	PoolSizesFile  *os.File
	PipeReader     *bufio.Reader
}

func (col *Worker) prepareMetricsFiles(metricsPath string) error {

	var err error

	// Attempt to create file for sent bytes metric.
	col.SentBytesFile, err = os.OpenFile(filepath.Join(metricsPath, "traffic_outgoing.evaluation"), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}

	// Attempt to create file for received bytes metric.
	col.RecvdBytesFile, err = os.OpenFile(filepath.Join(metricsPath, "traffic_incoming.evaluation"), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}

	// Attempt to create file for load metrics.
	col.LoadFile, err = os.OpenFile(filepath.Join(metricsPath, "load_unixnano.evaluation"), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}

	// Attempt to create file for memory usage.
	col.MemFile, err = os.OpenFile(filepath.Join(metricsPath, "mem_unixnano.evaluation"), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}

	if col.IsClient {

		// Attempt to create file for send times metric.
		col.SendTimeFile, err = os.OpenFile(filepath.Join(metricsPath, "send_unixnano.evaluation"), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
		if err != nil {
			return err
		}

		// Attempt to create file for receive times metric.
		col.RecvTimeFile, err = os.OpenFile(filepath.Join(metricsPath, "recv_unixnano.evaluation"), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
		if err != nil {
			return err
		}

	} else if col.System == "zeno" && col.IsServer {

		// Attempt to create file for pool sizes metrics.
		col.PoolSizesFile, err = os.OpenFile(filepath.Join(metricsPath, "pool-sizes_round.evaluation"), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func (col *Worker) collectSystemMetrics() {

	// Receive tick every second.
	secTicker := time.NewTicker(time.Second)

	for {

		select {

		// Leave function once a close
		// signal is received.
		case <-col.shutdownChan:
			return

		case <-secTicker.C:

			sentBytes := "n/a\n"
			recvdBytes := "n/a\n"
			load := "n/a\n"
			mem := "n/a\n"

			// Prepare various system metric collection commands.
			cmdSent := exec.Command("iptables", "-nvxL", "OUTPUT")
			cmdRecvd := exec.Command("iptables", "-nvxL", "INPUT")
			cmdLoad := exec.Command("mpstat")
			cmdMem := exec.Command("head", "-3", "/proc/meminfo")

			searchSent := "dpt:33000"
			searchRecvd := "dpt:33000"

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

				// If this is a pung server, we cannot look for
				// outgoing traffic to a specific port, but only
				// for packets sent via already established
				// connections by clients.
				if col.System == "pung" && col.IsServer {
					searchSent = "ESTABLISHED"
				}

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
				if col.System == "pung" && col.IsClient {
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
			load = fmt.Sprintf("%d usr:%s nice:%s sys:%s iowait:%s idle:%s\n", now, outLoadParts[2], outLoadParts[3], outLoadParts[4], outLoadParts[5], outLoadParts[11])

			// Extract memory usage values.
			outMemLines := strings.Split(strings.TrimSpace(outMem), "\n")
			memTotal := strings.Fields(outMemLines[0])
			memAvail := strings.Fields(outMemLines[2])
			mem = fmt.Sprintf("%d totalKB:%s availKB:%s\n", now, memTotal[1], memAvail[1])

			// Write all values to their respective
			// metrics files on disk.
			fmt.Fprint(col.SentBytesFile, sentBytes)
			_ = col.SentBytesFile.Sync()

			fmt.Fprint(col.RecvdBytesFile, recvdBytes)
			_ = col.RecvdBytesFile.Sync()

			fmt.Fprint(col.LoadFile, load)
			_ = col.LoadFile.Sync()

			fmt.Fprint(col.MemFile, mem)
			_ = col.MemFile.Sync()
		}
	}
}

func (col *Worker) collectTimingMetrics() {

	// Receive timing values come in two parts,
	// first the time for the subsequently
	// transmitted message ID. Hold buffer.
	recvTime := ""

	for metric := range col.MetricsChan {

		// Split at semicolon.
		metricParts := strings.Split(metric, ";")

		if metricParts[0] == "send" {

			// Write to file and sync to stable storage.
			fmt.Fprint(col.SendTimeFile, metricParts[1])
			_ = col.SendTimeFile.Sync()

		} else if metricParts[0] == "recv" {

			if recvTime == "" {

				// Stash receive time until message ID
				// associated with it is sent next.
				recvTime = strings.TrimSpace(metricParts[1])

			} else {

				// Write to file and sync to stable storage.
				fmt.Fprintf(col.RecvTimeFile, "%s %s\n", recvTime, strings.TrimSpace(metricParts[1]))
				_ = col.RecvTimeFile.Sync()

				// Reset buffer for receive timestamp.
				recvTime = ""
			}
		}
	}
}

func (col *Worker) collectPoolSizesMetrics() {

	for metric := range col.MetricsChan {

		// Write to file and sync to stable storage.
		fmt.Fprint(col.PoolSizesFile, metric)
		_ = col.PoolSizesFile.Sync()
	}
}

func main() {

	// Allow some command-line arguments.
	systemFlag := flag.String("system", "", "Specify system that is being evaluated ('zeno', 'vuvuzela', or 'pung').")
	isClientFlag := flag.Bool("client", false, "Append this flag if the worker is gathering metrics for a client that is being evaluated.")
	isServerFlag := flag.Bool("server", false, "Append this flag if the worker is gathering metrics for a server that is being evaluated.")
	pipeNameFlag := flag.String("pipe", "/tmp/collect", "Specify the named pipe to use for IPC with system being evaluated.")
	metricsPathFlag := flag.String("metricsPath", "./", "Specify the file system folder where the various metric files generated here should be placed.")
	flag.Parse()

	// System flag has to be one of three values.
	if *systemFlag != "zeno" && *systemFlag != "vuvuzela" && *systemFlag != "pung" {
		fmt.Printf("Flag '-system' requires one of the three values: 'zeno', 'vuvuzela', or 'pung'.")
		os.Exit(1)
	}

	// Enforce either client or server designation.
	if *isClientFlag == *isServerFlag {
		fmt.Printf("Please either specify '-client' or '-server'.\n")
		os.Exit(1)
	}

	metricsPath, err := filepath.Abs(*metricsPathFlag)
	if err != nil {
		fmt.Printf("Error converting metrics path '%s' into absolute path: %v\n", *metricsPathFlag, err)
		os.Exit(1)
	}

	// Initialize worker struct.
	col := &Worker{
		shutdownChan: make(chan struct{}),
		System:       strings.ToLower(*systemFlag),
		IsClient:     *isClientFlag,
		IsServer:     *isServerFlag,
		MetricsChan:  make(chan string, 100),
	}

	// Prepare the various metrics files.
	err = col.prepareMetricsFiles(metricsPath)
	if err != nil {
		fmt.Printf("Unable to prepare files for collecting metrics: %v\n", err)
		os.Exit(1)
	}

	if col.IsClient {

		// Spawn background process writing timing
		// values into metrics files.
		go col.collectTimingMetrics()

	} else if col.System == "zeno" && col.IsServer {

		// Spawn background process writing message
		// pool sizes to into metrics file.
		go col.collectPoolSizesMetrics()
	}

	// Spawn background process writing sent and
	// received bytes values to file every second.
	go col.collectSystemMetrics()

	defer func(col *Worker) {

		// Clean up when we leave this program.

		close(col.MetricsChan)
		col.shutdownChan <- struct{}{}

		col.SentBytesFile.Close()
		col.RecvdBytesFile.Close()
		col.LoadFile.Close()
		col.MemFile.Close()

		if col.IsClient {
			col.SendTimeFile.Close()
			col.RecvTimeFile.Close()
		} else if col.IsServer {
			col.PoolSizesFile.Close()
		}
	}(col)

	// Open named pipe for receiving metrics from
	// system under evaluation.
	pipe, err := os.OpenFile(*pipeNameFlag, os.O_RDONLY, 0600)
	if err != nil {
		fmt.Printf("Unable to open named pipe for passing metrics: %v\n", err)
		os.Exit(1)
	}
	col.PipeReader = bufio.NewReader(pipe)

	// Read next metric line from named pipe.
	metric, err := col.PipeReader.ReadString('\n')
	if err != nil {
		fmt.Printf("Failed reading from named pipe: %v\n", err)
		return
	}

	for metric != "done\n" {

		// Off-load metric line to file writer.
		col.MetricsChan <- metric

		// Read next metric line from named pipe.
		metric, err = col.PipeReader.ReadString('\n')
		if err != nil {
			fmt.Printf("Failed reading from named pipe: %v\n", err)
			return
		}
	}
}
