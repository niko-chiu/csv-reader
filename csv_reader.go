package csvreader

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

type CSV struct {
	Header map[string]int
	Rows   [][]string
}

func (csvData CSV) Get(row int, column string) (string, error) {
	if row < 0 || row >= len(csvData.Rows) {
		return "", fmt.Errorf("row number is not valid")
	}

	val, ok := csvData.Header[column]

	if !ok {
		return "", fmt.Errorf("column " + column + " is not exist")
	}

	if val >= len(csvData.Rows[row]) {
		return "", fmt.Errorf("column of the row is not exist")
	}

	return csvData.Rows[row][val], nil
}

func (csvData CSV) Scan(row int, v interface{}) error {
	vType := reflect.TypeOf(v)
	vValue := reflect.ValueOf(v)

	if vType.Kind() != reflect.Ptr {
		return fmt.Errorf("not a pointer")
	}

	if vType.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("not a struct")
	}

	for i := 0; i < vType.Elem().NumField(); i++ {
		fieldType := vType.Elem().Field(i).Type.Kind()
		tag := vType.Elem().Field(i).Tag.Get("csv")

		if tag == "" {
			continue
		}

		value, err := csvData.Get(row, tag)
		if err != nil {
			return err
		}

		switch fieldType {
		case reflect.String:
			vValue.Elem().Field(i).SetString(value)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if value == "" {
				vValue.Elem().Field(i).SetInt(int64(0))
				continue
			}

			val, err := strconv.Atoi(value)
			if err != nil {
				return err
			}

			vValue.Elem().Field(i).SetInt(int64(val))
		case reflect.Float32, reflect.Float64:
			if value == "" {
				vValue.Elem().Field(i).SetFloat(0)
				continue
			}

			val, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return err
			}

			vValue.Elem().Field(i).SetFloat(val)
		case reflect.Bool:
			if value == "" {
				vValue.Elem().Field(i).SetBool(false)
				continue
			}

			val, err := strconv.ParseBool(value)
			if err != nil {
				return err
			}

			vValue.Elem().Field(i).SetBool(val)
		}
	}

	return nil
}

func ReadFile(csvFile *os.File) (CSV, error) {
	csvData := CSV{
		Header: map[string]int{},
	}

	lines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		return csvData, err
	}

	csvData.Rows = lines[1:]

	for i := 0; i < len(lines[0]); i++ {
		csvData.Header[lines[0][i]] = i
	}

	return csvData, nil
}
