package tracker

import (
	"github.com/shldhll/hourglass/system"

	"time"
)

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

// NewTask creates a new task
func NewTask(appName string, recordedTime time.Time) *Task {
	return &Task{
		applicationName: appName,
		recordedTime:    recordedTime,
	}
}

// Ping returns window information in the form of Task struct.
func Ping(o system.OS) *Task {
	return NewTask(o.GetApplicationName(), o.Now())
}