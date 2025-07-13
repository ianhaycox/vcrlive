//go:generate mockgen -package irsdk -destination defines_mock.go -source defines.go
package irsdk

import "io"

const dataValidEventName string = "Local\\IRSDKDataValidEvent"
const fileMapName string = "Local\\IRSDKMemMapFileName"
const connTimeout = 30

const (
	stConnected int = 1
)

type reader interface {
	io.ReaderAt
	io.Closer
}
