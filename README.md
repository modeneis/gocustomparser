## Small Go package to parse CSV and PRN files with type checking and struct transformation.


## Go Get module
`go get -u github.com/modeneis/gocustomparser`


# Tests:
For tests, see prnparser_test.go and csvparser_test.go

# USAGE CSV
```
	type Workbook struct {
		Name        string `csv:"0" json:"name"`
		Address     string `csv:"1" json:"address"`
		Postcode    string `csv:"2" json:"postcode"`
		Phone       string `csv:"3" json:"phone"`
		CreditLimit string `csv:"4" json:"creditlimit"`
		Birthday    string `csv:"5" json:"birthday"`
	}
	csvParser := gocustomparser.CustomParser{
    		File:          "testfiles/Workbook2.csv",
            CustomDecoder: charmap.ISO8859_1.NewDecoder(),
    		Separator:     ',',
    		SkipFirstLine: true,
    }
    var csvParsedItems, err = csvParser.Parse(Workbook{})
```


# USAGE PRN
```
    type Workbook struct {
		Name        string `prn:"0" json:"name"`
		Address     string `prn:"1" json:"address"`
		Postcode    string `prn:"2" json:"postcode"`
		Phone       string `prn:"3" json:"phone"`
		CreditLimit string `prn:"4" json:"creditlimit"`
		Birthday    string `prn:"5" json:"birthday"`
	}

    prnParser := gocustomparser.CustomParser{
		File:          "testfiles/Workbook2.prn",
        CustomDecoder: charmap.ISO8859_1.NewDecoder(),
		SkipFirstLine: true,
		PRNReader: func(raw string) (line []string, err error) {
			runes := []rune(raw)
			if len(runes) < 74 {
				err = fmt.Errorf("ReadPrnLine detected Wrong data -> %s", raw)
				return
			}
			line = append(line, strings.TrimSpace(string(runes[0:16])))
			line = append(line, strings.TrimSpace(string(runes[16:38])))
			line = append(line, strings.TrimSpace(string(runes[38:47])))
			line = append(line, strings.TrimSpace(string(runes[47:61])))
			line = append(line, strings.TrimSpace(string(runes[61:74])))
			line = append(line, strings.TrimSpace(string(runes[74:])))
			return
		},
	}
	var prnParsedItems, err = prnParser.Parse(Workbook{})
```