package excel

import (
	"os"
	"testing"
)

func TestXls(t *testing.T) {
	path := "test.xlsx"

	var data = [][]string{
		{"head1", "head2"},
		{"v1a", "v1b"},
		{"v2a", "v2b"},
		{"v31", "v3b"},
	}

	if err := WriteToExcel(data, path); err != nil {
		t.Errorf("writeXlsx err %v", err)
	}
	dataRead, err := ReadExcel(path)
	if err != nil {
		t.Errorf("readcsv err %v", err)
	}

	if dataRead[1][0] != "v1a" || dataRead[2][1] != "v2b" {
		t.Error("read xlsx data error")
	}
	os.Remove(path)
}
