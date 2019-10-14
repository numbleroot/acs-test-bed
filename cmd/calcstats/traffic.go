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

// AddHighestSentBytes appends the node-individual
// highest outgoing traffic volume measurement.
// Respective values lie within the timestamp boundaries
// of interest of this run, are reduced by the warm-up
// offset, and MiB-normalized.
func (run *Run) AddHighestSentBytes(runNodesPath string, isClientMetric bool) error {

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

			// Ensure there are actually traffic values tracked
			// on outgoing connections from this node.
			if lastValue == 0 {
				return fmt.Errorf("outgoing traffic volume non-existent (reduceBy=%d, lastValue=%d, highestValue=%d, path: '%s')",
					reduceBy, lastValue, highestValue, path)
			}

			// If we did not see a value for the highest timestamp,
			// set it to the last seen value instead.
			if highestValue == 0 {
				highestValue = lastValue
			}

			// Calculate MiB-normalized highest outgoing traffic
			// volume value that is offset by what the warm-up
			// phase of the ACS produced in outgoing traffic.
			highestSentMiB := float64((highestValue - reduceBy)) / (1024.0 * 1024.0)

			// Append value to run-global sent bytes slice.
			if isClientMetric {
				run.ClientsSentMiBytesHighest = append(run.ClientsSentMiBytesHighest, highestSentMiB)
			} else {
				run.ServersSentMiBytesHighest = append(run.ServersSentMiBytesHighest, highestSentMiB)
			}
		}

		return nil
	})

	return err
}

// AddHighestRecvdBytes appends the node-individual
// highest incoming traffic volume measurement.
// Respective values lie within the timestamp boundaries
// of interest of this run, are reduced by the warm-up
// offset, and MiB-normalized.
func (run *Run) AddHighestRecvdBytes(runNodesPath string, isClientMetric bool) error {

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

			// Ensure there are actually traffic values tracked
			// on incoming connections to this node.
			if lastValue == 0 {
				return fmt.Errorf("incoming traffic volume non-existent (reduceBy=%d, lastValue=%d, highestValue=%d, path: '%s')",
					reduceBy, lastValue, highestValue, path)
			}

			// If we did not see a value for the highest timestamp,
			// set it to the last seen value instead.
			if highestValue == 0 {
				highestValue = lastValue
			}

			// Calculate MiB-normalized highest incoming traffic
			// volume value that is offset by what the warm-up
			// phase of the ACS produced in incoming traffic.
			highestRecvdMiB := float64((highestValue - reduceBy)) / (1024.0 * 1024.0)

			// Append value to run-global sent bytes slice.
			if isClientMetric {
				run.ClientsRecvdMiBytesHighest = append(run.ClientsRecvdMiBytesHighest, highestRecvdMiB)
			} else {
				run.ServersRecvdMiBytesHighest = append(run.ServersRecvdMiBytesHighest, highestRecvdMiB)
			}
		}

		return nil
	})

	return err
}

// TrafficToFiles calculates the per-client and per-server
// highest traffic volume average across all runs of this
// setting and writes each value out to its respective file.
func (set *Setting) TrafficToFiles(path string) error {

	clientsBandwidthAvg := float64(-1.0)
	allMetricsSum := float64(0.0)
	numMetrics := float64(0.0)

	// Add up the highest outgoing and incoming
	// traffic values for all clients and runs.
	for i := range set.Runs {

		for j := range set.Runs[i].ClientsSentMiBytesHighest {
			allMetricsSum = allMetricsSum + set.Runs[i].ClientsSentMiBytesHighest[j]
		}

		for k := range set.Runs[i].ClientsRecvdMiBytesHighest {
			allMetricsSum = allMetricsSum + set.Runs[i].ClientsRecvdMiBytesHighest[k]
		}

		numMetrics = numMetrics + set.Runs[i].NumClients
	}

	// Divide the total traffic volume for all clients
	// across all runs by the total number of clients
	// across all runs.
	clientsBandwidthAvg = allMetricsSum / numMetrics

	clientsBandwidthAvgFile, err := os.OpenFile(
		filepath.Join(path, "traffic-volume_mebibytes_highest-at-end-of-time-window_clients.data"),
		(os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}
	defer clientsBandwidthAvgFile.Close()
	defer clientsBandwidthAvgFile.Sync()

	// Write value to file for clients.
	fmt.Fprintf(clientsBandwidthAvgFile, "%.5f\n", clientsBandwidthAvg)

	// Do the same accordingly on server side.
	serversBandwidthAvg := float64(-1.0)
	allMetricsSum = float64(0.0)
	numMetrics = float64(0.0)

	// Add up the highest outgoing and incoming
	// traffic values for all servers and runs.
	for i := range set.Runs {

		for j := range set.Runs[i].ServersSentMiBytesHighest {
			allMetricsSum = allMetricsSum + set.Runs[i].ServersSentMiBytesHighest[j]
		}

		for k := range set.Runs[i].ServersRecvdMiBytesHighest {
			allMetricsSum = allMetricsSum + set.Runs[i].ServersRecvdMiBytesHighest[k]
		}

		numMetrics = numMetrics + set.Runs[i].NumServers
	}

	// Divide the total traffic volume for all servers
	// across all runs by the total number of servers
	// across all runs.
	serversBandwidthAvg = allMetricsSum / numMetrics

	serversBandwidthAvgFile, err := os.OpenFile(
		filepath.Join(path, "traffic-volume_mebibytes_highest-at-end-of-time-window_servers.data"),
		(os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}
	defer serversBandwidthAvgFile.Close()
	defer serversBandwidthAvgFile.Sync()

	// Write value to file for servers.
	fmt.Fprintf(serversBandwidthAvgFile, "%.5f\n", serversBandwidthAvg)

	return nil
}
