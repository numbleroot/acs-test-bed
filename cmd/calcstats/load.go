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

func (run *Run) AddCPULoad(runNodesPath string, isClientMetric bool) error {

	err := filepath.Walk(runNodesPath, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if filepath.Base(path) == "load_unixnano.evaluation" {

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

				if metric[0] == "n/a" {
					continue
				}

				// Convert first element to timestamp.
				timestamp, err := strconv.ParseInt(metric[0], 10, 64)
				if err != nil {
					return err
				}

				// Exclude metric for further consideration in
				// case it lies outside our zone of interest.
				if (timestamp < run.TimestampLowest) || (timestamp > run.TimestampHighest) {
					continue
				}

				// Convert specific element to idle metrics.
				loadIdle, err := strconv.ParseFloat(strings.TrimPrefix(metric[5], "idle:"), 64)
				if err != nil {
					return err
				}

				// Calculate difference ("busy" load metric).
				loadBusy := 100.0 - loadIdle

				// Append to corresponding slice of values.
				if isClientMetric {
					run.ClientsCPULoad = append(run.ClientsCPULoad, loadBusy)
				} else {
					run.ServersCPULoad = append(run.ServersCPULoad, loadBusy)
				}
			}
		}

		return nil
	})

	return err
}

func (run *Run) AddMemLoad(runNodesPath string, isClientMetric bool) error {

	err := filepath.Walk(runNodesPath, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if filepath.Base(path) == "mem_unixnano.evaluation" {

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

				if metric[0] == "n/a" {
					continue
				}

				// Convert first element to timestamp.
				timestamp, err := strconv.ParseInt(metric[0], 10, 64)
				if err != nil {
					return err
				}

				// Exclude metric for further consideration in
				// case it lies outside our zone of interest.
				if (timestamp < run.TimestampLowest) || (timestamp > run.TimestampHighest) {
					continue
				}

				// Convert following elements to memory metrics.

				memTotal, err := strconv.ParseFloat(strings.TrimPrefix(metric[1], "totalKB:"), 64)
				if err != nil {
					return err
				}

				memAvail, err := strconv.ParseFloat(strings.TrimPrefix(metric[2], "availKB:"), 64)
				if err != nil {
					return err
				}

				// Calculate difference ("used" memory metric).
				memUsed := memTotal - memAvail

				// Append to corresponding slice of values.
				if isClientMetric {
					run.ClientsMemLoad = append(run.ClientsMemLoad, memUsed)
				} else {
					run.ServersMemLoad = append(run.ServersMemLoad, memUsed)
				}
			}
		}

		return nil
	})

	return err
}

func (set *Setting) LoadToFiles(path string) error {

	clientsCPULoadFile, err := os.OpenFile(filepath.Join(path, "cpu_percentage-busy_all-values-in-time-window_clients.data"), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}
	defer clientsCPULoadFile.Close()
	defer clientsCPULoadFile.Sync()

	fmt.Fprintf(clientsCPULoadFile, "%.5f", set.Runs[0].ClientsCPULoad[0])

	for i := range set.Runs {

		for j := range set.Runs[i].ClientsCPULoad {

			if i == 0 && j == 0 {
				continue
			}

			fmt.Fprintf(clientsCPULoadFile, ",%.5f", set.Runs[i].ClientsCPULoad[j])
		}
	}

	fmt.Fprintf(clientsCPULoadFile, "\n")

	/*
		// Calculate average memory usage for clients.
		var clientsMemLoadAvg float64

		allMetricsSum := float64(0.0)
		numMetrics := float64(0.0)

		for i := range set.Runs {

			for j := range set.Runs[i].ClientsMemLoad {
				allMetricsSum = allMetricsSum + set.Runs[i].ClientsMemLoad[j]
			}

			numMetrics = numMetrics + float64(len(set.Runs[i].ClientsMemLoad))
		}

		clientsMemLoadAvg = float64(allMetricsSum / numMetrics)
	*/

	clientsMemLoadFile, err := os.OpenFile(filepath.Join(path, "memory_kilobytes-used_all-values-in-time-window_clients.data"), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}
	defer clientsMemLoadFile.Close()
	defer clientsMemLoadFile.Sync()

	fmt.Fprintf(clientsMemLoadFile, "%.5f", set.Runs[0].ClientsMemLoad[0])

	for i := range set.Runs {

		for j := range set.Runs[i].ClientsMemLoad {

			if i == 0 && j == 0 {
				continue
			}

			fmt.Fprintf(clientsMemLoadFile, ",%.5f", set.Runs[i].ClientsMemLoad[j])
		}
	}

	fmt.Fprintf(clientsMemLoadFile, "\n")

	/*
		// Write values to files for clients.
		fmt.Fprintf(clientsMemLoadFile, "%.5f\n", clientsMemLoadAvg)
	*/

	serversCPULoadFile, err := os.OpenFile(filepath.Join(path, "cpu_percentage-busy_all-values-in-time-window_servers.data"), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}
	defer serversCPULoadFile.Close()
	defer serversCPULoadFile.Sync()

	fmt.Fprintf(serversCPULoadFile, "%.5f", set.Runs[0].ServersCPULoad[0])

	for i := range set.Runs {

		for j := range set.Runs[i].ServersCPULoad {

			if i == 0 && j == 0 {
				continue
			}

			fmt.Fprintf(serversCPULoadFile, ",%.5f", set.Runs[i].ServersCPULoad[j])
		}
	}

	fmt.Fprintf(serversCPULoadFile, "\n")

	/*
		// Calculate average memory usage for servers.
		var serversMemLoadAvg float64

		allMetricsSum = float64(0.0)
		numMetrics = float64(0.0)

		for i := range set.Runs {

			for j := range set.Runs[i].ServersMemLoad {
				allMetricsSum = allMetricsSum + set.Runs[i].ServersMemLoad[j]
			}

			numMetrics = numMetrics + float64(len(set.Runs[i].ServersMemLoad))
		}

		serversMemLoadAvg = float64(allMetricsSum / numMetrics)
	*/

	serversMemLoadFile, err := os.OpenFile(filepath.Join(path, "memory_kilobytes-used_all-values-in-time-window_servers.data"), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}
	defer serversMemLoadFile.Close()
	defer serversMemLoadFile.Sync()

	fmt.Fprintf(serversMemLoadFile, "%.5f", set.Runs[0].ServersMemLoad[0])

	for i := range set.Runs {

		for j := range set.Runs[i].ServersMemLoad {

			if i == 0 && j == 0 {
				continue
			}

			fmt.Fprintf(serversMemLoadFile, ",%.5f", set.Runs[i].ServersMemLoad[j])
		}
	}

	fmt.Fprintf(serversMemLoadFile, "\n")

	/*
		// Write values to files for servers.
		fmt.Fprintf(serversMemLoadFile, "%.5f\n", serversMemLoadAvg)
	*/

	return nil
}
