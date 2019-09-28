package main

import (
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
)

// Exp contains all information relevant
// for monitoring an experiment.
type Exp struct {
	System                 string          `json:"system"`
	ZonesNetTroublesIfUsed map[string]bool `json:"zonesNetTroublesIfUsed"`
	Servers                []Worker        `json:"servers"`
	Clients                []Worker        `json:"clients"`
}

// Worker describes one compute instance
// exhaustively for reproducibility.
type Worker struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Zone           string `json:"zone"`
	MinCPUPlatform string `json:"minCPUPlatform"`
	MachineType    string `json:"machineType"`
	TypeOfNode     string `json:"typeOfNode"`
	BinaryName     string `json:"binaryName"`
	SourceImage    string `json:"sourceImage"`
	DiskType       string `json:"diskType"`
	DiskSize       string `json:"diskSize"`
}

// GCloudZones holds all but two GCloud zones.
var GCloudZones = [17]string{
	"asia-east1-b",
	"asia-northeast1-b",
	"asia-southeast1-b",
	"europe-west1-b",
	"europe-west4-b",
	//"us-central1-b",
	"us-east1-b",
	"us-east4-b",
	"us-west1-b",
	// Low storage zones below
	"asia-east2-b",
	"asia-northeast2-b",
	"asia-south1-b",
	"australia-southeast1-b",
	"europe-north1-b",
	"europe-west2-b",
	"europe-west6-b",
	"northamerica-northeast1-b",
	"us-west2-b",
}

// GCloudZonesLowStorage defines limits in terms
// of number of instances with a standard disk
// size of 10GB that can be spawned in a GCloud
// zone without risking overstepping the quota.
// This might change if a GCP project is assigned
// more lenient quotas.
var GCloudZonesLowStorage = map[string]int{
	"asia-east2-b":              290,
	"asia-northeast2-b":         290,
	"asia-south1-b":             290,
	"australia-southeast1-b":    290,
	"europe-north1-b":           290,
	"europe-west2-b":            290,
	"europe-west6-b":            290,
	"northamerica-northeast1-b": 290,
	"us-west2-b":                290,
}

func shuffleZones() {

	// Shuffle zones slice.
	for i := (len(GCloudZones) - 1); i > 0; i-- {

		// Generate new CSPRNG number smaller than i.
		jBig, err := rand.Int(rand.Reader, big.NewInt(int64(i)))
		if err != nil {
			fmt.Printf("Drawing random number failed: %v\n", err)
			os.Exit(1)
		}
		j := int(jBig.Int64())

		// Swap places i and j in zones slice.
		GCloudZones[i], GCloudZones[j] = GCloudZones[j], GCloudZones[i]
	}
}

func pickZone(zoneIdx int, gcloudZonesLowStorage map[string]int) (string, int) {

	// Pick next zone from randomized zones array.
	zone := GCloudZones[zoneIdx]

	// In case the storage quota for this zone has
	// been reached, pick the next zone.
	for gcloudZonesLowStorage[zone] > GCloudZonesLowStorage[zone] {

		// Increment counter. If we traversed zones array
		// once, shuffle it again and reset index.
		zoneIdx++
		if zoneIdx == len(GCloudZones) {
			shuffleZones()
			zoneIdx = 0
		}

		zone = GCloudZones[zoneIdx]
	}

	// Increment local map counter if this is a zone
	// constrained with persistent disk space.
	if GCloudZonesLowStorage[zone] == 290 {
		gcloudZonesLowStorage[zone]++
	}

	zoneIdx++
	if zoneIdx == len(GCloudZones) {
		shuffleZones()
		zoneIdx = 0
	}

	return zone, zoneIdx
}

