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

func (run *Run) AddLatency(runClientsPath string, numMsgsToCalc int64) error {

	err := filepath.Walk(runClientsPath, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return fmt.Errorf("path: '%s', err: %v", path, err)
		}

		if strings.HasSuffix(filepath.Base(path), "send_unixnano.evaluation") {

			partner := ""

			// Ingest supplied send times file.
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("path: '%s', err: %v", path, err)
			}
			content = bytes.TrimSpace(content)

			// Split file contents into lines.
			lines := strings.Split(string(content), "\n")

			// Prepare latencies state object.
			clientMsgLatencies := make([]*MetricLatency, len(lines))

			for i := range lines {

				// Split line at whitespace characters.
				metric := strings.Fields(lines[i])

				// Convert first element to timestamp.
				timestamp, err := strconv.ParseInt(metric[0], 10, 64)
				if err != nil {
					return fmt.Errorf("path: '%s', err: %v", path, err)
				}

				// Extract name of partner node.
				if partner == "" {
					partner = strings.Split(metric[1], "=>")[1]
				}

				// Convert third element to message ID.
				msgID, err := strconv.ParseInt(metric[2], 10, 64)
				if err != nil {
					return fmt.Errorf("path: '%s', err: %v", path, err)
				}

				// Append to temporary state object.
				clientMsgLatencies[i] = &MetricLatency{
					MsgID:         msgID,
					SendTimestamp: timestamp,
				}
			}

			// Sort state object by message IDs.
			sort.Slice(clientMsgLatencies, func(i, j int) bool {
				return clientMsgLatencies[i].MsgID < clientMsgLatencies[j].MsgID
			})

			clientStart := -1
			clientEnd := -1

			for i := range clientMsgLatencies {

				if clientMsgLatencies[i].MsgID < 3 {
					continue
				}

				if clientStart == -1 {
					clientStart = i
				}

				if clientMsgLatencies[i].MsgID >= (numMsgsToCalc + 3) {
					clientEnd = i
					break
				}
			}

			if clientStart == -1 || clientEnd == -1 {
				return fmt.Errorf("did not find adequate latency bounds on sender (start=%d, end=%d, path: '%s')",
					clientStart, clientEnd, path)
			}

			// Reslice to found bounds.
			clientMsgLatencies = clientMsgLatencies[clientStart:clientEnd]

			if int64(len(clientMsgLatencies)) != numMsgsToCalc {
				return fmt.Errorf("too few latency metrics in clientMsgLatencies (want %d, saw %d, path: '%s')",
					numMsgsToCalc, len(clientMsgLatencies), path)
			}

			// Find partner's receive time file.
			partnerFilePath := fmt.Sprintf("%s/*/%s_recv_unixnano.evaluation", runClientsPath, partner)
			candidates, err := filepath.Glob(partnerFilePath)
			if err != nil {
				return fmt.Errorf("path: '%s', err: %v", path, err)
			}

			if len(candidates) != 1 {
				return fmt.Errorf("client %s did not have unique conversation partner", path)
			}

			fmt.Printf("partner: '%s'\n", candidates[0])

			// Ingest partner's receive times file.
			content, err = ioutil.ReadFile(candidates[0])
			if err != nil {
				return fmt.Errorf("path: '%s', err: %v", path, err)
			}
			content = bytes.TrimSpace(content)
			lines = strings.Split(string(content), "\n")

			// Track all partner's message IDs.
			partnersMsgLatencies := make([]*MetricLatency, len(lines))

			for i := range lines {

				// Split line at whitespace characters.
				metric := strings.Fields(lines[i])

				// Convert first element to timestamp.
				timestamp, err := strconv.ParseInt(metric[0], 10, 64)
				if err != nil {
					return fmt.Errorf("path: '%s', err: %v", path, err)
				}

				// Convert third element to message ID.
				msgID, err := strconv.ParseInt(metric[2], 10, 64)
				if err != nil {
					return fmt.Errorf("path: '%s', err: %v", path, err)
				}

				partnersMsgLatencies[i] = &MetricLatency{
					MsgID:            msgID,
					ReceiveTimestamp: timestamp,
				}
			}

			// Sort slice of partner's receive message IDs.
			sort.Slice(partnersMsgLatencies, func(i, j int) bool {
				return partnersMsgLatencies[i].MsgID < partnersMsgLatencies[j].MsgID
			})

			partnerStart := -1
			partnerEnd := -1

			for i := range partnersMsgLatencies {

				if partnersMsgLatencies[i].MsgID < 3 {
					continue
				}

				if partnerStart == -1 {
					partnerStart = i
				}

				if partnersMsgLatencies[i].MsgID >= (numMsgsToCalc + 3) {
					partnerEnd = i
					break
				}
			}

			if partnerStart == -1 || partnerEnd == -1 {
				return fmt.Errorf("did not find adequate latency bounds on recipient (start=%d, end=%d, path: '%s')",
					partnerStart, partnerEnd, path)
			}

			// Reslice to found bounds.
			partnersMsgLatencies = partnersMsgLatencies[partnerStart:partnerEnd]

			if int64(len(partnersMsgLatencies)) != numMsgsToCalc {
				return fmt.Errorf("too few latency metrics in partnersMsgLatencies (want %d, saw %d, path: '%s')",
					numMsgsToCalc, len(partnersMsgLatencies), path)
			}

			msgLatencies := make([]*MetricLatency, 0, len(clientMsgLatencies))

			for i := range clientMsgLatencies {

				// Calculate this message's end-to-end latency in seconds.
				latencyNano := partnersMsgLatencies[i].ReceiveTimestamp - clientMsgLatencies[i].SendTimestamp
				latencySec := float64(latencyNano) / float64(1000000000)

				if latencyNano <= int64(0) {

					// Negative latencies should not be possible.
					return fmt.Errorf("incorrect latency: %dns (%.5fs) (ID: %d, %d => %d, %s)", latencyNano, latencySec,
						clientMsgLatencies[i].MsgID, clientMsgLatencies[i].SendTimestamp,
						partnersMsgLatencies[i].ReceiveTimestamp, path)
				}

				// Append to struct capturing all useful latency metrics.
				msgLatencies = append(msgLatencies, &MetricLatency{
					MsgID:            clientMsgLatencies[i].MsgID,
					SendTimestamp:    clientMsgLatencies[i].SendTimestamp,
					ReceiveTimestamp: partnersMsgLatencies[i].ReceiveTimestamp,
					Latency:          latencySec,
				})
			}

			// Reslice to desired size.
			msgLatencies = msgLatencies[:numMsgsToCalc]

			for i := range msgLatencies {

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

			// Append to list of client latencies.
			run.Latencies = append(run.Latencies, msgLatencies)
		}

		return nil
	})

	return err
}

