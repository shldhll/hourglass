package tracker_test

import (
	"github.com/shldhll/hourglass/data"

	"time"
	"errors"
)

var (
	stubTime = time.Date(1970, 01, 01, 0, 0, 0, 0, time.UTC)
	stubDBWriteErr = errors.New("Error occurred while writing data")
)

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

type stubDB struct {
	showErrorOK int
	write       int
	read        int
	writeList   int
	readList    int
}

func (s *stubDB) Write(entry data.Entry) error {
	s.write++
	if s.showErrorOK == 0 {
		return nil
	}

	return stubDBWriteErr
}

func (s *stubDB) Read(key string) (data.Entry, error) {
	s.read++
	if s.showErrorOK != 0 {
		return data.Entry{}, stubDBWriteErr
	}

	return data.Entry{}, nil
}

func (s *stubDB) WriteList(entry data.Entry) error {
	s.writeList++
	if s.showErrorOK != 0 {
		return stubDBWriteErr
	}

	return nil
}

func (s *stubDB) ReadList(date string) ([]data.Entry, error) {
	s.readList++
	if s.showErrorOK != 0 {
		return []data.Entry{}, stubDBWriteErr
	}

	return []data.Entry{}, nil
}
