package data

import "github.com/dgraph-io/badger"

// BadgerDB represents a Badger database
type BadgerDB struct {
	db      *badger.DB
}

// Write writes given entry to database
func (b BadgerDB) Write(entry Entry) error {
	return nil
}

// Read retrives entry with given id from the database
func (b BadgerDB) Read(id string) (Entry, error) {
	var e Entry
	return e, nil
}