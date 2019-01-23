package main

import (
	"fmt"
	"strings"
	"time"

	common "github.com/apiheat/akamai-cli-common"
	edgegrid "github.com/apiheat/go-edgegrid"
	log "github.com/sirupsen/logrus"

	"github.com/urfave/cli"
)

func cmdReport(c *cli.Context) error {
	if c.String("start") == "" && c.String("end") == "" {
		log.Info("[" + strings.Title(appName) + "]::Report will be generated for last 24h")
	}
	if c.String("cp-code") == "all" {
		return reportAllCpCodes(c)
	}
	return reportCpCode(c)
}

func reportAllCpCodes(c *cli.Context) error {
	log.Info("[" + strings.Title(appName) + "]::Report will be generated for all CP Codes")
	start, end := prepareTimeframes(c)
	requestOptions := edgegrid.ReportOptions{
		TypeOfReport: c.String("type"),
		Interval:     c.String("interval"),
		Start:        start,
		End:          end,
	}

	body := edgegrid.ReportingBodyAll{
		ObjectType: "cpcode",
		ObjectIds:  "all",
		Metrics:    common.StringToStringsArr(c.String("metrics")),
	}

	report, _ := apiClient.Reporting.GenerateReport(body, requestOptions)

	fmt.Println(report)

	return nil
}

func reportCpCode(c *cli.Context) error {
	start, end := prepareTimeframes(c)
	requestOptions := edgegrid.ReportOptions{
		TypeOfReport: c.String("type"),
		Interval:     c.String("interval"),
		Start:        start,
		End:          end,
	}

	body := edgegrid.ReportingBody{
		ObjectType: "cpcode",
		ObjectIds:  common.StringToStringsArr(c.String("cp-code")),
		Metrics:    common.StringToStringsArr(c.String("metrics")),
	}

	report, err := apiClient.Reporting.GenerateReport(body, requestOptions)
	common.ErrorCheck(err)

	common.PrintJSON(report.Body)

	return nil
}

// prepareTimeframes return start and end time for our query against OpsGenie
func prepareTimeframes(c *cli.Context) (start, end time.Time) {
	var (
		r time.Duration
	)

	switch c.String("interval") {
	case "FIVE_MINUTES":
		r = 5 * time.Minute
	case "HOUR":
		r = 1 * time.Hour
	}

	if c.String("start") == "" && c.String("end") == "" {
		start = time.Now().AddDate(0, 0, -1).Round(r)
		end = time.Now().Round(r)
		if c.String("interval") == "DAY" {
			start = specificDate(start)
			end = specificDate(end)
		}
	}

	return start, end

}

func specificDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
