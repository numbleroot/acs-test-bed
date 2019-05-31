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
			clientMsgLatencies := make([]*MetricLatency, len(lines))

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
				}

				// Convert third element to message ID.
				msgID, err := strconv.ParseInt(metric[2], 10, 64)
				if err != nil {
					return err
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

			// Reslice list of messages to capture messages [#3, ...).
			clientMsgLatencies = clientMsgLatencies[2:]

			// Ensure we have enough values available.
			for i := 0; int64(i) < numMsgsToCalc; i++ {

				if clientMsgLatencies[i].MsgID != int64((i + 3)) {
					return fmt.Errorf("%s does not have consecutive metrics (idx %d => msgID %d)", path, (i + 3), clientMsgLatencies[i].MsgID)
				}
			}

			// Find partner's receive time file.
			candidates, err := filepath.Glob(fmt.Sprintf("%s/%s_*/recv_unixnano.evaluation", runClientsPath, partner))
			if err != nil {
				return err
			}

			if len(candidates) != 1 {
				return fmt.Errorf("client %s did not have unique conversation partner", path)
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

			// Reslice list of messages to capture messages [#3, ...).
			partnersLatencies = partnersLatencies[2:]

			// Ensure we see enough shared metrics between
			// sending node and receiving.
			idxLastCommonMsg := -1
			for i := range clientMsgLatencies {

				if i < len(partnersLatencies) {

					if clientMsgLatencies[i].MsgID == partnersLatencies[i].MsgID {
						idxLastCommonMsg = i
					}
				}
			}

			if int64(idxLastCommonMsg) < (numMsgsToCalc - 1) {
				return fmt.Errorf("%s and %s only share %d metrics for messages (%d wanted)", path, partner, idxLastCommonMsg, numMsgsToCalc)
			}

			msgLatencies := make([]*MetricLatency, 0, len(clientMsgLatencies))

			for i := range clientMsgLatencies {

				if i < len(partnersLatencies) {

					// Calculate this message's end-to-end latency in seconds.
					latencyNano := partnersLatencies[i].ReceiveTimestamp - clientMsgLatencies[i].SendTimestamp

					/*
						if latencyNano <= int64(0) {

							// TODO: Potentially do something in case negative and
							//       "overly" positive latencies are not balancing
							//       each other out. Right now, I assume they do
							//       due to symmetry of communication partners and
							//       my larger data set.

							fmt.Printf("Incorrect latency: %d (ID: %d, %d => %d, %s)\n", latencyNano, clientMsgLatencies[i].MsgID,
								clientMsgLatencies[i].SendTimestamp, partnersLatencies[i].ReceiveTimestamp, path)
							continue
						}
					*/

					// Append to struct capturing all useful latency metrics.
					msgLatencies = append(msgLatencies, &MetricLatency{
						MsgID:            clientMsgLatencies[i].MsgID,
						SendTimestamp:    clientMsgLatencies[i].SendTimestamp,
						ReceiveTimestamp: partnersLatencies[i].ReceiveTimestamp,
						Latency:          float64(latencyNano) / float64(1000000000),
					})
				}
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

	clientsLatenciesFile, err := os.OpenFile(filepath.Join(path, "msg-latencies_lowest-to-highest_clients.data"), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
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
