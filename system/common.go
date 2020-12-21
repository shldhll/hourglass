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

// Cfg implements Config interface
type Cfg struct {
	cooldownTime  time.Duration
	minUsageTime  time.Duration
	loopCheckBool bool
}

// GetCooldownTime returns cooldown duration
func (c Cfg) GetCooldownTime() time.Duration {
	return c.cooldownTime
}

// GetMinUsageTime returns minimum time duration
func (c Cfg) GetMinUsageTime() time.Duration {
	return c.minUsageTime
}

// LoopCheck replicates a custom loop condition check
func (c Cfg) LoopCheck() bool {
	return c.loopCheckBool
}

// LoopNext is called before each loop
func (c Cfg) LoopNext() {
	time.Sleep(c.GetCooldownTime())
}
