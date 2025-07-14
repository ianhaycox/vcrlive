# iRacing Live Telemetry

Utility to expose [iRacing](https://www.iracing.com/) live race position data.

## iRacing SDK

The [iRacing SDK](https://forums.iracing.com/categories/iracing-api-s-and-development-discussions) to access the shared Windows memory to retrieve telemetry.

An iRacing account is required.

## Typical payload sent to the endpoint



## Usage

```
$ go run main.go --help
Usage of main: [flags]
  -file string
        Test data, e.g. race.bin
  -refresh int
        Refresh positions every n seconds (default 10)
  -wait int
        Delay in milliseconds to wait for iRacing data (default 100)
```