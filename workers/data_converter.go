package workers

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/oltdaniel/rwthcorona-bot/utils"
)

type DataConverterUpdate struct {
	DatasetType int
	DatasetPath string
}

var DataConverterAnnouncement chan (DataConverterUpdate) = make(chan DataConverterUpdate)

func DataConverter() {
	for {
		newDataset := <-DataConverterAnnouncement
		fmt.Println(newDataset)

		fp, err := os.Open(newDataset.DatasetPath)
		if err != nil {
			panic(err)
		}
		// TODO: this is cheating. there are some random 3 bytes at the start
		fp.Seek(3, 0)
		cr := csv.NewReader(fp)
		headers, err := cr.Read()
		if err != nil {
			panic(err)
		}
		data, err := cr.ReadAll()
		if err != nil {
			panic(err)
		}
		err = handleDataset(newDataset.DatasetType, &headers, &data)
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
	return 0, errors.New(fmt.Sprintf("'%v' column not found", header))
}

func handleDataset(datasetType int, headers *[]string, data *[][]string) error {
	switch datasetType {
	case DATASET_ALTER:
		return handleAlterDataset(headers, data)
	case DATASET_BASE:
		return handleBaseDataset(headers, data)
	}
	return nil
}

func handleAlterDataset(headers *[]string, data *[][]string) error {
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
	yesterday := today.AddDate(0, 0, -1)
	if err != nil && maxDate >= yesterday.Unix() {
		return errors.New("dataset too old")
	}
	// TODO: Custom csv parser using maps
	kreisIndex, err := headerIndex(headers, "kreis")
	if err != nil {
		return err
	}
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
		// convert to correct unit
		kreis, err := strconv.ParseInt(v[kreisIndex], 10, 64)
		if err != nil {
			return err
		}
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
		stm, err := utils.DATABASE.Prepare("INSERT OR IGNORE INTO corona_data(tag, plz, label, altersgruppe, anzahlWoche, rateWoche, anteilWoche) VALUES (?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			return err
		}
		// do it the stupid way, we only need to support these two
		label := "Aachen"
		if kreis == 5 {
			label = "NRW"
		}
		_, err = stm.Exec(date, kreis, label, altersgruppe, anzahlWoche, rateWoche, anteilWoche)
		if err != nil {
			return err
		}
	}

	return nil
}

func handleBaseDataset(headers *[]string, data *[][]string) error {
	// extract highest date entry
	var maxDate int64
	dateIndex, err := headerIndex(headers, "datumstd")
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
	yesterday := today.AddDate(0, 0, -1)
	if err != nil && maxDate >= yesterday.Unix() {
		return errors.New("dataset too old")
	}
	// TODO: Custom csv parser using maps
	kreisIndex, err := headerIndex(headers, "kreis")
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
	// iterate through each line
	for _, v := range *data {
		// get group details
		date := v[dateIndex]
		// convert to correct unit
		kreis, err := strconv.ParseInt(v[kreisIndex], 10, 64)
		if err != nil {
			return err
		}
		anzahlWoche, err := strconv.ParseFloat(v[anzahlIndex], 64)
		if err != nil {
			return err
		}
		rateWoche, err := strconv.ParseFloat(v[rateIndex], 64)
		if err != nil {
			return err
		}
		// append data
		stm, err := utils.DATABASE.Prepare("INSERT OR IGNORE INTO corona_data(tag, plz, label, altersgruppe, anzahlWoche, rateWoche, anteilWoche) VALUES (?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			return err
		}
		// do it the stupid way, we only need to support these two
		label := "Aachen"
		if kreis == 5 {
			label = "NRW"
		}
		_, err = stm.Exec(date, kreis, label, "gesamt", anzahlWoche, rateWoche, 100.0)
		if err != nil {
			return err
		}
	}

	return nil
}
