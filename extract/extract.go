package extract

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/shamhub/exercise/data"
)

type TimeSeriesJSONReader struct {
	filePath string
	file     *os.File
	scanner  *bufio.Scanner
}

func NewJSONReader(filePath string) *TimeSeriesJSONReader {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)

	return &TimeSeriesJSONReader{
		filePath: filePath,
		file:     file,
		scanner:  scanner,
	}
}

// Read a line and return the data
func (t *TimeSeriesJSONReader) ReadEntry() *data.TimeSeriesData {

	defer t.file.Close()
	var data data.TimeSeriesData

	// optionally, resize scanner's capacity for lines over 64K, see next example
	if t.scanner.Scan() {
		record := t.scanner.Text()
		json.Unmarshal([]byte(record), &data)
		fmt.Println(data)
	}

	if err := t.scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return &data
}
