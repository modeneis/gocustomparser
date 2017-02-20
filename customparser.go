package gocustomparser

import (
	"os"
	"io/ioutil"
	"fmt"
	"path/filepath"
	"log"
	"bytes"
)

func (parser CustomParser) Parse(f interface{}) (results []interface{}, error error) {
	csvFile, err := os.Open(parser.File)
	if err != nil {
		log.Println("Error opening file:", parser.File, err)
		return nil, err
	}
	defer csvFile.Close()
	fileBytes, err := ioutil.ReadAll(csvFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading from CSV file", err)
		os.Exit(1)
	}
	decoded, err := parser.CustomDecoder.Bytes(fileBytes)
	if err != nil {
		log.Println("Error converting from ISO-8859 to UTF-8:", err)
	}
	n := bytes.IndexByte(fileBytes, 0)
	if n == -1 {
		n = len(fileBytes)
	}

	extension := filepath.Ext(parser.File)
	switch extension {
	case ".csv":
		results, err = parser.GetCSVData(f, decoded, n)
	case ".prn":
		results, err = parser.GetPRNData(f, decoded, n)
	default:
		log.Println("Supports 'csv' and 'prn' files, This extension is not yet implemented ", extension)
	}

	return
}
