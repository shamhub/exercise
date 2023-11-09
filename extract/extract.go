package extract

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/shamhub/exercise/data"
)

const filePath = "./data.json"

type TimeSeriesJSONReader struct {
	filePath string
	file     *os.File
	scanner  *bufio.Scanner
}

func NewJSONReader() *TimeSeriesJSONReader {
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

func (t *TimeSeriesJSONReader) ReadEntry() *data.TimeSeriesData {
	// Read line by line from data.json and throw on channel

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
