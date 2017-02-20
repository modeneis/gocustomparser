package gocustomparser_test

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gocustomparser"
	"golang.org/x/text/encoding/charmap"
	"strings"
	"testing"
)

func TestParsePRN(t *testing.T) {

	type Workbook struct {
		Name        string `prn:"0" json:"name"`
		Address     string `prn:"1" json:"address"`
		Postcode    string `prn:"2" json:"postcode"`
		Phone       string `prn:"3" json:"phone"`
		CreditLimit string `prn:"4" json:"creditlimit"`
		Birthday    string `prn:"5" json:"birthday"`
	}

	expectedJSON := `[
	{"name":"Johnson, John","address":"Voorstraat 32","postcode":"3122gg","phone":"020 3849381","creditlimit":"1000000","birthday":"19870101"},
	{"name":"Anderson, Paul","address":"Dorpsplein 3A","postcode":"4532 AA","phone":"030 3458986","creditlimit":"10909300","birthday":"19651203"},
	{"name":"Wicket, Steve","address":"Mendelssohnstraat 54d","postcode":"3423 ba","phone":"0313-398475","creditlimit":"93400","birthday":"19640603"},
	{"name":"Benetar, Pat","address":"Driehoog 3zwart","postcode":"2340 CC","phone":"06-28938945","creditlimit":"54","birthday":"19640904"},
	{"name":"Gibson, Mal","address":"Vredenburg 21","postcode":"3209 DD","phone":"06-48958986","creditlimit":"5450","birthday":"19781109"},
	{"name":"Friendly, User","address":"Sint Jansstraat 32","postcode":"4220 EE","phone":"0885-291029","creditlimit":"6360","birthday":"19800810"},
	{"name":"Smith, John","address":"Børkestraße 32","postcode":"87823","phone":"+44 728 889838","creditlimit":"989830","birthday":"1999092"}]`

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
	if err != nil {
		t.Log("Error: Failed to parse Workbook ->", err)
		return
	}
	assert.Nil(t, err)
	assert.NotNil(t, prnParsedItems)

	js, err := json.Marshal(prnParsedItems)
	if err != nil {
		t.Error("Error:", err)
	}
	assert.JSONEq(t, expectedJSON, string(js))
}
