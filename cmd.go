package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"strings"
	"time"

	client "github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"github.com/urfave/cli"
)

type ReportBody struct {
	ObjectType string   `json:"objectType"`
	ObjectIds  []string `json:"objectIds"`
	Metrics    []string `json:"metrics"`
}

type ReportBodyAll struct {
	ObjectType string   `json:"objectType"`
	ObjectIds  string   `json:"objectIds"`
	Metrics    []string `json:"metrics"`
}

func cmdReport(c *cli.Context) error {
	if debug {
		if c.String("start") == "" && c.String("end") == "" {
			fmt.Println("# Will generate report for last 24h")
		}
		if c.String("cp-code") == "all" {
			fmt.Println("# Will generate report for all CP Codes")
		}
	}
	return reportProperty(c)
}

func strToStrArr(str string) (strArr []string) {
	for _, s := range strings.Split(str, ",") {
		strArr = append(strArr, s)
	}
	return strArr
}

func reportProperty(c *cli.Context) error {
	dateRange := prepareTimeframes(c)

	urlStr := fmt.Sprintf("%s/%s/versions/1/report-data?%s&interval=%s", URL, c.String("type"), dateRange, c.String("interval"))

	if debug {
		println(urlStr)
	}

	bodyStr, err := json.Marshal(ReportBodyAll{
		ObjectType: "cpcode",
		ObjectIds:  "all",
		Metrics:    strToStrArr(c.String("metrics")),
	})
	errorCheck(err)

	if c.String("cp-code") != "all" {
		bodyStr, err = json.Marshal(ReportBody{
			ObjectType: "cpcode",
			ObjectIds:  strToStrArr(c.String("cp-code")),
			Metrics:    strToStrArr(c.String("metrics")),
		})
		errorCheck(err)
	}

	if debug {
		println(string(bodyStr))
	}

	data, respCode := fetchData(urlStr, "POST", bytes.NewBufferString(string(bodyStr)))

	if debug {
		fmt.Printf("Response code: %d\n", respCode)
	}

	printJSON(data)

	return nil
}

func errorCheck(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func printJSON(str string) {
	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, []byte(str), "", "    ")
	if error != nil {
		log.Println("JSON parse error: ", error)
		return
	}
	fmt.Println(string(prettyJSON.Bytes()))
	return
}

func fetchData(urlPath, method string, body io.Reader) (string, int) {
	req, err := client.NewRequest(edgeConfig, method, urlPath, body)
	errorCheck(err)

	resp, err := client.Do(edgeConfig, req)
	errorCheck(err)

	defer resp.Body.Close()
	byt, _ := ioutil.ReadAll(resp.Body)

	return string(byt), resp.StatusCode
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
