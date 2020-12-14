package system

import (
	"os/exec"
	"strings"
	"time"
)

const (
	windowIDSplitSep   = "_NET_ACTIVE_WINDOW(WINDOW): window id # "
	windowNameSplitSep = "WM_NAME(UTF8_STRING) = "
)

// Current represents the current operating system
type Current struct{}

// GetApplicationName returns the name of current forground application
func (c Current) GetApplicationName() (appName string) {
	windowIDCmd, err := exec.Command("xprop", "-root", "_NET_ACTIVE_WINDOW").Output()
	if err != nil {
		return
	}

	windowIDCmdSplitRes := strings.Split(string(windowIDCmd), windowIDSplitSep)
	if len(windowIDCmdSplitRes) <= 1 {
		return
	}

	windowID := windowIDCmdSplitRes[1]
	windowNameCmd, err := exec.Command("xprop", "-id", windowID, "WM_NAME").Output()
	if err != nil {
		return
	}

	windowNameCmdSplitRes := strings.Split(string(windowNameCmd), windowNameSplitSep)
	if len(windowNameCmdSplitRes) <= 1 {
		return
	}

	windowName := windowNameCmdSplitRes[1]
	windowNameSplitRes := strings.Split(strings.ReplaceAll(windowName, "\"", ""), " - ")
	appName = windowNameSplitRes[len(windowNameSplitRes)-1]

	return appName
}

// Now returns current time
func (c Current) Now() time.Time {
	return time.Now()
}

// Log is used for system specific logging
func (c Current) Log(msg string) {}
