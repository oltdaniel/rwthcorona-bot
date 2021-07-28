package workers

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/oltdaniel/rwth-coronabot/utils"
)

var DataConverterAnnouncement chan (string) = make(chan string)

func DataConverter() {
	for {
		newDataset := <-DataConverterAnnouncement
		fmt.Println(newDataset)

		fp, err := os.Open(newDataset)
		if err != nil {
			panic(err)
		}
		cr := csv.NewReader(fp)
		headers, err := cr.Read()
		if err != nil {
			panic(err)
		}
		data, err := cr.ReadAll()
		if err != nil {
			panic(err)
		}
		err = handleDataset(&headers, &data)
		if err != nil {
			panic(err)
		}
	}
}

func headerIndex(headers *[]string, header string) (int, error) {
	for i, v := range *headers {
		if v == header {
			return i, nil
		}
	}
	return 0, errors.New("not found")
}

func handleDataset(headers *[]string, data *[][]string) error {
	// extract highest date entry
	var maxDate int64
	dateIndex, err := headerIndex(headers, "datumStd")
	if err != nil {
		return err
	}
	for _, v := range *data {
		date := v[dateIndex]
		td, err := time.Parse("2006-01-02", date)
		if err != nil {
			continue
		}
		if unix := td.Unix(); unix > maxDate {
			maxDate = unix
		}
	}
	// check if that is equal to yesterday
	today, err := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	yesterday := today.AddDate(0, -1, 0)
	if err != nil && maxDate >= yesterday.Unix() {
		return errors.New("dataset too old")
	}
	// re-map data
	dataset := make(utils.Dataset)
	// TODO: Custom csv parser using maps
	altersgruppeIndex, err := headerIndex(headers, "altersgruppe")
	if err != nil {
		return err
	}
	anzahlIndex, err := headerIndex(headers, "anzahlM7Tage")
	if err != nil {
		return err
	}
	rateIndex, err := headerIndex(headers, "rateM7Tage")
	if err != nil {
		return err
	}
	anteilIndex, err := headerIndex(headers, "anteilM7Tage")
	if err != nil {
		return err
	}
	// iterate through each line
	for _, v := range *data {
		// get group details
		date := v[dateIndex]
		altersgruppe := v[altersgruppeIndex]
		// convert to floats
		anzahlWoche, err := strconv.ParseFloat(v[anzahlIndex], 64)
		if err != nil {
			return err
		}
		rateWoche, err := strconv.ParseFloat(v[rateIndex], 64)
		if err != nil {
			return err
		}
		anteilWoche, err := strconv.ParseFloat(v[anteilIndex], 64)
		if err != nil {
			return err
		}
		// append data
		if dataset[date] == nil {
			dataset[date] = make(map[string][]*utils.DatasetEntry)
		}
		dataset[date][altersgruppe] = append(dataset[date][altersgruppe], &utils.DatasetEntry{
			AnzahlWoche:  anzahlWoche,
			RateWoche:    rateWoche,
			AnteiltWoche: anteilWoche,
		})
	}

	// store dataset
	utils.DATASET.Update(&dataset)

	return nil
}
