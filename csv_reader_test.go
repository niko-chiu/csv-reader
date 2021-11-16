package csvreader

import (
	"os"
	"testing"
)

func TestReadFile(t *testing.T) {
	var err error

	f, err := os.Open("./test.csv")
	if err != nil {
		t.Error(err)
	}

	csvData, err := ReadFile(f)
	if err != nil {
		t.Error(err)
	}

	if len(csvData.Header) != 5 {
		t.Errorf("header is not valid")
	}

	if len(csvData.Rows) != 2 {
		t.Errorf("number of rows is not valid")
	}

	var v string

	v, _ = csvData.Get(0, "b")
	if v != "1" {
		t.Errorf("value is not valid")
	}

	v, _ = csvData.Get(0, "c")
	if v != "2" {
		t.Errorf("value is not valid")
	}

	_, err = csvData.Get(3, "a")
	if err == nil {
		t.Errorf("error not occurs")
	}

	_, err = csvData.Get(-1, "a")
	if err == nil {
		t.Errorf("error not occurs")
	}

	_, err = csvData.Get(0, "z")
	if err == nil {
		t.Errorf("error not occurs")
	}
}

func TestScan(t *testing.T) {
	type test struct {
		A string  `csv:"a"`
		B int     `csv:"b"`
		C int64   `csv:"c"`
		D float32 `csv:"d"`
		E float64 `csv:"e"`
	}

	var err error

	f, err := os.Open("./test.csv")
	if err != nil {
		t.Error(err)
	}

	csvData, err := ReadFile(f)
	if err != nil {
		t.Error(err)
	}

	val := test{}
	err = csvData.Scan(3, &val)
	if err == nil {
		t.Errorf("enter row 3 with no error")
	}

	err = csvData.Scan(0, &val)
	if err != nil {
		t.Error(err)
	}

	if val.A != "0" {
		t.Errorf("value is not match")
	}

	if val.B != 1 {
		t.Errorf("value is not int")
	}

	if val.C != 2 {
		t.Errorf("value is not int")
	}

	if val.D != 3.1 {
		t.Errorf("value is not float")
	}

	if val.E != 4.56789 {
		t.Errorf("value is not float")
	}

	err = csvData.Scan(1, &val)
	if err != nil {
		t.Error(err)
	}

	if val.A != "" {
		t.Errorf("value is not match")
	}

	if val.B != 0 {
		t.Errorf("value is not match")
	}

	if val.C != 0 {
		t.Errorf("value is not match")
	}

	if val.D != 0 {
		t.Errorf("value is not match")
	}

	if val.E != 0 {
		t.Errorf("value is not match")
	}
}
