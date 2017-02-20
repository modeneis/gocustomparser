package gocustomparser

import (
	"encoding/csv"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)


// Method to get CSV data parsed into struct
func (parser CustomParser) GetCSVData(f interface{}, decoded []byte, n int) (parsed []interface{}, err error) {
	var csvReader = csv.NewReader(strings.NewReader(string(decoded[:n])))
	csvReader.Comma = parser.Separator

	resultType := reflect.ValueOf(f).Type()

	if parser.SkipFirstLine {
		csvReader.Read()
	}

	var rawLine []string
	for {
		rawLine, err = csvReader.Read()
		if err != nil {
			if fmt.Sprint(err) == "EOF" {
				break
			} else {
				return
			}
		}

		var newResult = reflect.New(resultType).Interface()

		// set all the struct fields
		for fieldIndex := 0; fieldIndex < resultType.NumField(); fieldIndex++ {
			var currentField = resultType.Field(fieldIndex)

			var csvTag = currentField.Tag.Get("csv")
			var csvColumnIndex, csvTagErr = strconv.Atoi(csvTag)

			if csvTagErr != nil {
				if csvTag == "" {
					csvColumnIndex = fieldIndex
				} else {
					return nil, csvTagErr
				}
			}

			if csvColumnIndex < 0 {
				err = fmt.Errorf("csv tag in struct field %v is less than zero", currentField.Name)
				return
			}

			if csvColumnIndex >= len(rawLine) {
				err = fmt.Errorf("Trying to access csv column %v for field %v, but csv has only %v column(s)", csvColumnIndex, currentField.Name, len(rawLine))
				return
			}

			var csvElement = rawLine[csvColumnIndex]
			var settableField = reflect.ValueOf(newResult).Elem().FieldByName(currentField.Name)

			if csvElement == "" && parser.SkipEmptyValues {
				continue
			}

			switch currentField.Type.Name() {

			case "bool":
				var parsedBool, err = strconv.ParseBool(csvElement)
				if err != nil {
					return nil,err
				}
				settableField.SetBool(parsedBool)

			case "uint", "uint8", "uint16", "uint32", "uint64":
				var parsedUint, err = strconv.ParseUint(csvElement, 10, 64)
				if err != nil {
					return nil,err
				}
				settableField.SetUint(uint64(parsedUint))

			case "int", "int32", "int64":
				var parsedInt, err = strconv.Atoi(csvElement)
				if err != nil {
					return nil,err
				}
				settableField.SetInt(int64(parsedInt))

			case "float32":
				var parsedFloat, err = strconv.ParseFloat(csvElement, 32)
				if err != nil {
					return nil, err
				}
				settableField.SetFloat(parsedFloat)

			case "float64":
				var parsedFloat, err = strconv.ParseFloat(csvElement, 64)
				if err != nil {
					return nil,err
				}
				settableField.SetFloat(parsedFloat)

			case "string":
				settableField.SetString(csvElement)

			case "Time":
				var date, err = time.Parse(currentField.Tag.Get("csvDate"), csvElement)
				if err != nil {
					return nil,err
				}
				settableField.Set(reflect.ValueOf(date))
			}
		}

		parsed = append(parsed, newResult)
	}

	return
}
