
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

var (
	jsonAll  []byte
	jsonSubs []string
	jsonStr  string
)

func main() {
	var counties []*County
	var county *County

	var fileName string
	switch len(os.Args) {
	case 1:
		fileName = "ca.2020-05-19.json"
	case 2:
		fileName = os.Args[1]
	default:
		// more parameters...
	}

	jsonFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer jsonFile.Close()

	newRead := io.Reader(jsonFile)
	jsonAll, err = ioutil.ReadAll(newRead)
	if err != nil {
		fmt.Println(err)
	}

	jsonSubs = strings.SplitAfter(string(jsonAll), "COUNTIES_TIMESERIES = [")
	jsonStr = "{\"thedata\":[" + jsonSubs[1]
	jsonStr = jsonStr[0:len(jsonStr)-2] + "}"
	bytz := []byte(jsonStr)

	var f1 interface{}
	err = json.Unmarshal(bytz, &f1)
	if err != nil {
		fmt.Println(err)
	}
	mp1 := f1.(map[string]interface{})

	for k, v := range mp1 {
		switch vv := v.(type) {
		case nil:
			// fmt.Println(k, "is a nil value (outer)", vv)
		case string:
			// fmt.Println(k, "is a string (outer)", vv)
		case float64:
			// fmt.Println(k, "is a float64 (outer)", vv)
		case []interface{}:
			// fmt.Println(k, "is an outer array: (outer)") // top level "thedata" array, placeholder

			for _, v2 := range vv {
				mp2 := v2.(map[string]interface{})
				// fmt.Println(k2)

				county = new(County)
				county.Date = fileName[3:13]
				county.CaseTable = make([]CaseTableStruct, 0)
				county.TimeSeries = make([]TimeSeriesStruct, 0)
				counties = append(counties, county)

				for x1, x2 := range mp2 {
					switch x3 := x2.(type) {
					case nil:
						saveCountyFloat(county, x1, nilToZero(x3))
					case string:
						saveCountyString(county, x1, x3)
					case float64:
						saveCountyFloat(county, x1, x3)

					case []interface{}:
						for _, x5 := range x3 {
							x5m, _ := x5.(map[string]interface{})

							switch x1 {
							case "case_table": // confirmed_cases:10 date:2020-04-02 days_since_10:0
								saveCountyCaseTable(county, x5m["date"].(string), nilToZero(x5m["confirmed_cases"]), nilToZero(x5m["days_since"]))
							case "time_series": // date:2020-05-18 new_confirmed_cases_seven_day_average:0 new_deaths_seven_day_average:0 total_patients:<nil>
								saveCountyTimeSeries(county, x5m["date"].(string), nilToZero(x5m["new_confirmed_cases_seven_day_average"]),
									nilToZero(x5m["new_deaths_seven_day_average"]), nilToZero(x5m["total_patients"]))
							default:
								fmt.Println("Inner other default ...")
							}
						}
						// fmt.Println("End Inner ----------")
					default:
						fmt.Println(x1, x2, "is of a type I don't know how to handle (inner)")
					}
				}
			}
		default:
			fmt.Println(k, vv, "is of a type I don't know how to handle (outer)")
		}
	}
	fmt.Println(counties)
	genSQLCode(counties)
	mysqlProcess(counties)
}

func nilToZero(mv interface{}) float64 {
	if reflect.ValueOf(mv).IsValid() {
		return mv.(float64)
	}
	return 0.0
}

func saveCountyFloat(c *County, k string, v float64) {
	switch k {
	case "new_confirmed_cases":
		c.NewConfirmedCases = v
	case "confirmed_cases":
		c.ConfirmedCases = v
	case "confirmed_cases_per_100k":
		c.ConfirmedCasesPer100K = v
	case "new_deaths":
		c.NewDeaths = v
	case "deaths":
		c.Deaths = v
	case "deaths_per_100k":
		c.DeathsPer100K = v
	default:
	}
}

func saveCountyString(c *County, k string, v string) {
	switch k {
	case "fips":
		c.Fips = v
	case "county":
		c.CountyName = v
	default:
	}
}

func saveCountyCaseTable(c *County, Date string, ConfirmedCases float64, DaysSince float64) {
	nct := CaseTableStruct{
		Date:           Date,
		ConfirmedCases: ConfirmedCases,
		DaysSince:      DaysSince,
	}
	c.CaseTable = append(c.CaseTable, nct)
}

func saveCountyTimeSeries(c *County, Date string, NewConfirmedCasesSevenDayAverage float64, NewDeathsSevenDayAverage float64, TotalPatients float64) {
	nts := TimeSeriesStruct{
		Date:                             Date,
		NewConfirmedCasesSevenDayAverage: NewConfirmedCasesSevenDayAverage,
		NewDeathsSevenDayAverage:         NewDeathsSevenDayAverage,
		TotalPatients:                    TotalPatients,
	}
	c.TimeSeries = append(c.TimeSeries, nts)
}
