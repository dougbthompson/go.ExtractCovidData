package main

// CaseTableStruct structure
type CaseTableStruct struct {
	Date           string
	ConfirmedCases float64
	DaysSince      float64
}

// TimeSeriesStruct structure
type TimeSeriesStruct struct {
	Date                             string
	NewConfirmedCasesSevenDayAverage float64
	NewDeathsSevenDayAverage         float64
	TotalPatients                    float64
}

// County structure
type County struct {
	Date                  string
	Fips                  string
	CountyName            string
	NewConfirmedCases     float64
	ConfirmedCases        float64
	ConfirmedCasesPer100K float64
	NewDeaths             float64
	Deaths                float64
	DeathsPer100K         float64
	DoublingTimeCases     float64
	CaseTable             []CaseTableStruct
	TimeSeries            []TimeSeriesStruct
	SQLCounty             []string
	SQLCaseTable          []string
	SQLTimeSeries         []string
}

var (
	cty           *County
	strHead       []string
	strCaseTable  []string
	strTimeSeries []string
)

// time.Sleep(100 * time.Millisecond)
