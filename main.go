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
	"sort"
	"time"

	"github.com/ianhaycox/vcrlive/irsdk"
	"github.com/ianhaycox/vcrlive/irsdk/iryaml"
	"github.com/ianhaycox/vcrlive/model"
)

const (
	defaultWaitMilliseconds = 100
	defaultRefreshSeconds   = 10
	url                     = ""
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

	var (
		irSession iryaml.IRSession
		session   model.Session
		weekend   model.Weekend
		drivers   model.Drivers
	)

	for {
		sdk.WaitForData(time.Duration(waitMilliseconds) * time.Millisecond)

		tick := sdk.GetLastVersion()
		if tick != latestTick {
			latestTick = tick
			irSession = sdk.GetSession()

			sessionNum, err := sdk.GetVarValue("SessionNum")
			if err != nil {
				session.SetState(model.Invalid)
				session.ErrorText = fmt.Sprintf("Can not determine SessionNum, err:%v, bailing...", err)

				break
			}

			weekend = model.NewWeekend(&irSession.WeekendInfo)
			session = model.NewSession(sessionNum.(int), irSession.SessionInfo.Sessions)
			drivers = model.NewDrivers(irSession.DriverInfo.Drivers)
		}

		state, err := sdk.GetVarValue("SessionState")
		if err != nil {
			session.SetState(model.Invalid)
			session.ErrorText = fmt.Sprintf("Can not determine SessionState, err:%v, bailing...", err)

			break
		}

		session.SetState(state.(int))

		if state.(int) == model.CoolDown {
			break
		}

		positions, err := sdk.GetVarValues("CarIdxClassPosition")
		if err != nil {
			session.SetState(model.Invalid)
			session.ErrorText = fmt.Sprintf("Can not determine CarIdxClassPosition, err:%v, bailing...", err)

			break
		}

		drivers.SetPositions(positions.([]int))

		laps, err := sdk.GetVarValues("CarIdxLap")
		if err != nil {
			session.SetState(model.Invalid)
			session.ErrorText = fmt.Sprintf("Can not determine CarIdxLap, err:%v, bailing...", err)

			break
		}

		drivers.SetLaps(laps.([]int))

		sortedDrivers := slices.Collect(maps.Values(drivers))
		sort.Slice(sortedDrivers, func(i, j int) bool { return sortedDrivers[i].CarIdx < sortedDrivers[j].CarIdx })

		livePositions := model.LivePositions{
			Weekend: weekend,
			Session: session,
			Drivers: sortedDrivers,
		}

		err = post(&livePositions)
		if err != nil {
			session.SetState(model.Invalid)
			session.ErrorText = fmt.Sprintf("Can not POST to endpoint:%s, err:%v, bailing...", url, err)
		}

		time.Sleep(time.Duration(refreshSeconds) * time.Second)
	}

	livePositions := model.LivePositions{
		Session: session,
	}

	err := post(&livePositions)
	if err != nil {
		log.Printf("can not final post to %s, err:%s", url, err)
	}
}

func post(livePositions *model.LivePositions) error { //nolint:unparam // TODO
	fmt.Println(dump(livePositions))

	return nil
}

func dump(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")

	return string(b)
}

func usage() {
	w := flag.CommandLine.Output()
	_, _ = fmt.Fprintf(w, "Usage of %s: [flags]\n", progName)

	flag.PrintDefaults()
}
