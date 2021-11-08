package excel

import (
	"os"
	"testing"
)

func TestCsv(t *testing.T) {
	path := "test.csv"

	var data = [][]string{
		{"head1", "head2"},
		{"v1a", "v1b"},
		{"v2a", "v2b"},
		{"v31", "v3b"},
	}

	if err := WriteToCsv(data, path); err != nil {
		t.Errorf("writeCsv err %v", err)
	}
	dataRead, err := ReadCsv(path)
	if err != nil {
		t.Errorf("readcsv err %v", err)
	}

	if dataRead[1][0] != "v1a" || dataRead[2][1] != "v2b" {
		t.Error("read csv data error")
	}

	os.Remove(path)
}
