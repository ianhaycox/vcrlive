package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ianhaycox/vcrlive/connectors/api"
	"github.com/ianhaycox/vcrlive/connectors/telemetry"
	"github.com/ianhaycox/vcrlive/connectors/vcrstandings"
	"github.com/ianhaycox/vcrlive/irsdk"
)

const (
	defaultWaitMilliseconds = 100
	defaultRefreshSeconds   = 10
)

var (
	progName         = filepath.Base(os.Args[0])
	ibtFile          string
	waitMilliseconds int
	refreshSeconds   int
	redact           bool
)

func main() {
	flag.StringVar(&ibtFile, "file", "", "Test data, e.g. race.bin")
	flag.IntVar(&waitMilliseconds, "wait", defaultWaitMilliseconds, "Delay in milliseconds to wait for iRacing data")
	flag.IntVar(&refreshSeconds, "refresh", defaultRefreshSeconds, "Refresh positions every n seconds")
	flag.BoolVar(&redact, "redact", false, "Obfuscate driver names for testing")
	flag.Usage = usage
	flag.Parse()

	client := vcrstandings.NewVcrStandingsService(nil, nil)

	args := flag.Args()
	if len(args) > 0 {
		client = vcrstandings.NewVcrStandingsService(api.NewAPIClient(api.NewConfiguration(args[0])), nil)
	}

	var sdk *irsdk.IRSDK

	if ibtFile == "" {
		sdk = irsdk.NewIrSDK(nil)
	} else {
		reader, err := os.Open(ibtFile) //nolint:gosec // for testing
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Init irSDK Linux(other)")

		sdk = irsdk.NewIrSDK(reader)
	}

	defer sdk.Close()

	telemetry := telemetry.NewTelemetry(sdk, client, redact)
	ctx := context.Background()

	// Keep sending telemetry data until the simulator session ends
	err := telemetry.Run(ctx, waitMilliseconds, refreshSeconds)
	if err != nil {
		log.Println(err)
	}
}

func usage() {
	w := flag.CommandLine.Output()
	_, _ = fmt.Fprintf(w, "Usage of %s: [flags] [url]\n", progName)

	flag.PrintDefaults()

	os.Exit(0)
}
