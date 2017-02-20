package gocustomparser

import (
	"bufio"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Method to get PRN data parsed into struct
func (parser CustomParser) GetPRNData(f interface{}, decoded []byte, n int) (parsed []interface{}, err error) {
	scanner := bufio.NewScanner(strings.NewReader(string(decoded[:n])))

	if parser.SkipFirstLine {
		scanner.Scan()
	}
	resultType := reflect.ValueOf(f).Type()

	var rawLine []string
	for scanner.Scan() {
		r := scanner.Text()


		rawLine, err = parser.PRNReader(r)
		if err != nil {
			log.Println("Error splitting fields, check input file:", err)
		}

		var newResult = reflect.New(resultType).Interface()

		// set all the struct fields
		for fieldIndex := 0; fieldIndex < resultType.NumField(); fieldIndex++ {
			var currentField = resultType.Field(fieldIndex)

			var prnTag = currentField.Tag.Get("prn")
			var prnColumnIndex, err = strconv.Atoi(prnTag)

			if err != nil {
				if prnTag == "" {
					prnColumnIndex = fieldIndex
				} else {
					return parsed, err
				}
			}

			if prnColumnIndex < 0 {
				err = fmt.Errorf("prn tag in struct field %v is less than zero", currentField.Name)
				return parsed, err
			}

			if prnColumnIndex >= len(rawLine) {
				err = fmt.Errorf("Trying to access prn column %v for field %v, but prn has only %v column(s)", prnColumnIndex, currentField.Name, len(rawLine))
				return parsed, err
			}

			var prnElement = rawLine[prnColumnIndex]
			var settableField = reflect.ValueOf(newResult).Elem().FieldByName(currentField.Name)

			if prnElement == "" && parser.SkipEmptyValues {
				continue
			}

			switch currentField.Type.Name() {

			case "bool":
				var parsedBool, err = strconv.ParseBool(prnElement)
				if err != nil {
					return parsed, err
				}
				settableField.SetBool(parsedBool)

			case "uint", "uint8", "uint16", "uint32", "uint64":
				var parsedUint, err = strconv.ParseUint(prnElement, 10, 64)
				if err != nil {
					return parsed, err
				}
				settableField.SetUint(uint64(parsedUint))

			case "int", "int32", "int64":
				var parsedInt, err = strconv.Atoi(prnElement)
				if err != nil {
					return parsed, err
				}
				settableField.SetInt(int64(parsedInt))

			case "float32":
				var parsedFloat, err = strconv.ParseFloat(prnElement, 32)
				if err != nil {
					return parsed, err
				}
				settableField.SetFloat(parsedFloat)

			case "float64":
				var parsedFloat, err = strconv.ParseFloat(prnElement, 64)
				if err != nil {
					return parsed, err
				}
				settableField.SetFloat(parsedFloat)

			case "string":
				settableField.SetString(prnElement)

			case "Time":
				var date, err = time.Parse(currentField.Tag.Get("prnDate"), prnElement)
				if err != nil {
					return parsed, err
				}
				settableField.Set(reflect.ValueOf(date))
			}
		}

		parsed = append(parsed, newResult)
	}
	if err = scanner.Err(); err != nil {
		log.Println("Error reading from CSV file:", err)
	}
	return
}
