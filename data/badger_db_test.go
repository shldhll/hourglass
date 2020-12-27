package data_test

import (
	"github.com/shldhll/hourglass/data"

	"testing"
	"os"
)

const dbLocation = "./db_test_dir"

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
