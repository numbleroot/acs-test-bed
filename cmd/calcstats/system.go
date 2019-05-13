package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
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
	Memory        []*MetricsInt64
	MemoryRaw     map[int64][]int64
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

		// Append to corresponding slice of values.
		sysM.RecvdBytesRaw[timestamp] = append(sysM.RecvdBytesRaw[timestamp], value)
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

		memTotal, err := strconv.ParseInt(strings.TrimPrefix(metric[1], "totalKB:"), 10, 64)
		if err != nil {
			return err
		}

		memAvail, err := strconv.ParseInt(strings.TrimPrefix(metric[2], "availKB:"), 10, 64)
		if err != nil {
			return err
		}

		// Calculate difference ("used" memory metric).
		memUsed := memTotal - memAvail

		// Append to corresponding slice of values.
		sysM.MemoryRaw[timestamp] = append(sysM.MemoryRaw[timestamp], memUsed)
	}

	return nil
}

func (sysM *SystemMetrics) OrderSystemMetrics() error {

	sysM.SentBytes = make([]*MetricsInt64, 0, len(sysM.SentBytesRaw))
	sysM.RecvdBytes = make([]*MetricsInt64, 0, len(sysM.RecvdBytesRaw))
	sysM.Memory = make([]*MetricsInt64, 0, len(sysM.MemoryRaw))
	sysM.Load = make([]*MetricsFloat64, 0, len(sysM.LoadRaw))

	// Insert metric values into slices for sorting.

	for ts := range sysM.SentBytesRaw {

		sysM.SentBytes = append(sysM.SentBytes, &MetricsInt64{
			Timestamp: ts,
			Values:    sysM.SentBytesRaw[ts],
		})
	}

	for ts := range sysM.RecvdBytesRaw {

		sysM.RecvdBytes = append(sysM.RecvdBytes, &MetricsInt64{
			Timestamp: ts,
			Values:    sysM.RecvdBytesRaw[ts],
		})
	}

	for ts := range sysM.MemoryRaw {

		sysM.Memory = append(sysM.Memory, &MetricsInt64{
			Timestamp: ts,
			Values:    sysM.MemoryRaw[ts],
		})
	}

	for ts := range sysM.LoadRaw {

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

	/*
		fmt.Printf("Sent bytes:\n")
		for i := range sysM.SentBytes {
			fmt.Printf("\t%03d:  %03d  =>  %#v\n", i, sysM.SentBytes[i].Timestamp, sysM.SentBytes[i].Values)
		}

		fmt.Printf("\nRecvd bytes:\n")
		for i := range sysM.RecvdBytes {
			fmt.Printf("\t%03d:  %03d  =>  %#v\n", i, sysM.RecvdBytes[i].Timestamp, sysM.RecvdBytes[i].Values)
		}

		fmt.Printf("\nMemory:\n")
		for i := range sysM.Memory {
			fmt.Printf("\t%03d:  %03d  =>  %#v\n", i, sysM.Memory[i].Timestamp, sysM.Memory[i].Values)
		}

		fmt.Printf("\nLoad:\n")
		for i := range sysM.Load {
			fmt.Printf("\t%03d:  %03d  =>  %#v\n", i, sysM.Load[i].Timestamp, sysM.Load[i].Values)
		}
	*/

	return nil
}

func (sysM *SystemMetrics) CalcAndWriteLoad(metricsPath string) error {

	first := true
	numValues := -1
	for i := range sysM.Load {

		if first {
			numValues = len(sysM.Load[i].Values)
			first = false
		}

		if len(sysM.Load[i].Values) != numValues {
			fmt.Printf("WARNING: unequal amount of values per timestamp (expected: %d, saw: %d at %d)\n", numValues, len(sysM.Load[i].Values), i)
		}
	}

	return nil
}
