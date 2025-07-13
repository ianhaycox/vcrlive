// Package events handle the shared memory
package events

import (
	"time"

	"golang.org/x/sys/windows"
)

var (
	eventHandle windows.Handle
	err         error
)

func OpenEvent(eventName string) error {
	utf16, err := windows.UTF16PtrFromString(eventName)
	if err != nil {
		return err
	}

	eventHandle, err = windows.OpenEvent(windows.READ_CONTROL, false, utf16)

	return err
}

func WaitForSingleObject(timeout time.Duration) bool {
	t0 := time.Now().UnixNano()
	timeoutMilli := uint32(timeout / time.Millisecond)
	r, err := windows.WaitForSingleObject(eventHandle, timeoutMilli)
	if err != nil {
		remainingTimeout := timeoutMilli - uint32((time.Now().UnixNano()-t0)/1000000)
		if remainingTimeout > 0 {
			time.Sleep(time.Duration(remainingTimeout) * time.Millisecond)
		}
		return false
	}
	return r == 0
}
