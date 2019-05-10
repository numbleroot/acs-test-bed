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

// Collector comprises all flags and values
// required for the metrics collector of our
// ACS evaluation to work correctly.
type Collector struct {
	shutdownChan   chan struct{}
	IsClient       bool
	IsMix          bool
	MetricsChan    chan string
	SentBytesFile  *os.File
	RecvdBytesFile *os.File
	SendTimeFile   *os.File
	RecvTimeFile   *os.File
	TimePipeReader *bufio.Reader
}

func (col *Collector) prepareMetricsFiles(metricsPath string, pipeName string) error {

	// Attempt to create file for sent bytes metric.
	SentBytesFile, err := os.OpenFile(filepath.Join(metricsPath, "traffic_incoming.evaluation"), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}

	// Attempt to create file for received bytes metric.
	RecvdBytesFile, err := os.OpenFile(filepath.Join(metricsPath, "traffic_outgoing.evaluation"), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}

	// Attempt to create file for send times metric.
	sendTimeFile, err := os.OpenFile(filepath.Join(metricsPath, "send_unixnano.evaluation"), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}

	// Attempt to create file for receive times metric.
	recvTimeFile, err := os.OpenFile(filepath.Join(metricsPath, "recv_unixnano.evaluation"), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}

	// Open named pipe for receiving timing values.
	pipe, err := os.OpenFile(pipeName, os.O_RDONLY, 0600)
	if err != nil {
		return err
	}

	col.SentBytesFile = SentBytesFile
	col.RecvdBytesFile = RecvdBytesFile
	col.SendTimeFile = sendTimeFile
	col.RecvTimeFile = recvTimeFile
	col.TimePipeReader = bufio.NewReader(pipe)

	return nil
}

func (col *Collector) writeToTrafficFiles() {

	// Prepare commands to extract current number
	// of sent and received bytes.
	cmdSent := exec.Command("iptables", "-nvxL", "OUTPUT")
	cmdRecvd := exec.Command("iptables", "-nvxL", "INPUT")

	// Receive tick every second.
	secTicker := time.NewTicker(time.Second)

	for {

		select {

		// Leave function once a close
		// signal is received.
		case <-col.shutdownChan:
			return

		case <-secTicker.C:

			sentBytes := "n/a"
			recvdBytes := "n/a"

			// Obtain current timestamp.
			now := time.Now().UnixNano()

			// Execute command to find sent bytes value.
			outSentRaw, err := cmdSent.CombinedOutput()
			if err != nil {
				fmt.Printf("Collecting sent bytes metric failed (error: %v):\n%s", err, outSentRaw)
			}
			outSent := string(outSentRaw)

			// Execute command to find received bytes value.
			outRecvdRaw, err := cmdRecvd.CombinedOutput()
			if err != nil {
				fmt.Printf("Collecting received bytes metric failed (error: %v):\n%s", err, outRecvdRaw)
			}
			outRecvd := string(outRecvdRaw)

			outSentLines := strings.Split(outSent, "\n")
			for i := range outSentLines {

				if strings.Contains(outSentLines[i], "dpt:33000") {

					// Split at one or more whitespace characters.
					// The bytes value is the second.
					outSentParts := strings.Fields(outSentLines[i])
					sentBytes = outSentParts[1]
				}
			}

			outRecvdLines := strings.Split(outRecvd, "\n")
			for i := range outRecvdLines {

				if strings.Contains(outRecvdLines[i], "dpt:33000") {

					// Split at one or more whitespace characters.
					// The bytes value is the second.
					outRecvdParts := strings.Fields(outRecvdLines[i])
					recvdBytes = outRecvdParts[1]
				}
			}

			// Write both values to their respective
			// metrics files on disk.
			fmt.Fprintf(col.SentBytesFile, "%d %s\n", now, sentBytes)
			fmt.Fprintf(col.RecvdBytesFile, "%d %s\n", now, recvdBytes)
			_ = col.SentBytesFile.Sync()
			_ = col.RecvdBytesFile.Sync()
		}
	}
}

func (col *Collector) writeToTimingFiles() {

	for metric := range col.MetricsChan {

		// Split at semicolon.
		metricParts := strings.Split(metric, ";")

		if metricParts[0] == "send" {

			// Write to file and sync to stable storage.
			fmt.Fprint(col.SendTimeFile, metricParts[1])
			_ = col.SendTimeFile.Sync()

		} else if metricParts[0] == "recv" {

			// Write to file and sync to stable storage.
			fmt.Fprint(col.RecvTimeFile, metricParts[1])
			_ = col.RecvTimeFile.Sync()
		}
	}
}

func main() {

	// Allow some command-line arguments.
	isClientFlag := flag.Bool("client", false, "Append this flag if the collector is gathering metrics for a client that is being evaluated.")
	isMixFlag := flag.Bool("mix", false, "Append this flag if the collector is gathering metrics for a mix that is being evaluated.")
	pipeNameFlag := flag.String("pipe", "/tmp/collect", "Specify the named pipe to use for IPC with system being evaluated.")
	metricsPathFlag := flag.String("metricsPath", "./", "Specify the file system folder where the various metric files generated here should be placed.")
	flag.Parse()

	// Enforce either client or mix designation.
	if *isClientFlag == *isMixFlag {
		fmt.Printf("Please either specify '-client' or '-mix'.\n")
		os.Exit(1)
	}

	metricsPath, err := filepath.Abs(*metricsPathFlag)
	if err != nil {
		fmt.Printf("Error converting results path '%s' into absolute path: %v\n", *metricsPathFlag, err)
		os.Exit(1)
	}

	// Initialize collector struct.
	col := &Collector{
		shutdownChan: make(chan struct{}),
		IsClient:     *isClientFlag,
		IsMix:        *isMixFlag,
		MetricsChan:  make(chan string, 100),
	}

	// Prepare the various metrics files.
	err = col.prepareMetricsFiles(metricsPath, *pipeNameFlag)
	if err != nil {
		fmt.Printf("Unable to prepare files for collecting metrics: %v\n", err)
		os.Exit(1)
	}

	// Spawn background process writing timing
	// values into metrics files.
	go col.writeToTimingFiles()

	// Spawn background process writing sent and
	// received bytes values to file every second.
	go col.writeToTrafficFiles()

	// Read next metric line from named pipe and
	// clean it up.
	metric, err := col.TimePipeReader.ReadString('\n')
	if err != nil {
		fmt.Printf("Failed reading from named pipe: %v\n", err)
	}

	for metric != "done" {

		// Off-load metric line to file writer.
		col.MetricsChan <- metric

		// Read next metric line from named pipe and
		// clean it up.
		metric, err = col.TimePipeReader.ReadString('\n')
		if err != nil {
			fmt.Printf("Failed reading from named pipe: %v\n", err)
			continue
		}
	}

	// Node being evaluated signaled that the
	// evaluation is completed, signal internally
	// via shutdown channel and wait for response.
	col.shutdownChan <- struct{}{}
}
