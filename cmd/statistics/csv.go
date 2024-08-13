package statistics

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

const filePath = "cmd/statistics/file.csv"

func storeAsCsv(value interface{}) error {
	t := reflect.TypeOf(value).Elem()
	v := reflect.ValueOf(value).Elem()
	headerRow := make([]string, t.NumField())
	valuesRow := make([]string, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		headerRow[i] = t.Field(i).Name
		valuesRow[i] = fmt.Sprintf("%.2f", v.Field(i).Interface())
	}
	return write([][]string{headerRow, valuesRow})
}

func readCsv() ([]*statistics, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	stats := make([]*statistics, 0, len(records)-1)
	for i := 1; i < len(records); i++ {
		numVals, err := strconv.ParseFloat(records[i][0], 64)
		minF, err := strconv.ParseFloat(records[i][1], 64)
		maxF, err := strconv.ParseFloat(records[i][2], 64)
		mean, err := strconv.ParseFloat(records[i][3], 64)
		std, err := strconv.ParseFloat(records[i][4], 64)
		if err != nil {
			return nil, err
		}
		stats = append(stats, &statistics{numVals, minF, maxF, mean, std})
	}
	return stats, nil
}

func write(rows [][]string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for i, row := range rows {
		if stat, err := file.Stat(); i == 0 && err == nil && stat.Size() != 0 {
			continue
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	return nil
}
