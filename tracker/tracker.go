package tracker

import "time"

// Task struct represents a running application.
type Task struct {
	applicationName string
	recordedTime    time.Time
}

// Time returns the capture Time
func (t Task) Time() time.Time {
	return t.recordedTime
}

// AppName returns the title of current application
func (t Task) AppName() string {
	return t.applicationName
}