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

// AddCPULoad traverses the supplied run's directory
// and extracts the "busy CPU" percentage for each
// second within the time interval of interest.
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
				// case it lies outside our interval of interest.
				if (timestamp < run.TimestampLowest) || (timestamp > run.TimestampHighest) {
					continue
				}

				// Convert specific element to idle metrics.
				loadIdle, err := strconv.ParseFloat(strings.TrimPrefix(metric[5], "idle:"), 64)
				if err != nil {
					return err
				}

				// Calculate difference ("busy" load metric).
				loadBusy := float64(100.0) - loadIdle

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

// AddMemLoad traverses the supplied run's directory
// and extracts the used RAM share in Megabytes for
// each second within the time interval of interest.
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
				// case it lies outside our interval of interest.
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

				// Append to corresponding slice of values,
				// in MB-normalized for clients and GB-normalized
				// on side of servers.
				if isClientMetric {
					run.ClientsMemLoad = append(run.ClientsMemLoad, (memUsed / float64(1000)))
				} else {
					run.ServersMemLoad = append(run.ServersMemLoad, (memUsed / float64(1000000)))
				}
			}
		}

		return nil
	})

	return err
}

// LoadToFiles writes all prepared load metrics
// for clients and servers into respective files.
func (set *Setting) LoadToFiles(path string) error {

	clientsCPULoadFile, err := os.OpenFile(
		filepath.Join(path, "cpu_percentage-busy_all-values-in-time-window_clients.data"),
		(os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}
	defer clientsCPULoadFile.Close()
	defer clientsCPULoadFile.Sync()

	// Write first CPU load value of first
	// run to file for clients.
	fmt.Fprintf(clientsCPULoadFile, "%.5f", set.Runs[0].ClientsCPULoad[0])

	for i := range set.Runs {

		for j := range set.Runs[i].ClientsCPULoad {

			// Do not write the first value of
			// the first run again.
			if i == 0 && j == 0 {
				continue
			}

			// Write all other values to file.
			fmt.Fprintf(clientsCPULoadFile, ",%.5f", set.Runs[i].ClientsCPULoad[j])
		}
	}

	fmt.Fprintf(clientsCPULoadFile, "\n")

	clientsMemLoadFile, err := os.OpenFile(
		filepath.Join(path, "memory_megabytes-used_all-values-in-time-window_clients.data"),
		(os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}
	defer clientsMemLoadFile.Close()
	defer clientsMemLoadFile.Sync()

	// Write first RAM load value of first
	// run to file for clients.
	fmt.Fprintf(clientsMemLoadFile, "%.5f", set.Runs[0].ClientsMemLoad[0])

	for i := range set.Runs {

		for j := range set.Runs[i].ClientsMemLoad {

			if i == 0 && j == 0 {
				continue
			}

			// Write all remaining values.
			fmt.Fprintf(clientsMemLoadFile, ",%.5f", set.Runs[i].ClientsMemLoad[j])
		}
	}

	fmt.Fprintf(clientsMemLoadFile, "\n")

	serversCPULoadFile, err := os.OpenFile(
		filepath.Join(path, "cpu_percentage-busy_all-values-in-time-window_servers.data"),
		(os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}
	defer serversCPULoadFile.Close()
	defer serversCPULoadFile.Sync()

	// Proceed similarly server-side.
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

	serversMemLoadFile, err := os.OpenFile(
		filepath.Join(path, "memory_gigabytes-used_all-values-in-time-window_servers.data"),
		(os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
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

	return nil
}
