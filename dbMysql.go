package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func mysqlProcess(c []*County) {
	pool, err := sql.Open("mysql", SQLConnectString)
	if err != nil {
		panic(err.Error())
	}
	defer pool.Close()

	pool.SetConnMaxLifetime(0)
	pool.SetMaxIdleConns(3)
	pool.SetMaxOpenConns(3)

	ctx1, stop := context.WithCancel(context.Background())
	defer stop()

	// the actual database connection
	con, err := pool.Conn(ctx1)
	if err != nil {
		panic(err.Error())
	}
	defer con.Close()
	// time.Sleep(100 * time.Millisecond)

	for _, cty = range c {
		result, err := con.ExecContext(ctx1, cty.SQLCounty[0])
		if err != nil {
			log.Fatal("unable to insert into [CACounty] table: ", err)
		}
		lastID, err := result.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}

		for _, sct := range cty.SQLCaseTable {
			sct = fmt.Sprintf("%s", strings.Replace(sct, "ID", fmt.Sprintf("%d", lastID), 1))
			_, err := con.ExecContext(ctx1, sct)
			if err != nil {
				log.Fatal("unable to insert into [CACaseTable] table: ", err)
			}
		}
		for _, sts := range cty.SQLTimeSeries {
			sts = fmt.Sprintf("%s", strings.Replace(sts, "ID", fmt.Sprintf("%d", lastID), 1))
			_, err := con.ExecContext(ctx1, sts)
			if err != nil {
				log.Fatal("unable to insert into [CATimeSeries] table: ", err)
			}
		}
	}
}

func genSQLCode(c []*County) {

	HeadColumns := "CDate,Fips,CountyName,NewConfirmedCases,ConfirmedCases,ConfirmedCasesPer100K,NewDeaths,Deaths,DeathsPer100K,DoublingTimeCases"
	CaseTableColumns := "CountyId,CTDate,ConfirmedCases,DaysSince"
	TimeSeriesColumns := "CountyId,TSDate,NewConfirmedCasesSevenDayAverage,NewDeathsSevenDayAverage,TotalPatients"

	for _, cty = range c {
		cty.SQLCounty = make([]string, 0)
		cty.SQLCaseTable = make([]string, 0)
		cty.SQLTimeSeries = make([]string, 0)

		cty.SQLCounty = append(cty.SQLCounty, fmt.Sprintf("insert into CACounty (%s) select \"%s\",\"%s\",\"%s\",%f,%f,%f,%f,%f,%f,%f;\n",
			HeadColumns, cty.Date, cty.Fips, cty.CountyName, cty.NewConfirmedCases, cty.ConfirmedCases, cty.ConfirmedCasesPer100K,
			cty.NewDeaths, cty.Deaths, cty.DeathsPer100K, cty.DoublingTimeCases))

		ct := cty.CaseTable
		for _, xct := range ct {
			cty.SQLCaseTable = append(cty.SQLCaseTable, fmt.Sprintf("insert into CACaseTable (%s) select ID,\"%s\",%f,%f;\n",
				CaseTableColumns, xct.Date, xct.ConfirmedCases, xct.DaysSince))
		}
		ts := cty.TimeSeries
		for _, xts := range ts {
			cty.SQLTimeSeries = append(cty.SQLTimeSeries, fmt.Sprintf("insert into CATimeSeries (%s) select ID,\"%s\",%f,%f,%f;\n",
				TimeSeriesColumns, xts.Date, xts.NewConfirmedCasesSevenDayAverage, xts.NewDeathsSevenDayAverage, xts.TotalPatients))
		}
	}
}
