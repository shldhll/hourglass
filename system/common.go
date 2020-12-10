package system

import "time"

// OS represents an operating system
type OS interface {
	GetApplicationName() string
	Now() time.Time
	Log(string)
}