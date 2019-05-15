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

type MetricsInt64 struct {
	Timestamp int64
	Values    []int64
}

type MetricsFloat64 struct {
	Timestamp int64
	Values    []float64
}

type SystemMetrics struct {
	SentBytes     []*MetricsInt64
	SentBytesRaw  map[int64][]int64
	RecvdBytes    []*MetricsInt64
	RecvdBytesRaw map[int64][]int64
	Memory        []*MetricsFloat64
	MemoryRaw     map[int64][]float64
	Load          []*MetricsFloat64
	LoadRaw       map[int64][]float64
}

func (sysM *SystemMetrics) AddSentBytes(path string) error {

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
		sysM.SentBytesRaw[timestamp] = append(sysM.SentBytesRaw[timestamp], value)
	}

	return nil
}

func (sysM *SystemMetrics) AddRecvdBytes(path string) error {

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
		sysM.RecvdBytesRaw[timestamp] = append(sysM.RecvdBytesRaw[timestamp], value)
	}

	return nil
}

func (sysM *SystemMetrics) AddMem(path string) error {

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
		sysM.MemoryRaw[timestamp] = append(sysM.MemoryRaw[timestamp], memUsedRatio)
	}

	return nil
}

func (sysM *SystemMetrics) AddLoad(path string) error {

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
		sysM.LoadRaw[timestamp] = append(sysM.LoadRaw[timestamp], loadBusy)
	}

	return nil
}

func (sysM *SystemMetrics) SystemSortByTimestamp(TimestampLowerBound int64, TimestampUpperBound int64) error {

	sysM.SentBytes = make([]*MetricsInt64, 0, len(sysM.SentBytesRaw))
	sysM.RecvdBytes = make([]*MetricsInt64, 0, len(sysM.RecvdBytesRaw))
	sysM.Memory = make([]*MetricsFloat64, 0, len(sysM.MemoryRaw))
	sysM.Load = make([]*MetricsFloat64, 0, len(sysM.LoadRaw))

	// Insert metric values into slices for sorting.

	for ts := range sysM.SentBytesRaw {

		// Exclude metric for further consideration in
		// case it lies outside our zone of interest.
		if (ts < TimestampLowerBound) || (ts > TimestampUpperBound) {
			continue
		}

		sysM.SentBytes = append(sysM.SentBytes, &MetricsInt64{
			Timestamp: ts,
			Values:    sysM.SentBytesRaw[ts],
		})
	}

	for ts := range sysM.RecvdBytesRaw {

		// Exclude metric for further consideration in
		// case it lies outside our zone of interest.
		if (ts < TimestampLowerBound) || (ts > TimestampUpperBound) {
			continue
		}

		sysM.RecvdBytes = append(sysM.RecvdBytes, &MetricsInt64{
			Timestamp: ts,
			Values:    sysM.RecvdBytesRaw[ts],
		})
	}

	for ts := range sysM.MemoryRaw {

		// Exclude metric for further consideration in
		// case it lies outside our zone of interest.
		if (ts < TimestampLowerBound) || (ts > TimestampUpperBound) {
			continue
		}

		sysM.Memory = append(sysM.Memory, &MetricsFloat64{
			Timestamp: ts,
			Values:    sysM.MemoryRaw[ts],
		})
	}

	for ts := range sysM.LoadRaw {

		// Exclude metric for further consideration in
		// case it lies outside our zone of interest.
		if (ts < TimestampLowerBound) || (ts > TimestampUpperBound) {
			continue
		}

		sysM.Load = append(sysM.Load, &MetricsFloat64{
			Timestamp: ts,
			Values:    sysM.LoadRaw[ts],
		})
	}

	// Sort resulting slices by timestamp.

	sort.Slice(sysM.SentBytes, func(i, j int) bool {
		return sysM.SentBytes[i].Timestamp < sysM.SentBytes[j].Timestamp
	})

	sort.Slice(sysM.RecvdBytes, func(i, j int) bool {
		return sysM.RecvdBytes[i].Timestamp < sysM.RecvdBytes[j].Timestamp
	})

	sort.Slice(sysM.Memory, func(i, j int) bool {
		return sysM.Memory[i].Timestamp < sysM.Memory[j].Timestamp
	})

	sort.Slice(sysM.Load, func(i, j int) bool {
		return sysM.Load[i].Timestamp < sysM.Load[j].Timestamp
	})

	return nil
}

func (sysM *SystemMetrics) SystemStoreForBoxplots(path string) error {

	sentBytesFile, err := os.OpenFile(filepath.Join(path, "sent-kibytes_per_second.boxplot"), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
	if err != nil {
		return err
	}
	defer sentBytesFile.Close()
	defer sentBytesFile.Sync()

	for ts := range sysM.SentBytes {

		var values string
		for i := range sysM.SentBytes[ts].Values {

			if values == "" {
				values = fmt.Sprintf("%d", sysM.SentBytes[ts].Values[i])
			} else {
				values = fmt.Sprintf("%s,%d", values, sysM.SentBytes[ts].Values[i])
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

	for ts := range sysM.RecvdBytes {

		var values string
		for i := range sysM.RecvdBytes[ts].Values {

			if values == "" {
				values = fmt.Sprintf("%d", sysM.RecvdBytes[ts].Values[i])
			} else {
				values = fmt.Sprintf("%s,%d", values, sysM.RecvdBytes[ts].Values[i])
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

	for ts := range sysM.Memory {

		var values string
		for i := range sysM.Memory[ts].Values {

			if values == "" {
				values = fmt.Sprintf("%f", sysM.Memory[ts].Values[i])
			} else {
				values = fmt.Sprintf("%s,%f", values, sysM.Memory[ts].Values[i])
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

	for ts := range sysM.Load {

		var values string
		for i := range sysM.Load[ts].Values {

			if values == "" {
				values = fmt.Sprintf("%f", sysM.Load[ts].Values[i])
			} else {
				values = fmt.Sprintf("%s,%f", values, sysM.Load[ts].Values[i])
			}
		}

		fmt.Fprintln(loadFile, values)
	}

	return nil
}
