package system

import "time"

// OS represents an operating system
type OS interface {
	GetApplicationName() string
	Now() time.Time
	Log(string)
}

// Config represents various configurations
type Config interface {
	GetCooldownTime() time.Duration
	GetMinUsageTime() time.Duration
	LoopCheck() bool
	LoopNext()
}
