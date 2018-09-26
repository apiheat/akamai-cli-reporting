package main

import (
	"os"
	"sort"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/fatih/color"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

var (
	colorOn, raw, debug       bool
	version, appName          string
	configSection, configFile string
	edgeConfig                edgegrid.Config
)

// Constants
const (
	URL     = "/reporting-api/v1/reports"
	padding = 3
)

func main() {
	_, inCLI := os.LookupEnv("AKAMAI_CLI")

	appName = "akamai-reporting"
	if inCLI {
		appName = "akamai reporting"
	}

	app := cli.NewApp()
	app.Name = appName
	app.HelpName = appName
	app.Usage = "A CLI to generate Akamai reports"
	app.Version = version
	app.Copyright = ""
	app.Authors = []cli.Author{
		{
			Name: "Petr Artamonov",
		},
		{
			Name: "Rafal Pieniazek",
		},
	}

	dir, _ := homedir.Dir()
	dir += string(os.PathSeparator) + ".edgerc"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "section, s",
			Value:       "default",
			Usage:       "`NAME` of section to use from credentials file",
			Destination: &configSection,
			EnvVar:      "AKAMAI_EDGERC_SECTION",
		},
		cli.StringFlag{
			Name:        "config, c",
			Value:       dir,
			Usage:       "Location of the credentials `FILE`",
			Destination: &configFile,
			EnvVar:      "AKAMAI_EDGERC",
		},
		cli.BoolFlag{
			Name:        "no-color",
			Usage:       "Disable color output",
			Destination: &colorOn,
		},
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "Debug info",
			Destination: &debug,
		},
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

	app.Before = func(c *cli.Context) error {
		if c.Bool("no-color") {
			color.NoColor = true
		}

		config(configFile, configSection)
		return nil
	}

	app.Run(os.Args)
}
