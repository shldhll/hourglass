package data

import (
	"time"
)

// DB represents a database
type DB interface {
	Write(entry Entry) error
	WriteList(entry Entry) error
	Read(id string) (Entry, error)
	ReadList(date string) ([]Entry, error)
}

// Entry represents a database entry
type Entry struct {
	ID       string
	AppName  string
	Duration time.Duration
}
