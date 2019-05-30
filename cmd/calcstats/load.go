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

func (run *Run) AddMemLoad(runNodesPath string, isClientMetric bool) error {

	memLoadRaw := make(map[int64][]float64)

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

				// Calculate ratio of used to total memory.
				memUsedRatio := (float64(memUsed / memTotal)) * 100.0

				// Append to corresponding slice of values.
				memLoadRaw[timestamp] = append(memLoadRaw[timestamp], memUsedRatio)
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	if isClientMetric {

		clientsMemLoad := make([]*MetricsFloat64, 0, len(memLoadRaw))

		for ts := range memLoadRaw {

			// Exclude metric for further consideration in
			// case it lies outside our zone of interest.
			if (ts < run.TimestampLowest) || (ts > run.TimestampHighest) {
				continue
			}

			clientsMemLoad = append(clientsMemLoad, &MetricsFloat64{
				Timestamp: ts,
				Values:    memLoadRaw[ts],
			})
		}

		sort.Slice(clientsMemLoad, func(i, j int) bool {
			return clientsMemLoad[i].Timestamp < clientsMemLoad[j].Timestamp
		})

		run.ClientsMemLoad = make([]float64, 0, (len(clientsMemLoad) * len(clientsMemLoad[0].Values)))

		for i := range clientsMemLoad {

			for j := range clientsMemLoad[i].Values {
				run.ClientsMemLoad = append(run.ClientsMemLoad, clientsMemLoad[i].Values[j])
			}
		}

	} else {

		serversMemLoad := make([]*MetricsFloat64, 0, len(memLoadRaw))

		for ts := range memLoadRaw {

			// Exclude metric for further consideration in
			// case it lies outside our zone of interest.
			if (ts < run.TimestampLowest) || (ts > run.TimestampHighest) {
				continue
			}

			serversMemLoad = append(serversMemLoad, &MetricsFloat64{
				Timestamp: ts,
				Values:    memLoadRaw[ts],
			})
		}

		sort.Slice(serversMemLoad, func(i, j int) bool {
			return serversMemLoad[i].Timestamp < serversMemLoad[j].Timestamp
		})

		run.ServersMemLoad = make([]float64, 0, (len(serversMemLoad) * len(serversMemLoad[0].Values)))

		for i := range serversMemLoad {

			for j := range serversMemLoad[i].Values {
				run.ServersMemLoad = append(run.ServersMemLoad, serversMemLoad[i].Values[j])
			}
		}
	}

	return nil
}

func (run *Run) AddCPULoad(runNodesPath string, isClientMetric bool) error {

	cpuLoadRaw := make(map[int64][]float64)

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

				// Convert specific element to idle metrics.
				loadIdle, err := strconv.ParseFloat(strings.TrimPrefix(metric[5], "idle:"), 64)
				if err != nil {
					return err
				}

				// Calculate difference ("busy" load metric).
				loadBusy := 100.0 - loadIdle

				// Append to corresponding slice of values.
				cpuLoadRaw[timestamp] = append(cpuLoadRaw[timestamp], loadBusy)
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	if isClientMetric {

		clientsCPULoad := make([]*MetricsFloat64, 0, len(cpuLoadRaw))

		for ts := range cpuLoadRaw {

			// Exclude metric for further consideration in
			// case it lies outside our zone of interest.
			if (ts < run.TimestampLowest) || (ts > run.TimestampHighest) {
				continue
			}

			clientsCPULoad = append(clientsCPULoad, &MetricsFloat64{
				Timestamp: ts,
				Values:    cpuLoadRaw[ts],
			})
		}

		sort.Slice(clientsCPULoad, func(i, j int) bool {
			return clientsCPULoad[i].Timestamp < clientsCPULoad[j].Timestamp
		})

		run.ClientsCPULoad = make([]float64, 0, (len(clientsCPULoad) * len(clientsCPULoad[0].Values)))

		for i := range clientsCPULoad {

			for j := range clientsCPULoad[i].Values {
				run.ClientsCPULoad = append(run.ClientsCPULoad, clientsCPULoad[i].Values[j])
			}
		}

	} else {

		serversCPULoad := make([]*MetricsFloat64, 0, len(cpuLoadRaw))

		for ts := range cpuLoadRaw {

			// Exclude metric for further consideration in
			// case it lies outside our zone of interest.
			if (ts < run.TimestampLowest) || (ts > run.TimestampHighest) {
				continue
			}

			serversCPULoad = append(serversCPULoad, &MetricsFloat64{
				Timestamp: ts,
				Values:    cpuLoadRaw[ts],
			})
		}

		sort.Slice(serversCPULoad, func(i, j int) bool {
			return serversCPULoad[i].Timestamp < serversCPULoad[j].Timestamp
		})

		run.ServersCPULoad = make([]float64, 0, (len(serversCPULoad) * len(serversCPULoad[0].Values)))

		for i := range serversCPULoad {

			for j := range serversCPULoad[i].Values {
				run.ServersCPULoad = append(run.ServersCPULoad, serversCPULoad[i].Values[j])
			}
		}
	}

	return nil
}

