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
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/shldhll/hourglass/data"
	"github.com/shldhll/hourglass/system"
	"github.com/shldhll/hourglass/tracker"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the time tracker",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := exec.Command("mkdir", "-p", os.Getenv("HOME")+"/.hourglass").Output()
		if err != nil {
			log.Fatal(err)
			return
		}
		db, err := data.GetBadgerDB(os.Getenv("HOME")+"/.hourglass/data", data.BadgerDBUtilsDefault{})
		if err != nil {
			log.Print("db error:", err)
			return
		}
		tracker.Start(system.Current{}, db, system.GetConfig(1*time.Second, 1*time.Second))
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
