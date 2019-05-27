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

func (run *Run) AddMem(runNodesPath string, isClientMetric bool) error {

	memRaw := make(map[int64][]float64)

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
				memRaw[timestamp] = append(memRaw[timestamp], memUsedRatio)
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	if isClientMetric {

		run.ClientsMemory = make([]*MetricsFloat64, 0, len(memRaw))

		for ts := range memRaw {

			// Exclude metric for further consideration in
			// case it lies outside our zone of interest.
			if (ts < run.TimestampLowest) || (ts > run.TimestampHighest) {
				continue
			}

			run.ClientsMemory = append(run.ClientsMemory, &MetricsFloat64{
				Timestamp: ts,
				Values:    memRaw[ts],
			})
		}

		sort.Slice(run.ClientsMemory, func(i, j int) bool {
			return run.ClientsMemory[i].Timestamp < run.ClientsMemory[j].Timestamp
		})

	} else {

		run.ServersMemory = make([]*MetricsFloat64, 0, len(memRaw))

		for ts := range memRaw {

			// Exclude metric for further consideration in
			// case it lies outside our zone of interest.
			if (ts < run.TimestampLowest) || (ts > run.TimestampHighest) {
				continue
			}

			run.ServersMemory = append(run.ServersMemory, &MetricsFloat64{
				Timestamp: ts,
				Values:    memRaw[ts],
			})
		}

		sort.Slice(run.ServersMemory, func(i, j int) bool {
			return run.ServersMemory[i].Timestamp < run.ServersMemory[j].Timestamp
		})
	}

	return nil
}

func (run *Run) AddLoad(runNodesPath string, isClientMetric bool) error {

	loadRaw := make(map[int64][]float64)

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
				loadRaw[timestamp] = append(loadRaw[timestamp], loadBusy)
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	if isClientMetric {

		run.ClientsLoad = make([]*MetricsFloat64, 0, len(loadRaw))

		for ts := range loadRaw {

			// Exclude metric for further consideration in
			// case it lies outside our zone of interest.
			if (ts < run.TimestampLowest) || (ts > run.TimestampHighest) {
				continue
			}

			run.ClientsLoad = append(run.ClientsLoad, &MetricsFloat64{
				Timestamp: ts,
				Values:    loadRaw[ts],
			})
		}

		sort.Slice(run.ClientsLoad, func(i, j int) bool {
			return run.ClientsLoad[i].Timestamp < run.ClientsLoad[j].Timestamp
		})

	} else {

		run.ServersLoad = make([]*MetricsFloat64, 0, len(loadRaw))

		for ts := range loadRaw {

			// Exclude metric for further consideration in
			// case it lies outside our zone of interest.
			if (ts < run.TimestampLowest) || (ts > run.TimestampHighest) {
				continue
			}

			run.ServersLoad = append(run.ServersLoad, &MetricsFloat64{
				Timestamp: ts,
				Values:    loadRaw[ts],
			})
		}

		sort.Slice(run.ServersLoad, func(i, j int) bool {
			return run.ServersLoad[i].Timestamp < run.ServersLoad[j].Timestamp
		})
	}

	return nil
}

/*
func (run *Run) SystemStoreForBoxplots(path string) error {

	sentBytesFile, err := os.OpenFile(filepath.Join(path, "sent-kibytes_per_second.boxplot"), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}
	defer sentBytesFile.Close()
	defer sentBytesFile.Sync()

	for ts := range run.SentBytes {

		var values string
		for i := range run.SentBytes[ts].Values {

			if values == "" {
				values = fmt.Sprintf("%d", run.SentBytes[ts].Values[i])
			} else {
				values = fmt.Sprintf("%s,%d", values, run.SentBytes[ts].Values[i])
			}
		}

		fmt.Fprintln(sentBytesFile, values)
	}

	recvdBytesFile, err := os.OpenFile(filepath.Join(path, "recvd-kibytes_per_second.boxplot"), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}
	defer recvdBytesFile.Close()
	defer recvdBytesFile.Sync()

	for ts := range run.RecvdBytes {

		var values string
		for i := range run.RecvdBytes[ts].Values {

			if values == "" {
				values = fmt.Sprintf("%d", run.RecvdBytes[ts].Values[i])
			} else {
				values = fmt.Sprintf("%s,%d", values, run.RecvdBytes[ts].Values[i])
			}
		}

		fmt.Fprintln(recvdBytesFile, values)
	}

	memoryFile, err := os.OpenFile(filepath.Join(path, "memory_per_second.boxplot"), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}
	defer memoryFile.Close()
	defer memoryFile.Sync()

	for ts := range run.Memory {

		var values string
		for i := range run.Memory[ts].Values {

			if values == "" {
				values = fmt.Sprintf("%f", run.Memory[ts].Values[i])
			} else {
				values = fmt.Sprintf("%s,%f", values, run.Memory[ts].Values[i])
			}
		}

		fmt.Fprintln(memoryFile, values)
	}

	loadFile, err := os.OpenFile(filepath.Join(path, "load_per_second.boxplot"), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}
	defer loadFile.Close()
	defer loadFile.Sync()

	for ts := range run.Load {

		var values string
		for i := range run.Load[ts].Values {

			if values == "" {
				values = fmt.Sprintf("%f", run.Load[ts].Values[i])
			} else {
				values = fmt.Sprintf("%s,%f", values, run.Load[ts].Values[i])
			}
		}

		fmt.Fprintln(loadFile, values)
	}

	return nil
}
*/
