package tracker_test

import (
	"github.com/shldhll/hourglass/data"
	"github.com/shldhll/hourglass/tracker"

	"time"
	"errors"
	"reflect"
	"testing"
	"strings"
	"fmt"
)

var (
	stubTime = time.Date(1970, 01, 01, 0, 0, 0, 0, time.UTC)
	stubDBWriteErr = errors.New("Error occurred while writing data")
)

const (
	stubName         = "App Name"
	stubCooldownTime = 1 * time.Second
	stubMinUsageTime = 1 * time.Nanosecond
	stubDuration     = 1 * time.Hour
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

type stubCfg struct {
	numLoops              int
	shouldLoop            bool
	getCooldownTimeCalled int
	getMinUsageTimeCalled int
	loopCheckCalled       int
	loopNextCalled        int
	cooldownTime          time.Duration
	minUsageTime          time.Duration
}

func (s *stubCfg) GetCooldownTime() time.Duration {
	s.getCooldownTimeCalled++
	return s.cooldownTime
}

func (s *stubCfg) GetMinUsageTime() time.Duration {
	s.getMinUsageTimeCalled++
	return s.minUsageTime
}

func (s *stubCfg) LoopCheck() bool {
	s.loopCheckCalled++
	return s.shouldLoop
}

func (s *stubCfg) LoopNext() {
	s.loopNextCalled++
	if s.numLoops == s.loopCheckCalled {
		s.shouldLoop = false
	}
}

func TestPing(t *testing.T) {
	system := stubOS{
		applicationName: stubName,
	}
	got := tracker.Ping(&system)
	want := tracker.NewTask(stubName, stubTime)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestCreateID(t *testing.T) {
	got := tracker.CreateID(stubName, stubTime)
	want := fmt.Sprintf(tracker.EntryIDStringFormat, stubTime.Format(tracker.EntryIDDateFormat), strings.ReplaceAll(stubName, tracker.EntryIDNameReplaceOld, tracker.EntryIDNameReplaceNew))

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNewTask(t *testing.T) {
	task := tracker.NewTask(stubName, stubTime)

	if task.AppName() != stubName {
		t.Errorf("got %v, want %v", task.AppName(), stubName)
	}

	if task.Time() != stubTime {
		t.Errorf("got %v, want %v", task.Time(), stubTime)
	}
}

func TestTaskAppName(t *testing.T) {
	task := tracker.NewTask(stubName, stubTime)
	got := task.AppName()
	want := stubName

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestTaskTime(t *testing.T) {
	task := tracker.NewTask(stubName, stubTime)
	got := task.Time()
	want := stubTime

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestStart(t *testing.T) {
	t.Run("OS functions called", func(t *testing.T) {
		system := stubOS{
			applicationName: stubName,
			realTime:        true,
		}
		db := stubDB{
			showErrorOK: 0,
		}
		config := stubCfg{
			shouldLoop:   true,
			numLoops:     1,
			cooldownTime: stubCooldownTime,
			minUsageTime: stubMinUsageTime,
		}

		tracker.Start(&system, &db, &config)
		if system.getAppNameCalled == 0 {
			t.Error("GetApplicationName() not called")
		}
		if system.nowCalled == 0 {
			t.Error("Now() not called")
		}
	})

	t.Run("DB write called", func(t *testing.T) {
		system := stubOS{
			applicationName: stubName,
			realTime:        true,
		}
		db := stubDB{
			showErrorOK: 0,
		}
		config := stubCfg{
			shouldLoop:   true,
			numLoops:     1,
			cooldownTime: stubCooldownTime,
			minUsageTime: stubMinUsageTime,
		}

		tracker.Start(&system, &db, &config)

		if db.write == 0 {
			t.Error("DB not called enough times")
		}
	})

	t.Run("Config functions called", func(t *testing.T) {
		system := stubOS{
			applicationName: stubName,
			realTime:        true,
		}
		db := stubDB{
			showErrorOK: 0,
		}
		config := stubCfg{
			shouldLoop:   true,
			numLoops:     1,
			cooldownTime: stubCooldownTime,
			minUsageTime: stubMinUsageTime,
		}

		tracker.Start(&system, &db, &config)

		if config.getCooldownTimeCalled == 0 {
			t.Error("GetCooldownTime() not called")
		}

		if config.getMinUsageTimeCalled == 0 {
			t.Error("MinUsageTime() not called")
		}

		if config.loopCheckCalled == 0 {
			t.Error("LoopCheck() not called")
		}

		if config.loopNextCalled == 0 {
			t.Error("LoopNext() not called")
		}
	})

	t.Run("Database write timeout", func(t *testing.T) {
		system := stubOS{
			applicationName: stubName,
			realTime:        true,
			shouldLog:       1,
			logChan:         make(chan string, 10),
		}
		db := stubDB{}
		config := stubCfg{
			shouldLoop:   true,
			numLoops:     1,
			cooldownTime: stubMinUsageTime,
			minUsageTime: stubMinUsageTime,
		}

		tracker.Start(&system, &db, &config)

		select {
		case msg := <-system.logChan:
			if msg != tracker.DBCallNoReturn {
				t.Errorf("got %v, want %v", msg, tracker.DBCallNoReturn)
			}
		case <-time.After(1 * time.Second):
			t.Errorf("timed out")
		}
	})

	t.Run("Database write error", func(t *testing.T) {
		system := stubOS{
			applicationName: stubName,
			realTime:        true,
			shouldLog:       1,
			logChan:         make(chan string, 10),
		}
		db := stubDB{showErrorOK: 1}
		config := stubCfg{
			shouldLoop:   true,
			numLoops:     1,
			cooldownTime: stubCooldownTime,
			minUsageTime: stubMinUsageTime,
		}

		tracker.Start(&system, &db, &config)

		select {
		case msg := <-system.logChan:
			if msg != stubDBWriteErr.Error() {
				t.Errorf("got %q, want %q", msg, tracker.DBCallNoReturn)
			}
		case <-time.After(1 * time.Second):
			t.Errorf("timed out")
		}
	})
}
