package utils

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var DATASET = NewDatasetStore()

type Dataset map[string]DatasetDay

func (d *Dataset) Yesterday() (*DatasetDay, error) {
	today := time.Now().AddDate(0, -1, 0).Format("2006-01-02")
	fmt.Println(today)
	if day, ok := (*d)[today]; ok {
		return &day, nil
	}
	return nil, errors.New("day missing")
}

type DatasetDay map[string][]*DatasetEntry

func (d *DatasetDay) Total() *DatasetEntry {
	for group, v := range *d {
		if group == "gesamt" {
			if len(v) == 1 {
				return v[0]
			}
		}
	}
	return nil
}

type DatasetEntry struct {
	AnzahlWoche  float64
	RateWoche    float64
	AnteiltWoche float64
}

type DatasetStore struct {
	dataset *Dataset
	lock    sync.Mutex
}

func NewDatasetStore() *DatasetStore {
	datasetStore := &DatasetStore{
		lock:    sync.Mutex{},
		dataset: nil,
	}
	return datasetStore
}

func (d *DatasetStore) Get() *Dataset {
	var dataset *Dataset
	d.lock.Lock()
	dataset = d.dataset
	d.lock.Unlock()
	return dataset
}

func (d *DatasetStore) Update(dataset *Dataset) {
	d.lock.Lock()
	d.dataset = dataset
	d.lock.Unlock()
}
