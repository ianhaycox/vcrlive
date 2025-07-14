//go:generate mockgen -package irsdk -destination irsdk_mock.go -source irsdk.go

// Package irsdk iRacing SDK
package irsdk

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hidez8891/shm"
	"github.com/ianhaycox/vcrlive/irsdk/iryaml"
	"github.com/ianhaycox/vcrlive/win/events"
	"gopkg.in/yaml.v3"
)

const fileMapSize int32 = 1164 * 1024

type SDK interface {
	RefreshSession()
	WaitForData(timeout time.Duration) bool
	GetVars() ([]Variable, error)
	GetVar(name string) (Variable, error)
	GetVarValue(name string) (interface{}, error)
	GetVarValues(name string) (interface{}, error)
	GetSession() iryaml.IRSession
	GetLastVersion() int
	IsConnected() bool
	GetYaml() string
	Close()
}

// IRSDK is the main SDK object clients must use
type IRSDK struct {
	SDK
	r             reader
	h             *header
	session       iryaml.IRSession
	s             []string
	tVars         *TelemetryVars
	lastValidData int64
}

// NewIrSDK creates a SDK instance to operate with
func NewIrSDK(r reader) *IRSDK {
	if r == nil {
		var err error

		r, err = shm.Open(fileMapName, fileMapSize)
		if err != nil {
			log.Fatalf("shared memory error. err:%s", err)
		}
	}

	sdk := &IRSDK{r: r, lastValidData: 0}

	err := events.OpenEvent(dataValidEventName)
	if err != nil {
		log.Fatal("Open event", err)
	}

	initIRSDK(sdk)

	return sdk
}

func (sdk *IRSDK) RefreshSession() {
	if sessionStatusOK(sdk.h.status) {
		sRaw := readSessionData(sdk.r, sdk.h)

		err := yaml.Unmarshal([]byte(sRaw), &sdk.session)
		if err != nil {
			log.Println(err)
		}

		sdk.s = strings.Split(sRaw, "\n")
	}
}

func (sdk *IRSDK) WaitForData(timeout time.Duration) bool {
	if !sdk.IsConnected() {
		initIRSDK(sdk)
	}

	if events.WaitForSingleObject(timeout) {
		sdk.RefreshSession()
		return readVariableValues(sdk)
	}

	return false
}

func (sdk *IRSDK) GetVars() ([]Variable, error) {
	if !sessionStatusOK(sdk.h.status) {
		return make([]Variable, 0), fmt.Errorf("session is not active")
	}

	results := make([]Variable, len(sdk.tVars.vars))

	sdk.tVars.mux.Lock()
	defer sdk.tVars.mux.Unlock()

	i := 0

	for _, variable := range sdk.tVars.vars {
		results[i] = variable
		i++
	}

	return results, nil
}

func (sdk *IRSDK) GetVar(name string) (Variable, error) {
	if !sessionStatusOK(sdk.h.status) {
		return Variable{}, fmt.Errorf("session is not active")
	}

	sdk.tVars.mux.Lock()
	defer sdk.tVars.mux.Unlock()

	if v, ok := sdk.tVars.vars[name]; ok {
		return v, nil
	}

	return Variable{}, fmt.Errorf("telemetry variable %q not found", name)
}

func (sdk *IRSDK) GetVarValue(name string) (interface{}, error) {
	var (
		r   Variable
		err error
	)

	if r, err = sdk.GetVar(name); err == nil {
		return r.Value, nil
	}

	return r, err
}

func (sdk *IRSDK) GetVarValues(name string) (interface{}, error) {
	var (
		r   Variable
		err error
	)

	if r, err = sdk.GetVar(name); err == nil {
		return r.Values, nil
	}

	return r, err
}

func (sdk *IRSDK) GetSession() iryaml.IRSession {
	return sdk.session
}

func (sdk *IRSDK) GetLastVersion() int {
	if !sessionStatusOK(sdk.h.status) {
		return -1
	}

	sdk.tVars.mux.Lock()
	defer sdk.tVars.mux.Unlock()

	last := sdk.tVars.lastVersion

	return last
}

func (sdk *IRSDK) GetSessionData(path string) (string, error) {
	if !sessionStatusOK(sdk.h.status) {
		return "", fmt.Errorf("session not connected")
	}

	return getSessionDataPath(sdk.s, path)
}

func (sdk *IRSDK) IsConnected() bool {
	if sdk.h != nil {
		if sessionStatusOK(sdk.h.status) && (sdk.lastValidData+connTimeout > time.Now().Unix()) {
			return true
		}
	}

	return false
}

func (sdk *IRSDK) GetYaml() string {
	return strings.Join(sdk.s, "\n")
}

// Close clean up sdk resources
func (sdk *IRSDK) Close() {
	_ = sdk.r.Close()
}

func initIRSDK(sdk *IRSDK) {
	h := readHeader(sdk.r)
	sdk.h = &h
	sdk.s = nil

	if sdk.tVars != nil {
		sdk.tVars.vars = nil
	}

	if sessionStatusOK(h.status) {
		sRaw := readSessionData(sdk.r, &h)

		err := yaml.Unmarshal([]byte(sRaw), &sdk.session)
		if err != nil {
			log.Println(err)
		}

		sdk.s = strings.Split(sRaw, "\n")
		sdk.tVars = readVariableHeaders(sdk.r, &h)
		readVariableValues(sdk)
	}
}

func sessionStatusOK(status int) bool {
	return (status & stConnected) > 0
}
