package tracker

import (
	"github.com/shldhll/hourglass/data"
	"github.com/shldhll/hourglass/system"

	"fmt"
	"strings"
	"time"
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

	// DBCallNoReturn is used when call to database times out
	DBCallNoReturn = "Call to DB did not return"
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

// Start is the entrypoint function
func Start(o system.OS, db data.DB, cfg system.Config) {
	var prevApp string
	prevTime := o.Now()
	cooldownTime := cfg.GetCooldownTime()
	minUsageTime := cfg.GetMinUsageTime()
	entryDict := make(map[string]data.Entry)

	for cfg.LoopCheck() {
		task := Ping(o)
		currApp := task.AppName()
		currTime := task.Time()
		errChan := make(chan error)

		if diff := currTime.Sub(prevTime); diff >= minUsageTime && currApp == prevApp {
			go func(database data.DB, appName string, startTime, endTime time.Time, errorChan chan<- error) {
				duration := endTime.Sub(startTime)
				id := CreateID(appName, startTime)
				entry := data.Entry{
					ID:       id,
					AppName:  appName,
					Duration: duration,
				}

				err := db.Write(entry)
				prevTime = currTime
				if _, ok := entryDict[entry.ID]; !ok && err == nil {
					err = db.WriteList(entry)
				}
				errorChan <- err
			}(db, currApp, prevTime, currTime, errChan)
		}

		if prevApp != currApp {
			prevApp = currApp
			prevApp = currApp
		}

		select {
		case err := <-errChan:
			if err != nil {
				o.Log(err.Error())
			}
		case <-time.After(cooldownTime):
			cfg.LoopNext()
		}

		cfg.LoopNext()
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