func main() {

	// Allow for control via command-line flags.
	configsPathFlag := flag.String("configsPath", "./gcloud-configs/", "Specify file system location where GCloud Compute configurations are supposed to be saved.")
	numClientsToGenFlag := flag.Int("numClientsToGen", 1000, "Specify the number of client nodes to generate. Number of conversing clients will be ten times as much.")
	numVuvuzelaMixesToGenFlag := flag.Int("numVuvuzelaMixesToGen", 7, "Specify the number of vuvuzela mix nodes to generate (number of zeno mixes is twice this number minus 1).")
	numZenoCascadesFlag := flag.Int("numZenoCascades", 1, "Specify the number of cascades in zeno to generate.")
	flag.Parse()

	numClientsToGen := *numClientsToGenFlag
	numVuvuzelaMixesToGen := *numVuvuzelaMixesToGenFlag
	numZenoCascades := *numZenoCascadesFlag

	configsPath, err := filepath.Abs(*configsPathFlag)
	if err != nil {
		fmt.Printf("Provided path to configuration files '%s' could not be converted to absolute path: %v\n", *configsPathFlag, err)
		os.Exit(1)
	}

	// Create configuration files folder
	// if it does not exist.
	err = os.MkdirAll(configsPath, 0755)
	if err != nil {
		fmt.Printf("Failed to create configurations folder %s: %v\n", configsPath, err)
		os.Exit(1)
	}

	// Create local map that tracks how many
	// instances have already been assigned
	// to each of the low storage zones.
	gcloudZonesLowStorage := map[string]int{
		"asia-east2-b":              0,
		"asia-northeast2-b":         0,
		"asia-south1-b":             0,
		"australia-southeast1-b":    0,
		"europe-north1-b":           0,
		"europe-west2-b":            0,
		"europe-west6-b":            0,
		"northamerica-northeast1-b": 0,
		"us-west2-b":                0,
	}

	// Shuffle zones array.
	shuffleZones()
	zoneIdx := 0

	// Pick three zones that will experience
	// network troubles if runexperiments is
	// instructed to emulate them.
	netTroubleZones := make(map[string]bool)
	netTroubleZones[GCloudZones[0]] = true
	netTroubleZones[GCloudZones[1]] = true
	netTroubleZones[GCloudZones[2]] = true

	// Prepare structures for each system.

	zenoExp := &Exp{
		System:                 "zeno",
		ZonesNetTroublesIfUsed: netTroubleZones,
		Servers:                make([]Worker, (numZenoCascades * ((2 * numVuvuzelaMixesToGen) - 1))),
		Clients:                make([]Worker, numClientsToGen),
	}

	vuvuzelaExp := &Exp{
		System:                 "vuvuzela",
		ZonesNetTroublesIfUsed: netTroubleZones,
		Servers:                make([]Worker, numVuvuzelaMixesToGen),
		Clients:                make([]Worker, numClientsToGen),
	}

	pungExp := &Exp{
		System:                 "pung",
		ZonesNetTroublesIfUsed: netTroubleZones,
		Servers:                make([]Worker, 1),
		Clients:                make([]Worker, numClientsToGen),
	}

	for i := range zenoExp.Clients {

		zone := ""
		zone, zoneIdx = pickZone(zoneIdx, gcloudZonesLowStorage)

		zenoExp.Clients[i] = Worker{
			ID:             (i + 1),
			Name:           fmt.Sprintf("client-%05d", (i + 1)),
			Zone:           zone,
			MinCPUPlatform: "Intel Skylake",
			MachineType:    "n1-standard-4",
			TypeOfNode:     "client",
			BinaryName:     "zeno",
			SourceImage:    "acs-eval",
			DiskType:       "pd-ssd",
			DiskSize:       "10",
		}
	}

	for i := range zenoExp.Servers {

		zone := ""
		zone, zoneIdx = pickZone(zoneIdx, gcloudZonesLowStorage)

		zenoExp.Servers[i] = Worker{
			ID:             (i + 1),
			Name:           fmt.Sprintf("server-%05d", (i + 1)),
			Zone:           zone,
			MinCPUPlatform: "Intel Skylake",
			MachineType:    "n1-standard-4",
			TypeOfNode:     "server",
			BinaryName:     "zeno",
			SourceImage:    "acs-eval",
			DiskType:       "pd-ssd",
			DiskSize:       "10",
		}
	}

	copy(vuvuzelaExp.Clients, zenoExp.Clients)
	for i := range vuvuzelaExp.Clients {
		vuvuzelaExp.Clients[i].BinaryName = "vuvuzela-client"
	}

	copy(pungExp.Clients, zenoExp.Clients)
	for i := range vuvuzelaExp.Clients {
		vuvuzelaExp.Clients[i].BinaryName = "pung-client"
	}

	copy(vuvuzelaExp.Servers, zenoExp.Servers[:numVuvuzelaMixesToGen])
	for i := range vuvuzelaExp.Servers {

		if i == 0 {
			vuvuzelaExp.Servers[i].TypeOfNode = "coordinator"
			vuvuzelaExp.Servers[i].BinaryName = "vuvuzela-coordinator"
		} else {
			vuvuzelaExp.Servers[i].BinaryName = "vuvuzela-mix"
		}
	}

	copy(pungExp.Servers, zenoExp.Servers[:1])
	for i := range pungExp.Servers {
		pungExp.Servers[i].BinaryName = "pung-server"
	}

	// Marshal slice of zeno experiment to JSON.
	zenoExpJSON, err := json.MarshalIndent(zenoExp, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshal zeno experiment to JSON: %v\n", err)
		os.Exit(1)
	}

	// Write zeno experiment to file.
	err = ioutil.WriteFile(filepath.Join(configsPath, "zeno.json"), zenoExpJSON, 0644)
	if err != nil {
		fmt.Printf("Error writing zeno experiment in JSON format to file: %v\n", err)
		os.Exit(1)
	}

	// Marshal slice of Vuvuzela experiment to JSON.
	vuvuzelaExpJSON, err := json.MarshalIndent(vuvuzelaExp, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshal Vuvuzela experiment to JSON: %v\n", err)
		os.Exit(1)
	}

	// Write Vuvuzela experiment to file.
	err = ioutil.WriteFile(filepath.Join(configsPath, "vuvuzela.json"), vuvuzelaExpJSON, 0644)
	if err != nil {
		fmt.Printf("Error writing Vuvuzela experiment in JSON format to file: %v\n", err)
		os.Exit(1)
	}

	// Marshal slice of pung experiment to JSON.
	pungExpJSON, err := json.MarshalIndent(pungExp, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshal Pung experiment to JSON: %v\n", err)
		os.Exit(1)
	}

	// Write pung experiment to file.
	err = ioutil.WriteFile(filepath.Join(configsPath, "pung.json"), pungExpJSON, 0644)
	if err != nil {
		fmt.Printf("Error writing Pung experiment in JSON format to file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("All done!\n")
}
