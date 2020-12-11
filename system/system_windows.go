package system

import "time"

// Current represents the current operating system
type Current struct{}

// GetApplicationName returns the name of current forground application
func (c Current) GetApplicationName() (appName string) {
	return appName
}

// Now returns current time
func (c Current) Now() time.Time {
	return time.Now()
}

// Log is used for system specific logging
func (c Current) Log(msg string) {}
