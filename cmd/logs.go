/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"time"

	"github.com/shldhll/hourglass/data"
	"github.com/shldhll/hourglass/tracker"
	"github.com/spf13/cobra"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Logs of current day",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := exec.Command("mkdir", "-p", os.Getenv("HOME")+"/.hourglass").Output()
		if err != nil {
			log.Fatal(err)
			return
		}
		db, err := data.GetBadgerDB(os.Getenv("HOME")+"/.hourglass/data", data.BadgerDBUtilsDefault{})
		if err != nil {
			log.Fatal(err)
			return
		}
		today := time.Now()
		entries, err := db.ReadList(today.Format(tracker.EntryIDDateFormat))
		if err != nil {
			log.Fatal(err)
			return
		}
		sort.Slice(entries[:], func(i, j int) bool {
			return entries[i].Duration > entries[j].Duration
		})
		for i, entry := range entries {
			fmt.Println(i+1, ")")
			fmt.Print("Name: \t\t", entry.AppName)
			duration := entry.Duration.Round(time.Second)
			h := duration / time.Hour
			duration -= h * time.Hour
			m := duration / time.Minute
			duration -= m * time.Minute
			s := duration / time.Second
			fmt.Println("Duration:\t", fmt.Sprintf("%02d:%02d:%02d", h, m, s))
			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
}
