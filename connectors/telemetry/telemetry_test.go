package telemetry

import (
	"context"
	"testing"
	"time"

	"github.com/ianhaycox/vcrlive/connectors/vcrstandings"
	"github.com/ianhaycox/vcrlive/irsdk"
	"github.com/ianhaycox/vcrlive/irsdk/iryaml"
	"github.com/ianhaycox/vcrlive/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestTelemetry(t *testing.T) {
	t.Run("Happy path should get two samples then exit with a CoolDown message", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()

		sdk := irsdk.NewMockSDK(ctrl)
		gomock.InOrder(
			sdk.EXPECT().WaitForData(time.Duration(10000000)),
			sdk.EXPECT().GetLastVersion().Return(1),
			sdk.EXPECT().GetSession().Return(iryaml.IRSession{
				WeekendInfo: iryaml.WeekendInfo{TrackID: 1},
				SessionInfo: iryaml.SessionInfo{Sessions: []iryaml.Session{{SessionNum: 1}}},
				DriverInfo: iryaml.DriverInfo{Drivers: []iryaml.Driver{
					{CarIdx: 1, UserName: "1", UserID: 1},
					{CarIdx: 2, UserName: "2", UserID: 2},
				}},
			}),
			sdk.EXPECT().GetVarValue("SessionNum").Return(1, nil),
			sdk.EXPECT().GetVarValue("SessionState").Return(4, nil),
			sdk.EXPECT().GetVarValues("CarIdxClassPosition").Return([]int{0, 12, 13}, nil),
			sdk.EXPECT().GetVarValues("CarIdxLap").Return([]int{0, 67, 68}, nil),

			// Second loop
			sdk.EXPECT().WaitForData(time.Duration(10000000)),
			sdk.EXPECT().GetLastVersion().Return(1),                 // not ticked over
			sdk.EXPECT().GetVarValue("SessionState").Return(6, nil), // cool down
		)

		vcr := vcrstandings.NewMockVcrStandingsAPI(ctrl)
		vcr.EXPECT().Post(ctx, &model.LivePositions{
			Weekend: model.Weekend{TrackID: 1},
			Session: model.Session{SessionNum: 1, SessionState: "Racing"},
			Drivers: []model.Driver{
				{CarIdx: 1, UserName: "1", UserID: 1, ClassPosition: 12, Lap: 67},
				{CarIdx: 2, UserName: "2", UserID: 2, ClassPosition: 13, Lap: 68},
			},
		})

		vcr.EXPECT().Post(ctx, &model.LivePositions{
			Session: model.Session{SessionNum: 1, SessionState: "Cool Down"},
		})

		tm := NewTelemetry(sdk, vcr)

		err := tm.Run(ctx, 10, 1)
		assert.NoError(t, err)
	})
}
