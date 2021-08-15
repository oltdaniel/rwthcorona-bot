package workers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/whiteshtef/clockwork"
)

var DataCrawlerUpdateRequest chan (string) = make(chan string)

const (
	DATASET_ALTER = iota
	DATASET_BASE
)

var DATASETS = map[int]string{
	DATASET_ALTER: "https://www.lzg.nrw.de/covid19/daten/covid19_alter_5334.csv",
	DATASET_BASE:  "https://www.lzg.nrw.de/covid19/daten/covid19_5334.csv",
}

func DataCrawler() {
	crawlForNewDataset()

	sched := clockwork.NewScheduler()

	// structure allows crawling until new dataset is there
	sched.Schedule().Every().Hour().Do(crawlForNewDataset)

	// "manual" update
	go func() {
		for {
			oldFilepath := <-DataCrawlerUpdateRequest
			_ = os.Remove(oldFilepath)
			go crawlForNewDataset()
		}
	}()

	sched.Run()
}

func crawlForNewDataset() {
	for i, v := range DATASETS {
		crawlDataset(i, v)
	}
}

var DATA_BASE = os.Getenv("DATA_DIR")

func getFullDataFilepath(filename string) string {
	return filepath.Join(DATA_BASE, filename)
}

func crawlDataset(datasetType int, datasetUrl string) {
	// create name and check if already exists
	y, m, d := time.Now().AddDate(0, 0, -1).Date()
	h := time.Now().Hour()
	targetFilepath := getFullDataFilepath(
		fmt.Sprintf("data.type%d.%02d-%02d-%02d_%02d00.csv", datasetType, y, int(m), d, h))
	if _, err := os.Stat(targetFilepath); os.IsExist(err) {
		fmt.Println("Data already loaded")
		return
	}
	// load csv
	resp, err := http.Get(datasetUrl)
	if err != nil {
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = os.WriteFile(targetFilepath, body, 0744)
	if err != nil {
		return
	}
	DataConverterAnnouncement <- DataConverterUpdate{DatasetType: datasetType, DatasetPath: targetFilepath}
}
