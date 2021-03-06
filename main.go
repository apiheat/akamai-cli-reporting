package main

import (
	"fmt"
	"os"
	"sort"

	common "github.com/apiheat/akamai-cli-common"
	edgegrid "github.com/apiheat/go-edgegrid"
	log "github.com/sirupsen/logrus"

	"github.com/urfave/cli"
)

var (
	apiClient       *edgegrid.Client
	apiClientOpts   *edgegrid.ClientOptions
	appVer, appName string
)

func main() {
	app := common.CreateNewApp(appName, "A CLI to generate Akamai reports", appVer)
	app.Flags = common.CreateFlags()
	app.Before = func(c *cli.Context) error {
		var err error

		apiClientOpts := &edgegrid.ClientOptions{}
		apiClientOpts.ConfigPath = c.GlobalString("config")
		apiClientOpts.ConfigSection = c.GlobalString("section")
		apiClientOpts.DebugLevel = c.GlobalString("debug")
		apiClientOpts.AccountSwitchKey = c.GlobalString("ask")

		apiClient, err = edgegrid.NewClient(nil, apiClientOpts)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "Generate a report.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "cp-code",
					Value: "all",
					Usage: "Set of CP Codes to generate report for. Comma separated list if many",
				},
				cli.StringFlag{
					Name:  "interval",
					Value: "FIVE_MINUTES",
					Usage: "The duration of each data record, either FIVE_MINUTES, HOUR, DAY",
				},
				cli.StringFlag{
					Name:  "start",
					Value: "",
					Usage: "Specifies the start of the reported period as an ISO–8601 date with timezone.",
				},
				cli.StringFlag{
					Name:  "end",
					Value: "",
					Usage: "Specifies the end of the reported period as an ISO–8601 date with timezone. Any data that matches the end value’s timestamp is excluded from the report.",
				},
				cli.StringFlag{
					Name:  "type",
					Value: "todaytraffic-by-time",
					Usage: "Default report type",
				},
				cli.StringFlag{
					Name:  "metrics",
					Value: "bytesOffload,hitsOffload,bytesOffloadAvg,hitsOffloadAvg,edgeHitsPerSecond,edgeBitsPerSecond,originBitsPerSecond,originHitsPerSecond,originBytesTotal,originHitsTotal,edgeBytesTotal,edgeHitsTotal",
					Usage: "Default metrics for report",
				},
			},
			Action: cmdReport,
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
