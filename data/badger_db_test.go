package data_test

import (
	"github.com/shldhll/hourglass/data"
	"github.com/shldhll/hourglass/tracker"

	"testing"
	"os"
	"reflect"
	"time"
	"errors"
	"fmt"

	"github.com/dgraph-io/badger"
)

const (
	dbLocation      = "./db_test_dir"
	stubName        = "App Name"
	stubDuration    = 1 * time.Hour
	multiWriteCount = 3
)

var stubTime = time.Date(1970, 01, 01, 0, 0, 0, 0, time.UTC)

type stubDBUtils struct {
	errDelayCount      int
	encodeErrCount     int
	decodeErrCount     int
	encodeListErrCount int
	decodeListErrCount int
	encodeErr          error
	decodeErr          error
	encodeListErr      error
	decodeListErr      error
	showEncodeErr      bool
	showDecodeErr      bool
	showEncodeListErr  bool
	showDecodeListErr  bool
}

func (s *stubDBUtils) Encode(data.Entry) ([]byte, error) {
	if s.showEncodeErr && s.encodeErrCount >= s.errDelayCount {
		return []byte{}, s.encodeErr
	}
	s.encodeErrCount++
	return nil, nil
}

func (s *stubDBUtils) Decode([]byte) (data.Entry, error) {
	if s.showDecodeErr && s.decodeErrCount >= s.errDelayCount {
		return data.Entry{}, s.decodeErr
	}
	s.decodeErrCount++
	return data.Entry{}, nil
}

func (s *stubDBUtils) EncodeList([]data.Entry) ([]byte, error) {
	if s.showEncodeListErr && s.encodeListErrCount >= s.errDelayCount {
		return nil, s.encodeListErr
	}
	s.encodeListErrCount++
	return nil, nil
}

func (s *stubDBUtils) DecodeList([]byte) ([]data.Entry, error) {
	if s.showDecodeListErr && s.decodeListErrCount >= s.errDelayCount {
		return nil, errors.New("DecodeList error")
	}
	s.decodeListErrCount++
	return nil, nil
}

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

	t.Run("Encode error write test", func(t *testing.T) {
		defer clean()
		encodeErr := errors.New("Encode error")
		dbUtils := stubDBUtils{showEncodeErr: true, encodeErr: encodeErr}
		db, err := data.GetBadgerDB(dbLocation, &dbUtils)
		assertErrorFatal(t, err)
		defer db.Close()

		err = db.Write(data.Entry{ID: "id"})
		assertErrorEqual(t, err, encodeErr)
	})
}

func TestBadgerDBWriteList(t *testing.T) {
	t.Run("WriteList test", func(t *testing.T) {
		defer clean()
		db, err := data.GetBadgerDB(dbLocation, nil)
		assertErrorFatal(t, err)
		defer db.Close()

		num := 5
		entryList := createEntryList(num)
		want := make([]string, num)

		for i, entry := range entryList {
			err = db.WriteList(entry)
			assertErrorFatal(t, err)
			want[i] = entry.ID
		}

		for i, entry := range entryList {
			err = db.WriteList(entry)
			assertErrorFatal(t, err)
			want[i] = entry.ID
		}

		got, err := db.ReadIDList(db.GetDate(entryList[num-1]))
		assertErrorFatal(t, err)

		for _, entry1 := range got {
			ok := false
			for _, entry2 := range want {
				if entry1 == entry2 {
					ok = true
					break
				}
			}
			if !ok {
				t.Errorf("Entry ID %v of got not found in want", entry1)
			}
		}
	})

	t.Run("Adding/overwriting to existing list", func(t *testing.T) {
		defer clean()
		db, err := data.GetBadgerDB(dbLocation, nil)
		assertErrorFatal(t, err)
		defer db.Close()

		num := 5
		entryList := createEntryList(num)
		want := make([]string, num)

		err = db.WriteList(entryList[num-1])
		assertErrorFatal(t, err)

		for i, entry := range entryList {
			err = db.WriteList(entry)
			assertErrorFatal(t, err)
			want[i] = entry.ID
		}

		got, err := db.ReadIDList(db.GetDate(entryList[num-1]))
		assertErrorFatal(t, err)

		for _, entry1 := range got {
			ok := false
			for _, entry2 := range want {
				if entry1 == entry2 {
					ok = true
					break
				}
			}
			if !ok {
				t.Errorf("Entry ID %v of got not found in want", entry1)
			}
		}
	})

	t.Run("Empty key", func(t *testing.T) {
		defer clean()
		db, err := data.GetBadgerDB(dbLocation, nil)
		assertErrorFatal(t, err)
		defer db.Close()

		entry := data.Entry{}

		err = db.WriteList(entry)
		if err != badger.ErrEmptyKey {
			assertError(t, err)
		}
	})

	t.Run("EncodeList error", func(t *testing.T) {
		defer clean()
		encodeListErr := errors.New("EncodeList error")
		dbUtils := stubDBUtils{showEncodeListErr: true, encodeListErr: encodeListErr}
		db, err := data.GetBadgerDB(dbLocation, &dbUtils)
		assertErrorFatal(t, err)
		defer db.Close()

		entry := createEntry()

		err = db.WriteList(entry)
		assertErrorEqual(t, err, encodeListErr)
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

func assertErrorEqual(t *testing.T, got, want error) {
	t.Helper()
	if got.Error() != want.Error() {
		t.Errorf("want %q, got %q", want.Error(), got.Error())
	}
}

func assertNotNil(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func assertEqual(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %v", got, want)
	}
}

func clean() error {
	return os.RemoveAll(dbLocation)
}

func createEntry() data.Entry {
	id := tracker.CreateID(stubName, stubTime)
	return data.Entry{id, stubName, stubDuration}
}

func createEntryList(num int) []data.Entry {
	entryList := make([]data.Entry, num)
	for i := 0; i < num; i++ {
		id := tracker.CreateID(fmt.Sprint(i), stubTime)
		entryList[i] = data.Entry{id, stubName, stubDuration}
	}
	return entryList
}
