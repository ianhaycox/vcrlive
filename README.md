# iRacing Live Telemetry

Utility to expose [iRacing](https://www.iracing.com/) live race position data.

## Running

With iRacing running,

      vcrlive.exe https://example.com/

will sample the telemetry data every 10 seconds and POST the current car positions to the specified URL

Alternatively, omit the URL and it will output the payload to the console.

For testing you can use a Race replay and jump backwards and forwards between Practice, Qualifying and the Race.

## Typical payloads sent to the endpoint

```
{
  "weekend": {
    "track_id": 168,
    "track_display_name": "Suzuka International Racing Course",
    "track_config_name": "Grand Prix",
    "series_id": 285,
    "season_id": 5582,
    "sub_session_id": 78289018,
    "official": 1,
    "race_week": 3,
    "event_type": "Race",
    "category": "SportsCar",
    "num_car_classes": 2,
    "num_car_types": 3
  },
  "session": {
    "session_num": 1,
    "session_laps": "2",
    "session_type": "Race",
    "session_name": "RACE",
    "session_state": "Racing",
    "error_text": ""
  },
  "drivers": [
    {
      "car_idx": 1,
      "user_name": "Test driver",
      "user_id": 996799,
      "car_class_id": 84,
      "car_id": 77,
      "class_position": 3,
      "laps_completed": 3,
      "irating": 6176,
      "car_number_raw": 2
    }
  ]
}
```

`session_state` can be one of,
- Invalid
- Get In Car
- Warmup
- Parade Laps
- Racing
- Checkered
- Cool Down

`session_type` can be one of,
- PRACTICE
- QUALIFY
- RACE

at the end of the Race, the final payload will be,

```
{
  "weekend": {},
  "session": {
    "session_num": 2,
    "session_laps": "unlimited",
    "session_type": "Race",
    "session_name": "RACE",
    "session_state": "Cool Down"
  }
}
```

[An abbreviated race example with JSON payloads for Practice, Qualifying and Race](./example.json.txt)

## iRacing SDK

The [iRacing SDK](https://forums.iracing.com/categories/iracing-api-s-and-development-discussions) documentation to access the shared Windows memory
to retrieve telemetry. An iRacing account is required to view.

## Usage

Requires [go v24.1 or greater](https://go.dev/doc/install) to be installed.

From a `cmd` or `PowerShell` prompt build a Windows executable,

`go build -o vcrlive.exe  main.go`

Start iRacing

Then run via,

`vcrlive.exe http://my-site.com/`

## Options

```
go run main.go --help
Usage of main: [flags] [url]
  -file string
    	Test data, e.g. race.bin
  -redact
    	Obfuscate driver names for testing
  -refresh int
    	Refresh positions every n seconds (default 10)
  -wait int
    	Delay in milliseconds to wait for iRacing data (default 100)
```

## Development

`go run main.go -file testdata.bin`

See [pyirsdk](https://github.com/kutu/pyirsdk/blob/master/tutorials/02%20Using%20irsdk%20script.md) for creating `.bin` telemetry files.

