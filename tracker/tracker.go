package tracker

import (
	"github.com/shldhll/hourglass/system"

	"time"
	"strings"
	"fmt"
)

const (
	// EntryIDStringFormat is the format used for Sprintf function to create the ID for database entry
	EntryIDStringFormat = "%v_%s"
	// EntryIDDateFormat is the date format used in the ID
	EntryIDDateFormat = "2006-01-02"
	// EntryIDNameReplaceOld is the string which is to be replaced by EntryIDNameReplaceNew
	EntryIDNameReplaceOld = " "
	// EntryIDNameReplaceNew is the string which replaces EntryIDNameReplaceOld
	EntryIDNameReplaceNew = ""
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

// CreateID creates an ID string using formatted date and string
func CreateID(appName string, date time.Time) string {
	formattedDate := date.Format(EntryIDDateFormat)
	formattedAppName := strings.ReplaceAll(appName, EntryIDNameReplaceOld, EntryIDNameReplaceNew)

	return fmt.Sprintf(EntryIDStringFormat, formattedDate, formattedAppName)
}
