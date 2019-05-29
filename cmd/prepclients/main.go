package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {

	resultsPathFlag := flag.String("resultsPath", "", "Specify file system location where results to prepare are stored.")
	numClientsToPrepFlag := flag.Int("numClientsToPrep", 1000, "Specify the number of client nodes to prepare. Has to be an even number.")
	flag.Parse()

	if *resultsPathFlag == "" {
		fmt.Printf("Supply path to results files ('-resultsPath').\n")
		os.Exit(1)
	}

	resultsPath := *resultsPathFlag
	searchPath := filepath.Join(resultsPath, "clients")
	numClientsToPrep := *numClientsToPrepFlag

	correctClientsFound := 0
	deleteThese := make([]string, 0, 20)

	allFolders, err := ioutil.ReadDir(searchPath)
	if err != nil {
		fmt.Printf("Failed to retrieve all folders: %v\n", err)
		os.Exit(1)
	}

	for i := range allFolders {

		clientIDString := strings.Split(strings.Split(filepath.Base(allFolders[i].Name()), "_")[0], "-")[1]

		clientID, err := strconv.Atoi(clientIDString)
		if err != nil {
			fmt.Printf("Error converting ID to string: %v\n", err)
			os.Exit(1)
		}
		clientName := fmt.Sprintf("client-%05d", clientID)

		clientCands, err := filepath.Glob(fmt.Sprintf("%s/%s_*", searchPath, clientName))
		if err != nil {
			fmt.Printf("Failed retrieving client's folder: %v\n", err)
			os.Exit(1)
		}

		// Determine partner client of this client.
		partner := ""
		if (clientID % 2) == 0 {
			partner = fmt.Sprintf("client-%05d", (clientID - 1))
		} else {
			partner = fmt.Sprintf("client-%05d", (clientID + 1))
		}

		// Attempt to find partner folder. If unavailable,
		// mark client as removable.
		partnerCands, err := filepath.Glob(fmt.Sprintf("%s/%s_*/log.evaluation", searchPath, partner))
		if err != nil {
			fmt.Printf("Failed retrieving partner's folder: %v\n", err)
			os.Exit(1)
		}

		if (len(partnerCands) == 1) && (correctClientsFound < numClientsToPrep) {
			correctClientsFound++
		} else {
			deleteThese = append(deleteThese, clientCands[0])
		}
	}

	fmt.Printf("Delete these client folders:\n%s\n", strings.Join(deleteThese, " "))
}
