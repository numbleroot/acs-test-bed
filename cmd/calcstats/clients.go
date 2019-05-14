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

type MetricLatency struct {
	MsgID            int64
	SendTimestamp    int64
	ReceiveTimestamp int64
	Latency          int64
}

type ClientMetrics struct {
	*SystemMetrics
	MetricsPath       string
	NumMsgsToCalc     int64
	LatencyLowerBound int64
	LatencyUpperBound int64
	Latencies         [][]*MetricLatency
}

func (clM *ClientMetrics) AddLatency(path string) error {

	// Ingest supplied send times file.
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	content = bytes.TrimSpace(content)

	// Split file contents into lines.
	lines := strings.Split(string(content), "\n")

	var partner string

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

		// Convert second element to partner.
		if partner == "" {
			partner = fmt.Sprintf("mixnet-%s", strings.Split(metric[1], "=>")[1])
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

	// Determine whether we have all required data points.
	if (int64(len(msgLatencies)) < clM.NumMsgsToCalc) || (msgLatencies[(clM.NumMsgsToCalc-1)].MsgID != clM.NumMsgsToCalc) {
		fmt.Printf("Not enough consecutive send message time metrics available (want: %d).\n", clM.NumMsgsToCalc)
		os.Exit(1)
	}

	// Find partner's receive time file.
	candidates, err := filepath.Glob(fmt.Sprintf("%s/%s_*/recv_unixnano.evaluation", clM.MetricsPath, partner))
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

	// Determine whether we have all required data points.
	if (int64(len(partnersLatencies)) < clM.NumMsgsToCalc) || (partnersLatencies[(clM.NumMsgsToCalc-1)].MsgID != clM.NumMsgsToCalc) {
		fmt.Printf("Not enough consecutive receive message time metrics available (want: %d).\n", clM.NumMsgsToCalc)
		os.Exit(1)
	}

	for i := 0; int64(i) < clM.NumMsgsToCalc; i++ {

		// Combine and calculate latency metrics.
		msgLatencies[i].ReceiveTimestamp = partnersLatencies[i].ReceiveTimestamp
		msgLatencies[i].Latency = msgLatencies[i].ReceiveTimestamp - msgLatencies[i].SendTimestamp

		if msgLatencies[i].Latency <= int64(0) {
			fmt.Printf("Non-existent or negative message latency, impossible. Corrupted data or system clocks?\n")
			os.Exit(1)
		}
	}

	// Reslice to desired size.
	msgLatencies = msgLatencies[:clM.NumMsgsToCalc]

	// Append to list of client latencies.
	clM.Latencies = append(clM.Latencies, msgLatencies)

	return nil
}

func (clM *ClientMetrics) CalcAndWriteLatency(path string) error {
	return nil
}
