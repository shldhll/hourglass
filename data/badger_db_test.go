package data_test

import (
	"github.com/shldhll/hourglass/data"
	"github.com/shldhll/hourglass/tracker"

	"testing"
	"os"
	"reflect"
	"time"
)

const (
	dbLocation      = "./db_test_dir"
	stubName        = "App Name"
	stubDuration    = 1 * time.Hour
	multiWriteCount = 3
)

var stubTime = time.Date(1970, 01, 01, 0, 0, 0, 0, time.UTC)

func TestGetBadgerDB(t *testing.T) {
	defer clean()
	db, err := data.GetBadgerDB(dbLocation, nil)
	assertError(t, err)
	defer db.Close()
}

func TestBadgerDBClose(t *testing.T) {
	defer clean()
	db, err := data.GetBadgerDB(dbLocation, nil)
	assertErrorFatal(t, err)

	err = db.Close()
	assertError(t, err)
}

func TestBadgerDBWrite(t *testing.T) {
	t.Run("Simple write test", func(t *testing.T) {
		defer clean()
		db, err := data.GetBadgerDB(dbLocation, nil)
		assertErrorFatal(t, err)
		defer db.Close()

		entry := createEntry()
		err = db.Write(entry)
		assertErrorFatal(t, err)

		got, err := db.Read(db.GetKey(entry))
		assertErrorFatal(t, err)

		want := entry
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("Multi write test", func(t *testing.T) {
		defer clean()
		db, err := data.GetBadgerDB(dbLocation, nil)
		assertErrorFatal(t, err)
		defer db.Close()

		entry := createEntry()

		for i := 0; i < multiWriteCount; i++ {
			err = db.Write(entry)
			assertErrorFatal(t, err)
		}

		got, err := db.Read(db.GetKey(entry))
		assertErrorFatal(t, err)

		entry.Duration = stubDuration * multiWriteCount
		want := entry
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("Empty key write test", func(t *testing.T) {
		defer clean()
		db, err := data.GetBadgerDB(dbLocation, nil)
		assertErrorFatal(t, err)
		defer db.Close()

		entry := data.Entry{}

		err = db.Write(entry)
		if err != badger.ErrEmptyKey {
			assertError(t, err)
		}
	})
}

func assertError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("No error expected, got %v", err)
	}
}

func assertErrorFatal(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("No error expected, got %v", err)
	}
}

func clean() error {
	return os.RemoveAll(dbLocation)
}

func createEntry() data.Entry {
	id := tracker.CreateID(stubName, stubTime)
	return data.Entry{id, stubName, stubDuration}
}