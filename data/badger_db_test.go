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

func assertError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("No error expected, got %v", err)
	}
}

func clean() error {
	return os.RemoveAll(dbLocation)
}
