package data

import (
	"bytes"
	"encoding/gob"

	"github.com/dgraph-io/badger"
)

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

// Close closes connection to database
func (b BadgerDB) Close() error {
	return b.db.Close()
}

// GetKey returns key of the entry
func (b BadgerDB) GetKey(entry Entry) string {
	return entry.ID
}

// GetBadgerDB returns a reference to BadgerDB struct
func GetBadgerDB(location string) (*BadgerDB, error) {
	db, err := badger.Open(badger.DefaultOptions(location))
	badgerDB := &BadgerDB{
		db:	db,
	}
	return badgerDB, err
}

// BadgerDBUtils represents functions required for calling DB functions
type BadgerDBUtils interface {
	Encode(Entry) ([]byte, error)
	Decode([]byte) (Entry, error)
	EncodeList([]Entry) ([]byte, error)
	DecodeList([]byte) ([]Entry, error)
}

// BadgerDBUtilsDefault represents default implementation of BadgerDBUtils
type BadgerDBUtilsDefault struct{}

// Encode returns encoded value of the given entry
func (b BadgerDBUtilsDefault) Encode(entry Entry) ([]byte, error) {
	var buff bytes.Buffer
	e := gob.NewEncoder(&buff)
	err := e.Encode(entry)
	return buff.Bytes(), err
}

// Decode returns entry after decoding the given value
func (b BadgerDBUtilsDefault) Decode(value []byte) (Entry, error) {
	return nil, nil
}

// EncodeList returns encoded value of the given entry list
func (b BadgerDBUtilsDefault) EncodeList(entryList []Entry) ([]byte, error) {
	return nil, nil
}

// DecodeList returns entry after decoding the given value
func (b BadgerDBUtilsDefault) DecodeList(value []byte) ([]Entry, error) {
	return nil, nil
}