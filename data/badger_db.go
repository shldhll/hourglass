package data

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"strings"

	"github.com/dgraph-io/badger"
)

const (
	// EntryIDDateSeparator is used for separating date from the rest
	EntryIDDateSeparator = "_"
	// ErrReadPrefixText is used as prefix text for read errors
	ErrReadPrefixText = "Following errors occurred while reading the entries:"
)

// BadgerDB represents a Badger database
type BadgerDB struct {
	db      *badger.DB
	dbUtils BadgerDBUtils
}

// Write writes given entry to database
func (b BadgerDB) Write(entry Entry) error {
	key := []byte(b.GetKey(entry))

	existingEntry, err := b.Read(string(key))
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}

	if err != badger.ErrKeyNotFound {
		entry.Duration += existingEntry.Duration
	}

	value, err := b.dbUtils.Encode(entry)
	if err != nil {
		return err
	}

	err = b.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, value)
		return err
	})

	return err
}

// Read retrives entry with given id from the database
func (b BadgerDB) Read(id string) (Entry, error) {
	var e Entry

	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(id))
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			entry, err := b.dbUtils.Decode(val)
			if err != nil {
				return err
			}
			e = entry
			return nil
		})

		return err
	})

	return e, err
}

// ReadIDList retrieves list of IDs matching the given date
func (b BadgerDB) ReadIDList(date string) ([]string, error) {
	var list []string

	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(date))
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			entryList, err := b.dbUtils.DecodeList(val)
			if err != nil {
				return err
			}
			for _, entry := range entryList {
				list = append(list, entry.ID)
			}
			return nil
		})

		return err
	})
	return list, err
}

// ReadList returns list of entries matching the given date
func (b BadgerDB) ReadList(date string) ([]Entry, error) {
	var errStr strings.Builder
	entryList := []Entry{}

	idList, err := b.ReadIDList(date)
	if err != nil {
		return entryList, err
	}

	for _, entryID := range idList {
		entry, err := b.Read(entryID)
		if err != nil {
			fmt.Fprintf(&errStr, "%q, ", err.Error())
			continue
		}

		entryList = append(entryList, entry)
	}

	if len(errStr.String()) != 0 {
		err := errors.New(fmt.Sprint(ErrReadPrefixText, errStr.String()))
		return entryList, err
	}

	return entryList, err
}

// WriteList writes given entry into matching list of entries
func (b BadgerDB) WriteList(entry Entry) error {
	var entryList []Entry = make([]Entry, 1)
	entryList[0] = entry

	key := []byte(b.GetDate(entry))

	existingEntryIDList, err := b.ReadIDList(string(key))
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}
	var existingEntryList []Entry
	for _, entryID := range existingEntryIDList {
		existingEntryList = append(existingEntryList, Entry{ID: entryID})
	}

	if err != badger.ErrKeyNotFound {
		shouldAdd := true
		for i := range existingEntryList {
			if existingEntryList[i].ID == entry.ID {
				shouldAdd = false
				break
			}
		}

		if !shouldAdd {
			return nil
		}

		entryList = append(entryList, existingEntryList...)
	}

	value, err := b.dbUtils.EncodeList(entryList)
	if err != nil {
		return err
	}
	err = b.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, value)
		return err
	})

	return err
}

// Close closes connection to database
func (b BadgerDB) Close() error {
	return b.db.Close()
}

// GetKey returns key of the entry
func (b BadgerDB) GetKey(entry Entry) string {
	return entry.ID
}

// GetDate extracts date from the given entry
func (b BadgerDB) GetDate(entry Entry) string {
	splitID := strings.Split(entry.ID, EntryIDDateSeparator)
	date := splitID[0]
	return date
}

// GetBadgerDB returns a reference to BadgerDB struct
func GetBadgerDB(location string, dbUtils BadgerDBUtils) (*BadgerDB, error) {
	var utils BadgerDBUtils = BadgerDBUtilsDefault{}
	if dbUtils != nil {
		utils = dbUtils
	}
	options := badger.DefaultOptions(location)
	options.Logger = nil
	db, err := badger.Open(options)
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
