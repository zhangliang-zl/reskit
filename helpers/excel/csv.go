package excel

import (
	"encoding/csv"
	"os"
)

func WriteToCsv(records [][]string, filename string) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	data := make([][]string, 0)
	data = append(data, []string{"head1", "head2"})
	data = append(data, []string{"v1", "v2"})

	f.WriteString("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(f)
	err = w.WriteAll(records)
	return err
}

func ReadCsv(filename string) (records [][]string, err error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0444)
	if err != nil {
		return
	}

	defer f.Close()
	reader := csv.NewReader(f)
	return reader.ReadAll()
}
