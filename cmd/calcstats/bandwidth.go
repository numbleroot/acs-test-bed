package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (run *Run) AddSentBytes(runNodesPath string, isClientMetric bool) error {

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

			var reduceBy int64
			var lastValue int64
			var highestValue int64

			// Split file contents into lines.
			lines := strings.Split(string(content), "\n")
			for i := range lines {

				// Split line at whitespace characters.
				metric := strings.Fields(lines[i])

				if metric[0] == "n/a" {
					continue
				}

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

				// Stash value at one second before cutoff.
				if timestamp == (run.TimestampLowest - 1) {
					reduceBy = value
				}

				// Only save the value at highest relevant timestamp.
				if timestamp == run.TimestampHighest {
					highestValue = value
				}

				lastValue = value
			}

			if highestValue == 0 {
				highestValue = lastValue
			}

			// Append as reduced and MiB-normalized highest
			// value to run-global sent bytes slice.
			if isClientMetric {
				run.ClientsSentMiBytesHighest = append(run.ClientsSentMiBytesHighest, (float64((highestValue - reduceBy)) / (1024.0 * 1024.0)))
			} else {
				run.ServersSentMiBytesHighest = append(run.ServersSentMiBytesHighest, (float64((highestValue - reduceBy)) / (1024.0 * 1024.0)))
			}
		}

		return nil
	})

	return err
}

func (run *Run) AddRecvdBytes(runNodesPath string, isClientMetric bool) error {

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

			var reduceBy int64
			var lastValue int64
			var highestValue int64

			// Split file contents into lines.
			lines := strings.Split(string(content), "\n")
			for i := range lines {

				// Split line at whitespace characters.
				metric := strings.Fields(lines[i])

				if metric[0] == "n/a" {
					continue
				}

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

				// Stash value at one second before cutoff.
				if timestamp == (run.TimestampLowest - 1) {
					reduceBy = value
				}

				// Only save the value at highest relevant timestamp.
				if timestamp == run.TimestampHighest {
					highestValue = value
				}

				lastValue = value
			}

			if highestValue == 0 {
				highestValue = lastValue
			}

			// Append as reduced and MiB-normalized highest
			// value to run-global sent bytes slice.
			if isClientMetric {
				run.ClientsRecvdMiBytesHighest = append(run.ClientsRecvdMiBytesHighest, (float64((highestValue - reduceBy)) / (1024.0 * 1024.0)))
			} else {
				run.ServersRecvdMiBytesHighest = append(run.ServersRecvdMiBytesHighest, (float64((highestValue - reduceBy)) / (1024.0 * 1024.0)))
			}
		}

		return nil
	})

	return err
}

func (set *Setting) BandwidthToFiles(path string) error {

	// Calculate bandwidth average for clients.
	var clientsBandwidthAvg float64

	allMetricsSum := float64(0.0)
	numMetrics := float64(len(set.Runs))

	for i := range set.Runs {

		for j := range set.Runs[i].ClientsSentMiBytesHighest {
			allMetricsSum = allMetricsSum + (set.Runs[i].ClientsSentMiBytesHighest[j] + set.Runs[i].ClientsRecvdMiBytesHighest[j])
		}
	}

	clientsBandwidthAvg = float64(allMetricsSum / numMetrics)

	clientsBandwidthAvgFile, err := os.OpenFile(filepath.Join(path, "bandwidth_highest_avg_clients.data"), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}
	defer clientsBandwidthAvgFile.Close()
	defer clientsBandwidthAvgFile.Sync()

	// Write values to files for clients.
	fmt.Fprintf(clientsBandwidthAvgFile, "%.5f\n", clientsBandwidthAvg)

	// Calculate bandwidth average for servers.
	var serversBandwidthAvg float64

	allMetricsSum = 0.0
	numMetrics = float64(len(set.Runs))

	for i := range set.Runs {

		for j := range set.Runs[i].ServersSentMiBytesHighest {
			allMetricsSum = allMetricsSum + (set.Runs[i].ServersSentMiBytesHighest[j] + set.Runs[i].ServersRecvdMiBytesHighest[j])
		}
	}

	serversBandwidthAvg = float64(allMetricsSum / numMetrics)

	serversBandwidthAvgFile, err := os.OpenFile(filepath.Join(path, "bandwidth_highest_avg_servers.data"), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}
	defer serversBandwidthAvgFile.Close()
	defer serversBandwidthAvgFile.Sync()

	// Write values to files for servers.
	fmt.Fprintf(serversBandwidthAvgFile, "%.5f\n", serversBandwidthAvg)

	return nil
}
