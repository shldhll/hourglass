package tracker

import "time"

// Task struct represents a running application.
type Task struct {
	applicationName string
	recordedTime    time.Time
}