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

func (run *Run) AddMsgsPerMix(runServersPath string) error {

	err := filepath.Walk(runServersPath, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if filepath.Base(path) == "pool-sizes_round.evaluation" {

			// Extract this mix' name.
			name := ""
			pathParts := strings.Split(path, "/")
			for i := range pathParts {

				if strings.HasPrefix(pathParts[i], "mixnet-") {
					name = strings.Split(pathParts[i], "_")[0]
				}
			}

			// Ingest supplied metrics file.
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			content = bytes.TrimSpace(content)

			// Split file contents into lines.
			lines := strings.Split(string(content), "\n")

			// Prepare messages state object.
			numMsgs := make([]int64, len(lines))

			for i := range lines {

				// Split line at whitespace characters.
				metric := strings.Fields(lines[i])

				// Convert second to fifth element to numbers.
				firstPool, err := strconv.ParseInt(strings.TrimPrefix(metric[1], "1st:"), 10, 64)
				if err != nil {
					return err
				}

				secPool, err := strconv.ParseInt(strings.TrimPrefix(metric[2], "2nd:"), 10, 64)
				if err != nil {
					return err
				}

				thirdPool, err := strconv.ParseInt(strings.TrimPrefix(metric[3], "3rd:"), 10, 64)
				if err != nil {
					return err
				}

				outPool, err := strconv.ParseInt(strings.TrimPrefix(metric[4], "out:"), 10, 64)
				if err != nil {
					return err
				}

				// Add to slice of message counts.
				numMsgs[i] = (firstPool + secPool + thirdPool + outPool)
			}

			// Add data of mix to global state.
			run.Mixes = append(run.Mixes, name)
			run.MsgsPerMix = append(run.MsgsPerMix, numMsgs)
		}

		return nil
	})

	return err
}

func (set *Setting) MsgsPerMixToFile(path string) error {

	if strings.Contains(path, "/zeno/") {

		for i := range set.Runs {

			runPlotPath := filepath.Join(path, (fmt.Sprintf("run-%02d", (i + 1))), "msgs-per-mix_first-to-last-round.data")

			msgsPerMixFile, err := os.OpenFile(runPlotPath, (os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_APPEND), 0644)
			if err != nil {
				return err
			}
			defer msgsPerMixFile.Close()
			defer msgsPerMixFile.Sync()

			// Prefix list of metrics with labels.
			fmt.Fprintf(msgsPerMixFile, "%s\n", strings.Join(set.Runs[i].Mixes, ","))

			for j := range set.Runs[i].MsgsPerMix {

				// Start each subsequent line with first metrics
				// value for that corresponding mix.
				fmt.Fprintf(msgsPerMixFile, "%d", set.Runs[i].MsgsPerMix[j][0])

				for k := 1; k < len(set.Runs[i].MsgsPerMix[j]); k++ {
					fmt.Fprintf(msgsPerMixFile, ",%d", set.Runs[i].MsgsPerMix[j][k])
				}

				fmt.Fprintf(msgsPerMixFile, "\n")
			}

			fmt.Fprintf(msgsPerMixFile, "\n")
		}
	}

	return nil
}
