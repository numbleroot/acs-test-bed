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

			// Append as reduced and KiB-normalized highest
			// value to run-global sent bytes slice.
			if isClientMetric {
				run.ClientsSentKiBytesHighest = append(run.ClientsSentKiBytesHighest, (float64((highestValue - reduceBy)) / 1024.0))
			} else {
				run.ServersSentKiBytesHighest = append(run.ServersSentKiBytesHighest, (float64((highestValue - reduceBy)) / 1024.0))
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

			// Append as reduced and KiB-normalized highest
			// value to run-global sent bytes slice.
			if isClientMetric {
				run.ClientsRecvdKiBytesHighest = append(run.ClientsRecvdKiBytesHighest, (float64((highestValue - reduceBy)) / 1024.0))
			} else {
				run.ServersRecvdKiBytesHighest = append(run.ServersRecvdKiBytesHighest, (float64((highestValue - reduceBy)) / 1024.0))
			}
		}

		return nil
	})

	return err
}

func (set *Setting) BandwidthToFiles(path string) error {

	// Calculate combined bandwidth average
	// and median values for clients.
	var clientsBandwidthAvg float64
	var clientsBandwidthMed float64

	var numMetrics float64
	var allMetricsSum float64
	allMetrics := make([]float64, 0, (len(set.Runs) * len(set.Runs[0].ClientsSentKiBytesHighest)))

	for i := range set.Runs {

		for j := range set.Runs[i].ClientsSentKiBytesHighest {
			allMetrics = append(allMetrics, (set.Runs[i].ClientsSentKiBytesHighest[j] + set.Runs[i].ClientsRecvdKiBytesHighest[j]))
		}

		numMetrics += float64(len(set.Runs[i].ClientsSentKiBytesHighest))
	}

	// Sort slice holding all metrics for clients
	// by size for median determination.
	sort.Slice(allMetrics, func(i, j int) bool {
		return allMetrics[i] < allMetrics[j]
	})

	for i := range allMetrics {
		allMetricsSum += allMetrics[i]
	}

	clientsBandwidthAvg = float64(allMetricsSum / numMetrics)
	clientsBandwidthMed = allMetrics[(len(allMetrics) / 2)]

	clientsBandwidthAvgFile, err := os.OpenFile(filepath.Join(path, "bandwidth_highest_avg_clients.data"), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}
	defer clientsBandwidthAvgFile.Close()
	defer clientsBandwidthAvgFile.Sync()

	clientsBandwidthMedFile, err := os.OpenFile(filepath.Join(path, "bandwidth_highest_med_clients.data"), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}
	defer clientsBandwidthMedFile.Close()
	defer clientsBandwidthMedFile.Sync()

	// Write values to files for clients.
	fmt.Fprintf(clientsBandwidthAvgFile, "%.5f\n", clientsBandwidthAvg)
	fmt.Fprintf(clientsBandwidthMedFile, "%.5f\n", clientsBandwidthMed)

	// Calculate combined bandwidth average
	// and median values for servers.
	var serversBandwidthAvg float64
	var serversBandwidthMed float64

	numMetrics = 0.0
	allMetricsSum = 0.0
	allMetrics = make([]float64, 0, (len(set.Runs) * len(set.Runs[0].ServersSentKiBytesHighest)))

	for i := range set.Runs {

		for j := range set.Runs[i].ServersSentKiBytesHighest {
			allMetrics = append(allMetrics, (set.Runs[i].ServersSentKiBytesHighest[j] + set.Runs[i].ServersRecvdKiBytesHighest[j]))
		}

		numMetrics += float64(len(set.Runs[i].ServersSentKiBytesHighest))
	}

	// Sort slice holding all metrics for servers
	// by size for median determination.
	sort.Slice(allMetrics, func(i, j int) bool {
		return allMetrics[i] < allMetrics[j]
	})

	for i := range allMetrics {
		allMetricsSum += allMetrics[i]
	}

	serversBandwidthAvg = float64(allMetricsSum / numMetrics)
	serversBandwidthMed = allMetrics[(len(allMetrics) / 2)]

	serversBandwidthAvgFile, err := os.OpenFile(filepath.Join(path, "bandwidth_highest_avg_servers.data"), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}
	defer serversBandwidthAvgFile.Close()
	defer serversBandwidthAvgFile.Sync()

	serversBandwidthMedFile, err := os.OpenFile(filepath.Join(path, "bandwidth_highest_med_servers.data"), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}
	defer serversBandwidthMedFile.Close()
	defer serversBandwidthMedFile.Sync()

	// Write values to files for servers.
	fmt.Fprintf(serversBandwidthAvgFile, "%.5f\n", serversBandwidthAvg)
	fmt.Fprintf(serversBandwidthMedFile, "%.5f\n", serversBandwidthMed)

	return nil
}
