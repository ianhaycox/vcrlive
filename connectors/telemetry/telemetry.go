// Package telemetry gets iRacing telemetry data from the simulator
package telemetry

import (
	"context"
	"fmt"
	"log"
	"maps"
	"slices"
	"sort"
	"time"

	"github.com/ianhaycox/vcrlive/connectors/vcrstandings"
	"github.com/ianhaycox/vcrlive/irsdk"
	"github.com/ianhaycox/vcrlive/irsdk/iryaml"
	"github.com/ianhaycox/vcrlive/model"
)

type Telemetry struct {
	sdk     irsdk.SDK
	service vcrstandings.VcrStandingsAPI
	redact  bool
}

func NewTelemetry(sdk irsdk.SDK, service vcrstandings.VcrStandingsAPI, redact bool) *Telemetry {
	return &Telemetry{
		sdk:     sdk,
		service: service,
		redact:  redact,
	}
}

func (t *Telemetry) Run(ctx context.Context, waitMilliseconds int, refreshSeconds int) error {
	var (
		irSession iryaml.IRSession
		session   model.Session
		weekend   model.Weekend
		drivers   model.Drivers
	)

	latestTick := -1

	for {
		t.sdk.WaitForData(time.Duration(waitMilliseconds) * time.Millisecond)

		tick := t.sdk.GetLastVersion()
		if tick != latestTick {
			latestTick = tick
			irSession = t.sdk.GetSession()

			sessionNum, err := t.sdk.GetVarValue("SessionNum")
			if err != nil {
				session.SetState(model.Invalid)
				session.ErrorText = fmt.Sprintf("Can not determine SessionNum, err:%v, bailing...", err)

				break
			}

			weekend = model.NewWeekend(&irSession.WeekendInfo)
			session = model.NewSession(sessionNum.(int), irSession.SessionInfo.Sessions)
			drivers = model.NewDrivers(irSession.DriverInfo.Drivers, t.redact)
		}

		state, err := t.sdk.GetVarValue("SessionState")
		if err != nil {
			session.SetState(model.Invalid)
			session.ErrorText = fmt.Sprintf("Can not determine SessionState, err:%v, bailing...", err)

			break
		}

		session.SetState(state.(int))

		if state.(int) == model.Invalid {
			log.Printf("State invalid at tick:%d, ignored", tick)
			continue
		}

		if state.(int) == model.CoolDown {
			break
		}

		positions, err := t.sdk.GetVarValues("CarIdxClassPosition")
		if err != nil {
			session.SetState(model.Invalid)
			session.ErrorText = fmt.Sprintf("Can not determine CarIdxClassPosition, err:%v, bailing...", err)

			break
		}

		drivers.SetPositions(positions.([]int))

		laps, err := t.sdk.GetVarValues("CarIdxLapCompleted")
		if err != nil {
			session.SetState(model.Invalid)
			session.ErrorText = fmt.Sprintf("Can not determine CarIdxLapCompleted, err:%v, bailing...", err)

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

		err = t.service.Post(ctx, &livePositions)
		if err != nil {
			session.SetState(model.Invalid)
			session.ErrorText = fmt.Sprintf("Can not POST to endpoint, err:%v, bailing...", err)
		}

		time.Sleep(time.Duration(refreshSeconds) * time.Second)
	}

	livePositions := model.LivePositions{
		Session: session,
	}

	// POSTs either CoolDown or error
	err := t.service.Post(ctx, &livePositions)
	if err != nil {
		return fmt.Errorf("can not final post, err:%s", err)
	}

	return err
}
