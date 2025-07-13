// Package events fake windows events to ease Linux development
package events

import "time"

func OpenEvent(eventName string) error {
	return nil
}

func WaitForSingleObject(timeout time.Duration) bool {
	return true
}