func (set *Setting) LoadToFiles(path string) error {

	metrics := ""
	for i := range set.Runs {

		for j := range set.Runs[i].ClientsCPULoad {

			if metrics == "" {
				metrics = fmt.Sprintf("%f", set.Runs[i].ClientsCPULoad[j])
			} else {
				metrics = fmt.Sprintf("%s,%f", metrics, set.Runs[i].ClientsCPULoad[j])
			}
		}
	}

	clientsCPULoadFile, err := os.OpenFile(filepath.Join(path, "load_cpu_lowest-to-highest_clients.data"), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}
	defer clientsCPULoadFile.Close()
	defer clientsCPULoadFile.Sync()

	fmt.Fprintf(clientsCPULoadFile, "%s\n", metrics)

	metrics = ""
	for i := range set.Runs {

		for j := range set.Runs[i].ClientsMemLoad {

			if metrics == "" {
				metrics = fmt.Sprintf("%f", set.Runs[i].ClientsMemLoad[j])
			} else {
				metrics = fmt.Sprintf("%s,%f", metrics, set.Runs[i].ClientsMemLoad[j])
			}
		}
	}

	clientsMemLoadFile, err := os.OpenFile(filepath.Join(path, "load_mem_lowest-to-highest_clients.data"), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}
	defer clientsMemLoadFile.Close()
	defer clientsMemLoadFile.Sync()

	fmt.Fprintf(clientsMemLoadFile, "%s\n", metrics)

	metrics = ""
	for i := range set.Runs {

		for j := range set.Runs[i].ServersCPULoad {

			if metrics == "" {
				metrics = fmt.Sprintf("%f", set.Runs[i].ServersCPULoad[j])
			} else {
				metrics = fmt.Sprintf("%s,%f", metrics, set.Runs[i].ServersCPULoad[j])
			}
		}
	}

	serversCPULoadFile, err := os.OpenFile(filepath.Join(path, "load_cpu_lowest-to-highest_servers.data"), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}
	defer serversCPULoadFile.Close()
	defer serversCPULoadFile.Sync()

	fmt.Fprintf(serversCPULoadFile, "%s\n", metrics)

	metrics = ""
	for i := range set.Runs {

		for j := range set.Runs[i].ServersMemLoad {

			if metrics == "" {
				metrics = fmt.Sprintf("%f", set.Runs[i].ServersMemLoad[j])
			} else {
				metrics = fmt.Sprintf("%s,%f", metrics, set.Runs[i].ServersMemLoad[j])
			}
		}
	}

	serversMemLoadFile, err := os.OpenFile(filepath.Join(path, "load_mem_lowest-to-highest_servers.data"), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}
	defer serversMemLoadFile.Close()
	defer serversMemLoadFile.Sync()

	fmt.Fprintf(serversMemLoadFile, "%s\n", metrics)

	return nil
}
