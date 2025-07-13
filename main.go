package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/ianhaycox/vcrlive/irsdk"
	"github.com/ianhaycox/vcrlive/irsdk/iryaml"
	"github.com/ianhaycox/vcrlive/model"
)

const (
	defaultWaitMilliseconds = 100
	defaultRefreshSeconds   = 3
)

var (
	progName         = filepath.Base(os.Args[0])
	ibtFile          string
	waitMilliseconds int
	refreshSeconds   int
)

func main() {
	flag.StringVar(&ibtFile, "file", "", "Test data, e.g. race.bin")
	flag.IntVar(&waitMilliseconds, "wait", defaultWaitMilliseconds, "Delay in milliseconds to wait for iRacing data")
	flag.IntVar(&refreshSeconds, "refresh", defaultRefreshSeconds, "Refresh positions every n seconds")
	flag.Usage = usage
	flag.Parse()

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

	latestTick := -1

	for {
		var (
			irSession iryaml.IRSession
			session   model.Session
			weekend   model.Weekend
			drivers   model.Drivers
		)

		sdk.WaitForData(time.Duration(waitMilliseconds) * time.Millisecond)

		tick := sdk.GetLastVersion()
		if tick != latestTick {
			irSession = sdk.GetSession()

			sessionNum, err := sdk.GetVarValue("SessionNum")
			if err != nil {
				log.Printf("SessionNum:%v", err)
			}

			weekend = model.NewWeekend(&irSession.WeekendInfo)
			session = model.NewSession(sessionNum.(int), irSession.SessionInfo.Sessions)
			drivers = model.NewDrivers(irSession.DriverInfo.Drivers)
		}

		state, err := sdk.GetVarValue("SessionState")
		if err != nil {
			log.Printf("SessionState:%v", err)
		}

		session.SetState(state.(int))

		positions, err := sdk.GetVarValues("CarIdxClassPosition")
		if err != nil {
			log.Printf("CarIdxClassPosition:%v", err)
		}

		drivers.SetPositions(positions.([]int))

		livePositions := model.LivePositions{
			Weekend: weekend,
			Session: session,
			Drivers: slices.Collect(maps.Values(drivers)),
		}

		fmt.Println(dump(livePositions))

		time.Sleep(time.Duration(refreshSeconds) * time.Second)
	}
}

func dump(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")

	return string(b)
}

func usage() {
	w := flag.CommandLine.Output()
	_, _ = fmt.Fprintf(w, "Usage of %s: [flags] filename \n", progName)

	flag.PrintDefaults()
}
