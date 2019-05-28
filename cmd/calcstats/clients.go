package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func (run *Run) AddLatency(runClientsPath string, systemUnderEval string, numMsgsToCalc int64) error {

	err := filepath.Walk(runClientsPath, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if filepath.Base(path) == "send_unixnano.evaluation" {

			partner := ""

			// Ingest supplied send times file.
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			content = bytes.TrimSpace(content)

			// Split file contents into lines.
			lines := strings.Split(string(content), "\n")

			// Prepare latencies state object.
			msgLatencies := make([]*MetricLatency, len(lines))

			for i := range lines {

				// Split line at whitespace characters.
				metric := strings.Fields(lines[i])

				// Convert first element to timestamp.
				timestamp, err := strconv.ParseInt(metric[0], 10, 64)
				if err != nil {
					return err
				}

				// Extract name of partner node.
				if partner == "" {
					partner = strings.Split(metric[1], "=>")[1]
					// fmt.Printf("Partner of '%s': '%s'\n", path, partner)
				}

				// Convert third element to message ID.
				msgID, err := strconv.ParseInt(metric[2], 10, 64)
				if err != nil {
					return err
				}

				// Append to temporary state object.
				msgLatencies[i] = &MetricLatency{
					MsgID:         msgID,
					SendTimestamp: timestamp,
				}
			}

			// Sort state object by message IDs.
			sort.Slice(msgLatencies, func(i, j int) bool {
				return msgLatencies[i].MsgID < msgLatencies[j].MsgID
			})

			/*
				fmt.Printf("[AAA] Found:\n")
				for i := range msgLatencies {
					fmt.Printf("\t%d: %d\n", msgLatencies[i].MsgID, msgLatencies[i].SendTimestamp)
				}
			*/

			if systemUnderEval == "pung" {

				// Reslice list of messages to capture
				// messages [#2, (#2 + #numMsgs)).
				fmt.Printf("[PUNG] Reslicing msgLatencies to [2, %d)\n", (2 + int(numMsgsToCalc)))
				msgLatencies = msgLatencies[2:(2 + int(numMsgsToCalc))]
			}

			/*
				fmt.Printf("[BBB] Found:\n")
				for i := range msgLatencies {
					fmt.Printf("\t%d: %d\n", msgLatencies[i].MsgID, msgLatencies[i].SendTimestamp)
				}
			*/

			// Determine whether we have all required data points.
			if (int64(len(msgLatencies)) < numMsgsToCalc) || (msgLatencies[(numMsgsToCalc-1)].MsgID != numMsgsToCalc) {
				fmt.Printf("Not enough consecutive send message time metrics available (found: %d, want: %d).\n", msgLatencies[(numMsgsToCalc-1)].MsgID, numMsgsToCalc)
				os.Exit(1)
			}

			// Find partner's receive time file.
			candidates, err := filepath.Glob(fmt.Sprintf("%s/%s_*/recv_unixnano.evaluation", runClientsPath, partner))
			if err != nil {
				return err
			}

			if len(candidates) != 1 {
				fmt.Printf("Client at '%s' did not have unique conversation partner.\n", path)
				os.Exit(1)
			}

			// Ingest partner's receive times file.
			content, err = ioutil.ReadFile(candidates[0])
			if err != nil {
				return err
			}
			content = bytes.TrimSpace(content)
			lines = strings.Split(string(content), "\n")

			// Track all partner's message IDs.
			partnersLatencies := make([]*MetricLatency, len(lines))

			for i := range lines {

				// Split line at whitespace characters.
				metric := strings.Fields(lines[i])

				// Convert first element to timestamp.
				timestamp, err := strconv.ParseInt(metric[0], 10, 64)
				if err != nil {
					return err
				}

				// Convert third element to message ID.
				msgID, err := strconv.ParseInt(metric[2], 10, 64)
				if err != nil {
					return err
				}

				partnersLatencies[i] = &MetricLatency{
					MsgID:            msgID,
					ReceiveTimestamp: timestamp,
				}
			}

			// Sort slice of partner's receive message IDs.
			sort.Slice(partnersLatencies, func(i, j int) bool {
				return partnersLatencies[i].MsgID < partnersLatencies[j].MsgID
			})

			/*
				fmt.Printf("[CCC] Found:\n")
				for i := range partnersLatencies {
					fmt.Printf("\t%d: %d\n", partnersLatencies[i].MsgID, partnersLatencies[i].ReceiveTimestamp)
				}
			*/

			if systemUnderEval == "pung" {

				// Reslice list of messages to capture
				// messages [#2, (#2 + #numMsgs)).
				fmt.Printf("[PUNG] Reslicing partnersLatencies to [2, %d)\n", (2 + int(numMsgsToCalc)))
				partnersLatencies = partnersLatencies[2:(2 + int(numMsgsToCalc))]
			}

			/*
				fmt.Printf("[DDD] Found:\n")
				for i := range partnersLatencies {
					fmt.Printf("\t%d: %d\n", partnersLatencies[i].MsgID, partnersLatencies[i].ReceiveTimestamp)
				}
			*/

			// Determine whether we have all required data points.
			if (int64(len(partnersLatencies)) < numMsgsToCalc) || (partnersLatencies[(numMsgsToCalc-1)].MsgID != numMsgsToCalc) {
				fmt.Printf("Not enough consecutive receive (partner: %s) message time metrics available (found: %d, want: %d).\n", partner,
					partnersLatencies[(numMsgsToCalc-1)].MsgID, numMsgsToCalc)
				os.Exit(1)
			}

			for i := 0; int64(i) < numMsgsToCalc; i++ {

				// Integrate temporarily-stored receive timestamps of partner.
				msgLatencies[i].ReceiveTimestamp = partnersLatencies[i].ReceiveTimestamp

				// Calculate this message's end-to-end latency in seconds.
				latencyNano := msgLatencies[i].ReceiveTimestamp - msgLatencies[i].SendTimestamp
				msgLatencies[i].Latency = float64(latencyNano) / float64(1000000000)

				if msgLatencies[i].Latency <= float64(0.0) {
					fmt.Printf("Non-existent or negative message latency, impossible. Corrupted data or system clocks?\n")
					os.Exit(1)
				}

				// In case one of this client's send timestamps
				// holds a lower value than the previous lowest
				// send timestamp, update the global bound.
				sendTimestampSec := (msgLatencies[i].SendTimestamp / 1000000000) - 1
				if sendTimestampSec < run.TimestampLowest {
					run.TimestampLowest = sendTimestampSec
				}

				// In case one of this client's receive timestamps
				// holds a higher value than the previous highest
				// receive timestamp, update the global bound.
				recvTimestampSec := (msgLatencies[i].ReceiveTimestamp / 1000000000) + 1
				if recvTimestampSec > run.TimestampHighest {
					run.TimestampHighest = recvTimestampSec
				}
			}

			// Reslice to desired size.
			msgLatencies = msgLatencies[:numMsgsToCalc]

			// Append to list of client latencies.
			run.Latencies = append(run.Latencies, msgLatencies)
		}

		return nil
	})

	return err
}

func (set *Setting) LatenciesToFile(path string) error {

	return nil
}

/*
func (clM *ClientMetrics) ClientStoreForBoxplot() error {

	latencyFile, err := os.OpenFile(filepath.Join(clM.MetricsPath, "latency_per_message.boxplot"), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}
	defer latencyFile.Close()
	defer latencyFile.Sync()

	for i := 0; int64(i) < numMsgsToCalc; i++ {

		var values string
		for client := range clM.Latencies {

			if values == "" {
				values = fmt.Sprintf("%.5f", clM.Latencies[client][i].Latency)
			} else {
				values = fmt.Sprintf("%s,%.5f", values, clM.Latencies[client][i].Latency)
			}
		}

		fmt.Fprintln(latencyFile, values)
	}

	return nil
}
*/
