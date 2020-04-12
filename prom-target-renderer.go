package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

const inputDir = "/volumes_in/"
const outputDir = "/volumes_out/"

type Labels map[string]string

type MinionsLabels map[string]Labels

type Exporters map[string]int

type Minion struct {
	Exporters Exporters `json:"exporters"`
}

type Minions map[string]Minion

type Out struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

func check(e error, herr string) {
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		fmt.Fprintln(os.Stderr, herr)
	}
}

func getDownMinions() []string {
	dat, err := ioutil.ReadFile(inputDir + "down_minions.json")
	check(err, "Can't open down_minions.json")
	var downMinions []string
	err = json.Unmarshal(dat, &downMinions)
	check(err, "Can't parse down_minions.json")
	return downMinions
}

func getExporters() Minions {
	dat, err := ioutil.ReadFile(inputDir + "exporters.json")
	check(err, "Can't open exporters.json")
	var minions Minions
	err = json.Unmarshal(dat, &minions)
	check(err, "")
	return minions
}

func getLabels() MinionsLabels {
	dat, err := ioutil.ReadFile(inputDir + "labels.json")
	check(err, "Can't open labels.json")
	var minionLabels MinionsLabels
	err = json.Unmarshal(dat, &minionLabels)
	check(err, "")
	return minionLabels
}

func renderPrometheusTargets() {

	minions := getExporters()
	minionLabels := getLabels()

	var minionsOut []Out
	var minionOut Out
	for minionName, minion := range minions {
		minionOut.Labels = map[string]string{}
		for exporterName, exporterPort := range minion.Exporters {
			minionOut.Targets = []string{minionName + ":" + strconv.Itoa(exporterPort)}
			minionOut.Labels["minion_id"] = minionName
			minionOut.Labels["exporter_name"] = exporterName
			for key, value := range minionLabels[minionName] {
				if key != "retcode" {
					minionOut.Labels[key] = value
				}
			}
			minionsOut = append(minionsOut, minionOut)
		}
	}

	//add info about unavailable minions
	downMinions := getDownMinions()
	for _, downMinion := range downMinions {
		minionOut.Labels = map[string]string{}
		minionOut.Targets = []string{downMinion + ":9200"}
		minionOut.Labels["minion_id"] = downMinion
		minionOut.Labels["exporter_name"] = "node_exporter"
		minionsOut = append(minionsOut, minionOut)
	}
	b, _ := json.Marshal(minionsOut)
	ioutil.WriteFile(outputDir+"targets.json", b, os.ModePerm)
	//fmt.Fprintf(os.Stdout, "%s", b)
}

func main() {
	fmt.Println("Starting renderer...")
	for {
		renderPrometheusTargets()
		time.Sleep(60 * time.Second)
	}
}
