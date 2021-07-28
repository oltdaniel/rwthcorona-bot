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

const DATA_URL = "https://www.lzg.nrw.de/covid19/daten/covid19_alter_5334.csv"

func DataCrawler() {
	sched := clockwork.NewScheduler()

	// structure allows crawling until new dataset is there
	sched.Schedule().Every().Day().At("01:00").Do(crawlForNewDataset)

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

var DATA_BASE = os.Getenv("DATA_DIR")

func getFullDataFilepath(filename string) string {
	return filepath.Join(DATA_BASE, filename)
}

func crawlForNewDataset() {
	// create name and check if already exists
	y, m, d := time.Now().AddDate(0, -1, 0).Date()
	targetFilepath := getFullDataFilepath(
		fmt.Sprintf("data.%02d-%02d-%02d.csv", y, int(m), d))
	if _, err := os.Stat(targetFilepath); os.IsExist(err) {
		fmt.Println("Data already loaded")
		return
	}
	// load csv
	resp, err := http.Get(DATA_URL)
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
	DataConverterAnnouncement <- targetFilepath
}
