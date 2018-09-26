# Akamai CLI for Adaptive Acceleration
*NOTE:* This tool is intended to be installed via the Akamai CLI package manager, which can be retrieved from the releases page of the [Akamai CLI](https://github.com/akamai/cli) tool.

### Local Install, if you choose not to use the akamai package manager
If you want to compile it from source, you will need Go 1.9 or later, and the [Glide](https://glide.sh) package manager installed:
1. Fetch the package:
   `go get https://github.com/partamonov/akamai-cli-reporting`
1. Change to the package directory:
   `cd $GOPATH/src/github.com/partamonov/akamai-cli-reporting`
1. Install dependencies using Glide:
   `glide install`
1. Compile the binary:
   `go build -ldflags="-s -w -X main.version=X.X.X" -o akamai-reporting`

### Credentials
In order to use this configuration, you need to:
* Set up your credential files as described in the [authorization](https://developer.akamai.com/introduction/Prov_Creds.html) and [credentials](https://developer.akamai.com/introduction/Conf_Client.html) sections of the getting started guide on developer.akamai.com.

Expects `default` section in .edgerc, can be changed via --section parameter

```
[default]
client_secret = XXXXXXXXXXXX
host = XXXXXXXXXXXX
access_token = XXXXXXXXXXXX
client_token = XXXXXXXXXXXX
```

## Overview
Generate Akamai reports

## Main Command Usage
```shell
NAME:
   akamai-reporting - A CLI to generate Akamai reports

USAGE:
   akamai-reporting [global options] command [command options] [arguments...]

VERSION:
   0.0.1

AUTHORS:
   Petr Artamonov
   Rafal Pieniazek

COMMANDS:
     generate, g  Generate a report.
     help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config FILE, -c FILE   Location of the credentials FILE (default: "/Users/partamonov/.edgerc") [$AKAMAI_EDGERC]
   --debug                  Debug info
   --no-color               Disable color output
   --section NAME, -s NAME  NAME of section to use from credentials file (default: "default") [$AKAMAI_EDGERC_SECTION]
   --help, -h               show help
   --version, -v            print the version
```

### Report command

You can get report for CP Code

```shell
NAME:
   akamai-reporting generate - Generate a report.

USAGE:
   akamai-reporting generate [command options] [arguments...]

OPTIONS:
   --cp-code value   Set of CP Codes to generate report for. Comma separated list if many (default: "all")
   --interval value  The duration of each data record, either FIVE_MINUTES, HOUR, DAY (default: "FIVE_MINUTES")
   --start value     Specifies the start of the reported period as an ISO–8601 date with timezone.
   --end value       Specifies the end of the reported period as an ISO–8601 date with timezone. Any data that matches the end value’s timestamp is excluded from the report.
   --type value      Default report type (default: "todaytraffic-by-time")
   --metrics value   Default metrics for report (default: "bytesOffload,hitsOffload,bytesOffloadAvg,hitsOffloadAvg,edgeHitsPerSecond,edgeBitsPerSecond,originBitsPerSecond,originHitsPerSecond,originBytesTotal,originHitsTotal,edgeBytesTotal,edgeHitsTotal")

> akamai reporting generate --cp-code XXXXXX --interval HOUR

{
    "metadata": {
        "name": "todaytraffic-by-time",
        "version": "1",
        "groupBy": [
            "startdatetime"
        ],
        "interval": "HOUR",
        "start": "2018-09-25T14:00:00+02:00",
        "end": "2018-09-26T14:00:00+02:00",
        "availableDataEnds": "2018-09-26T13:40:00+02:00",
        "suggestedRetryTime": "2018-09-26T14:15:24.316+02:00",
        "rowCount": 24,
        "filters": [],
        "columns": [
            {
                "name": "groupBy",
                "label": "startdatetime"
            },
            {
                "name": "bytesOffload",
                "label": "Offloaded Bytes"
            },
            {
                "name": "edgeBitsPerSecond",
                "label": "Edge Bits/Sec"
            },
            {
                "name": "edgeHitsPerSecond",
                "label": "Edge Hits/Sec"
            },
            {
                "name": "hitsOffload",
                "label": "Offloaded Hits"
            },
            {
                "name": "originBitsPerSecond",
                "label": "Origin Bits/Sec"
            },
            {
                "name": "originHitsPerSecond",
                "label": "Origin Hits/Sec"
            }
        ],
        "objectType": "cpcode",
        "objectIds": [
            "XXXXXX"
        ]
    },
    "data": [
        {
            "startdatetime": "2018-09-25T14:00:00+02:00",
            "bytesOffload": "92.06264019706161",
            "edgeBitsPerSecond": "1932708.124444444",
            "edgeHitsPerSecond": "2.9",
            "hitsOffload": "79.02298850574713",
            "originBitsPerSecond": "153405.9977777778",
            "originHitsPerSecond": "0.6083333333333333"
        },
        .....
        {
            "startdatetime": "2018-09-26T12:00:00+02:00",
            "bytesOffload": "81.76398889560950",
            "edgeBitsPerSecond": "2773834.568888889",
            "edgeHitsPerSecond": "3.066666666666667",
            "hitsOffload": "82.60869565217391",
            "originBitsPerSecond": "505836.78",
            "originHitsPerSecond": "0.5333333333333333"
        }
    ],
    "summaryStatistics": {
        "bytesOffloadAvg": {
            "value": "81.91146941726757",
            "details": {}
        },
        "edgeBytesTotal": {
            "value": "15329986742",
            "details": {}
        },
        "edgeHitsTotal": {
            "value": "159900",
            "details": {}
        },
        "hitsOffloadAvg": {
            "value": "49.11393986598426",
            "details": {}
        },
        "originBytesTotal": {
            "value": "2467536323",
            "details": {}
        },
        "originHitsTotal": {
            "value": "56670",
            "details": {}
        }
    }
}
```

#### Reset command


```shell
> akamai reporting reset PROPERTY_ID
...
```
