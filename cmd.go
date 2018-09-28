package main

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	common "github.com/apiheat/akamai-cli-common"
	edgegrid "github.com/apiheat/go-edgegrid"
	logs "github.com/sirupsen/logrus"

	"github.com/urfave/cli"
)

func cmdReport(c *cli.Context) error {
	if c.String("start") == "" && c.String("end") == "" {
		logs.Info("[" + strings.Title(appName) + "]::Report will be generated for last 24h")
	}
	if c.String("cp-code") == "all" {
		return reportAllCpCodes(c)
	}
	return reportCpCode(c)
}

func strToStrArr(str string) (strArr []string) {
	for _, s := range strings.Split(str, ",") {
		strArr = append(strArr, s)
	}
	return strArr
}

func reportAllCpCodes(c *cli.Context) error {
	logs.Info("[" + strings.Title(appName) + "]::Report will be generated for all CP Codes")
	requestOptions := edgegrid.AkamaiReportOptions{
		TypeOfReport: c.String("type"),
		Interval:     c.String("interval"),
		DateRange:    prepareTimeframes(c),
	}

	body := edgegrid.AkamaiReportingBodyAll{
		ObjectType: "cpcode",
		ObjectIds:  "all",
		Metrics:    common.StringToStringsArr(c.String("metrics")),
	}

	report, _ := apiClient.ReportingAPI.GenerateReport(body, requestOptions)

	fmt.Println(report)

	return nil
}

func reportCpCode(c *cli.Context) error {
	requestOptions := edgegrid.AkamaiReportOptions{
		TypeOfReport: c.String("type"),
		Interval:     c.String("interval"),
		DateRange:    prepareTimeframes(c),
	}

	body := edgegrid.AkamaiReportingBody{
		ObjectType: "cpcode",
		ObjectIds:  common.StringToStringsArr(c.String("cp-code")),
		Metrics:    common.StringToStringsArr(c.String("metrics")),
	}

	report, err := apiClient.ReportingAPI.GenerateReport(body, requestOptions)
	common.ErrorCheck(err)

	common.PrintJSON(report.Body)

	return nil
}

// prepareTimeframes return start and end time for our query against OpsGenie
func prepareTimeframes(c *cli.Context) string {
	var (
		start, end time.Time
		r          time.Duration
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
	}

	return fmt.Sprintf("start=%s&end=%s", url.QueryEscape(start.Format(time.RFC3339)), url.QueryEscape(end.Format(time.RFC3339)))

}

func specificDate(t time.Time, hourAt, minAt, secAt int) time.Time {
	ams, _ := time.LoadLocation("Europe/Amsterdam")
	year, month, day := t.Date()
	return time.Date(year, month, day, hourAt, minAt, secAt, 0, ams)
}
