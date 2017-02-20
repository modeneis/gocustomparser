package gocustomparser_test

import (
	"testing"
	"gocustomparser"
	"github.com/stretchr/testify/assert"
	"encoding/json"
	"golang.org/x/text/encoding/charmap"
)

func TestParseCSV(t *testing.T) {

	type Workbook struct {
		Name        string `csv:"0" json:"name"`
		Address     string `csv:"1" json:"address"`
		Postcode    string `csv:"2" json:"postcode"`
		Phone       string `csv:"3" json:"phone"`
		CreditLimit string `csv:"4" json:"creditlimit"`
		Birthday    string `csv:"5" json:"birthday"`
	}

	expectedJSON := `[
	{"name":"Johnson, John","address":"Voorstraat 32","postcode":"3122gg","phone":"020 3849381","creditlimit":"10000","birthday":"01/01/1987"},
	{"name":"Anderson, Paul","address":"Dorpsplein 3A","postcode":"4532 AA","phone":"030 3458986","creditlimit":"109093","birthday":"03/12/1965"},
	{"name":"Wicket, Steve","address":"Mendelssohnstraat 54d","postcode":"3423 ba","phone":"0313-398475","creditlimit":"934","birthday":"03/06/1964"},
	{"name":"Benetar, Pat","address":"Driehoog 3zwart","postcode":"2340 CC","phone":"06-28938945","creditlimit":"54","birthday":"04/09/1964"},
	{"name":"Gibson, Mal","address":"Vredenburg 21","postcode":"3209 DD","phone":"06-48958986","creditlimit":"54.5","birthday":"09/11/1978"},
	{"name":"Friendly, User","address":"Sint Jansstraat 32","postcode":"4220 EE","phone":"0885-291029","creditlimit":"63.6","birthday":"10/08/1980"},
	{"name":"Smith, John","address":"Børkestraße 32","postcode":"87823","phone":"+44 728 889838","creditlimit":"9898.3","birthday":"20/09/199"}]`

	csvParser := gocustomparser.CustomParser{
		File:          "testfiles/Workbook2.csv",
		CustomDecoder: charmap.ISO8859_1.NewDecoder(),
		Separator:     ',',
		SkipFirstLine: true,
	}
	var csvParsedItems, err = csvParser.Parse(Workbook{})
	if err != nil {
		t.Log("Error: Failed to parse Workbook ->", err)
		return
	}
	assert.Nil(t, err)
	assert.NotNil(t, csvParsedItems)

	js, err := json.Marshal(csvParsedItems)
	if err != nil {
		t.Error("Error:", err)
	}
	assert.JSONEq(t, expectedJSON, string(js))
}