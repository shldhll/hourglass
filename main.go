package main

import (
	"log"
	"time"

	"github.com/shldhll/hourglass/data"
	"github.com/shldhll/hourglass/system"
	"github.com/shldhll/hourglass/tracker"
)

const (
	minUsageTime = 1 * time.Second
	cooldownTime = 1 * time.Second
)

func main() {
	os := system.Current{}
	cfg := system.GetConfig(cooldownTime, minUsageTime)
	dbUtils := data.BadgerDBUtilsDefault{}
	db, err := data.GetBadgerDB("./data/store", dbUtils)
	if err != nil {
		log.Fatal("database error")
	}
	tracker.Start(os, db, cfg)
}
