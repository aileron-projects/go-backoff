package examples

import (
	"encoding/csv"
	"os"
)

func NewCsvWriter(filename string) (writer *csv.Writer, close func()) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}
	writer = csv.NewWriter(f)
	return writer, func() {
		writer.Flush()
		f.Close()
	}
}