func (set *Setting) LatenciesToFile(path string) error {

	clientsLatenciesFile, err := os.OpenFile(filepath.Join(path,
		"transmission-latencies_seconds_all-values-in-time-window.data"),
		(os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}
	defer clientsLatenciesFile.Close()
	defer clientsLatenciesFile.Sync()

	fmt.Fprintf(clientsLatenciesFile, "%.5f", set.Runs[0].Latencies[0][0].Latency)

	for i := range set.Runs {

		for j := range set.Runs[i].Latencies {

			for k := range set.Runs[i].Latencies[j] {

				if i == 0 && j == 0 && k == 0 {
					continue
				}

				fmt.Fprintf(clientsLatenciesFile, ",%.5f", set.Runs[i].Latencies[j][k].Latency)
			}
		}
	}

	fmt.Fprintf(clientsLatenciesFile, "\n")

	return nil
}

func (set *Setting) TotalExpTimesToFile(path string) error {

	clientsTotalExpTimesFile, err := os.OpenFile(filepath.Join(path,
		"total-experiment-times_seconds.data"), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}
	defer clientsTotalExpTimesFile.Close()
	defer clientsTotalExpTimesFile.Sync()

	allTotals := make([]string, len(set.Runs))

	for i := range set.Runs {
		allTotals[i] = fmt.Sprintf("%d", ((set.Runs[i].TimestampHighest - set.Runs[i].TimestampLowest) + 2))
	}

	fmt.Fprintf(clientsTotalExpTimesFile, "%s\n", strings.Join(allTotals, ","))

	return nil
}
