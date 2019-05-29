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

func (run *Run) AddSentBytes(runNodesPath string, isClientMetric bool) error {

	sentBytesRaw := make(map[int64][]int64)

	err := filepath.Walk(runNodesPath, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if filepath.Base(path) == "traffic_outgoing.evaluation" {

			// Ingest supplied file.
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			content = bytes.TrimSpace(content)

			// Split file contents into lines.
			lines := strings.Split(string(content), "\n")
			for i := range lines {

				// Split line at whitespace characters.
				metric := strings.Fields(lines[i])

				// Convert first element to timestamp.
				timestamp, err := strconv.ParseInt(metric[0], 10, 64)
				if err != nil {
					return err
				}

				// Convert second element to metric
				// we are interested in.
				value, err := strconv.ParseInt(metric[1], 10, 64)
				if err != nil {
					return err
				}

				// Normalize to KiB.
				value = value / 1024

				// Append to corresponding slice of values.
				sentBytesRaw[timestamp] = append(sentBytesRaw[timestamp], value)
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	sentBytes := make([]*MetricsInt64, 0, len(sentBytesRaw))

	for ts := range sentBytesRaw {

		// Exclude metric for further consideration in
		// case it lies outside our zone of interest.
		if (ts < run.TimestampLowest) || (ts > run.TimestampHighest) {
			continue
		}

		sentBytes = append(sentBytes, &MetricsInt64{
			Timestamp: ts,
			Values:    sentBytesRaw[ts],
		})
	}

	sort.Slice(sentBytes, func(i, j int) bool {
		return sentBytes[i].Timestamp < sentBytes[j].Timestamp
	})

	// Extract number of values in last relevant bucket.
	lastBucket := sentBytes[(len(sentBytes) - 1)].Values

	fmt.Printf("Sent buckets:    highest ts: %d    found num: %d\n", (sentBytes[(len(sentBytes) - 1)].Timestamp), len(lastBucket))

	if isClientMetric {

		// Copy to final structure.
		run.ClientsSentBytesHighest = make([]int64, len(lastBucket))
		for i := range lastBucket {
			run.ClientsSentBytesHighest[i] = lastBucket[i]
		}

	} else {

		// Copy to final structure.
		run.ServersSentBytesHighest = make([]int64, len(lastBucket))
		for i := range lastBucket {
			run.ServersSentBytesHighest[i] = lastBucket[i]
		}
	}

	return nil
}

func (run *Run) AddRecvdBytes(runNodesPath string, isClientMetric bool) error {

	recvdBytesRaw := make(map[int64][]int64)

	err := filepath.Walk(runNodesPath, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if filepath.Base(path) == "traffic_incoming.evaluation" {

			// Ingest supplied file.
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			content = bytes.TrimSpace(content)

			// Split file contents into lines.
			lines := strings.Split(string(content), "\n")
			for i := range lines {

				// Split line at whitespace characters.
				metric := strings.Fields(lines[i])

				// Convert first element to timestamp.
				timestamp, err := strconv.ParseInt(metric[0], 10, 64)
				if err != nil {
					return err
				}

				// Convert second element to metric
				// we are interested in.
				value, err := strconv.ParseInt(metric[1], 10, 64)
				if err != nil {
					return err
				}

				// Normalize to KiB.
				value = value / 1024

				// Append to corresponding slice of values.
				recvdBytesRaw[timestamp] = append(recvdBytesRaw[timestamp], value)
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	recvdBytes := make([]*MetricsInt64, 0, len(recvdBytesRaw))

	for ts := range recvdBytesRaw {

		// Exclude metric for further consideration in
		// case it lies outside our zone of interest.
		if (ts < run.TimestampLowest) || (ts > run.TimestampHighest) {
			continue
		}

		recvdBytes = append(recvdBytes, &MetricsInt64{
			Timestamp: ts,
			Values:    recvdBytesRaw[ts],
		})
	}

	sort.Slice(recvdBytes, func(i, j int) bool {
		return recvdBytes[i].Timestamp < recvdBytes[j].Timestamp
	})

	// Extract number of values in last relevant bucket.
	lastBucket := recvdBytes[(len(recvdBytes) - 1)].Values

	fmt.Printf("Recvd buckets:    highest ts: %d    found num: %d\n", (recvdBytes[(len(recvdBytes) - 1)].Timestamp), len(lastBucket))

	if isClientMetric {

		// Copy to final structure.
		run.ClientsRecvdBytesHighest = make([]int64, len(lastBucket))
		for i := range lastBucket {
			run.ClientsRecvdBytesHighest[i] = lastBucket[i]
		}

	} else {

		// Copy to final structure.
		run.ServersRecvdBytesHighest = make([]int64, len(lastBucket))
		for i := range lastBucket {
			run.ServersRecvdBytesHighest[i] = lastBucket[i]
		}
	}

	return nil
}

func (set *Setting) BandwidthToFiles(path string) error {

	// Calculate combined bandwidth average
	// and median values for clients.
	var clientsBandwidthAvg float64
	var clientsBandwidthMed int64

	var numMetrics int64
	var allMetricsSum int64
	allMetrics := make([]int64, 0, (len(set.Runs) * len(set.Runs[0].ClientsSentBytesHighest)))

	for i := range set.Runs {

		for j := range set.Runs[i].ClientsSentBytesHighest {
			allMetrics = append(allMetrics, (set.Runs[i].ClientsSentBytesHighest[j] + set.Runs[i].ClientsRecvdBytesHighest[j]))
		}

		numMetrics += int64(len(set.Runs[i].ClientsSentBytesHighest))
	}

	// Sort slice holding all metrics for clients
	// by size for median determination.
	sort.Slice(allMetrics, func(i, j int) bool {
		return allMetrics[i] < allMetrics[j]
	})

	for i := range allMetrics {
		fmt.Printf("allClientsMetrics[%d]: %d\n", i, allMetrics[i])
		allMetricsSum += allMetrics[i]
	}

	clientsBandwidthAvg = float64(allMetricsSum / numMetrics)
	clientsBandwidthMed = allMetrics[(len(allMetrics) / 2)]

	clientsBandwidthFile, err := os.OpenFile(filepath.Join(path, "bandwidth_highest_avg_med_clients.data"), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}
	defer clientsBandwidthFile.Close()
	defer clientsBandwidthFile.Sync()

	// Write values to file for clients.
	fmt.Fprintf(clientsBandwidthFile, "%.5f,%d\n", clientsBandwidthAvg, clientsBandwidthMed)

	// Calculate combined bandwidth average
	// and median values for servers.
	var serversBandwidthAvg float64
	var serversBandwidthMed int64

	numMetrics = 0
	allMetricsSum = 0
	allMetrics = make([]int64, 0, (len(set.Runs) * len(set.Runs[0].ServersSentBytesHighest)))

	for i := range set.Runs {

		for j := range set.Runs[i].ServersSentBytesHighest {
			allMetrics = append(allMetrics, (set.Runs[i].ServersSentBytesHighest[j] + set.Runs[i].ServersRecvdBytesHighest[j]))
		}

		numMetrics += int64(len(set.Runs[i].ServersSentBytesHighest))
	}

	// Sort slice holding all metrics for servers
	// by size for median determination.
	sort.Slice(allMetrics, func(i, j int) bool {
		return allMetrics[i] < allMetrics[j]
	})

	for i := range allMetrics {
		fmt.Printf("allServersMetrics[%d]: %d\n", i, allMetrics[i])
		allMetricsSum += allMetrics[i]
	}

	serversBandwidthAvg = float64(allMetricsSum / numMetrics)
	serversBandwidthMed = allMetrics[(len(allMetrics) / 2)]

	serversBandwidthFile, err := os.OpenFile(filepath.Join(path, "bandwidth_highest_avg_med_servers.data"), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}
	defer serversBandwidthFile.Close()
	defer serversBandwidthFile.Sync()

	// Write values to file for clients.
	fmt.Fprintf(serversBandwidthFile, "%.5f,%d\n", serversBandwidthAvg, serversBandwidthMed)

	return nil
}
