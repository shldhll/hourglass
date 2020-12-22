package tracker_test

import "time"

var stubTime = time.Date(1970, 01, 01, 0, 0, 0, 0, time.UTC)

type stubOS struct {
	applicationName  string
	realTime         bool
	getAppNameCalled int
	nowCalled        int
	shouldLog        int
	logChan          chan string
}

func (s *stubOS) GetApplicationName() string {
	s.getAppNameCalled++
	return s.applicationName
}

func (s *stubOS) Now() time.Time {
	s.nowCalled++
	if s.realTime {
		return time.Now()
	}
	return stubTime
}

func (s *stubOS) Log(msg string) {
	if s.shouldLog != 0 {
		s.logChan <- msg
	}
}