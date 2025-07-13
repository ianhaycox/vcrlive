package irsdk

import (
	"testing"

	gomock "go.uber.org/mock/gomock"
)

func TestInit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mmapFile := NewMockreader(ctrl)
	mmapFile.EXPECT().ReadAt(gomock.Any(), int64(0))
	mmapFile.EXPECT().Close()

	sdk := NewIrSDK(mmapFile)
	sdk.Close()
}
