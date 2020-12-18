package data

import (
	"bytes"
	"encoding/gob"

	"github.com/dgraph-io/badger"
)

// BadgerDB represents a Badger database
type BadgerDB struct {
	db      *badger.DB
	dbUtils BadgerDBUtils
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
func GetBadgerDB(location string, dbUtils BadgerDBUtils) (*BadgerDB, error) {
	var utils BadgerDBUtils = BadgerDBUtilsDefault{}
	if dbUtils != nil {
		utils = dbUtils
	}
	db, err := badger.Open(badger.DefaultOptions(location))
	badgerDB := &BadgerDB{
		db:      db,
		dbUtils: utils,
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
	var entry Entry
	d := gob.NewDecoder(bytes.NewReader(value))
	err := d.Decode(&entry)
	return entry, err
}

// EncodeList returns encoded value of the given entry list
func (b BadgerDBUtilsDefault) EncodeList(entryList []Entry) ([]byte, error) {
	var buff bytes.Buffer
	idList := []string{}

	for _, entry := range entryList {
		idList = append(idList, entry.ID)
	}
	e := gob.NewEncoder(&buff)
	err := e.Encode(idList)
	return buff.Bytes(), err
}

// DecodeList returns entry after decoding the given value
func (b BadgerDBUtilsDefault) DecodeList(value []byte) ([]Entry, error) {
	var idList []string
	var entryList []Entry
	d := gob.NewDecoder(bytes.NewReader(value))
	err := d.Decode(&idList)
	for _, id := range idList {
		entryList = append(entryList, Entry{ID: id})
	}
	return entryList, err
}