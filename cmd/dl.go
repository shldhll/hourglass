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
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"

	"github.com/shldhll/hourglass/data"
	"github.com/shldhll/hourglass/tracker"
	"github.com/spf13/cobra"
)

const idStringSeparator = "_"
const htmlCode = `<!DOCTYPE html>
<html>
    <head>
        <title>Tracking Data</title>
        <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.1.3/css/bootstrap.min.css" integrity="sha384-MCw98/SFnGE8fJT3GXwEOngsV7Zt27NXFoaoApmYm81iuXoPkFOJwJ8ERdknLPMO" crossorigin="anonymous">
    </head>
    <body>
        {{with .Entries}}
        <table class="table">
            <thead>
              <tr>
                <th scope="col">Date</th>
                <th scope="col">Application</th>
                <th scope="col">Duration</th>
              </tr>
            </thead>
            <tbody>
              {{range .}}
              <tr>
                <th scope="row">{{.Date}}</th>
                <td>{{.AppName}}</td>
                <td>{{.Duration}}</td>
              </tr>
              {{end}}
            </tbody>
          </table>
          {{end}}
    </body>
</html>`

type task struct {
	Date     string
	AppName  string
	Duration string
}

// dlCmd represents the dl command
var dlCmd = &cobra.Command{
	Use:   "dl",
	Short: "Download the tracking data",
	Run: func(cmd *cobra.Command, args []string) {
		var startTime time.Time
		if len(args) < 2 {
			println("usage: hourglass dl [today|week|month] <filename.html>")
			return
		}

		if args[0] == "today" {
			startTime = time.Now()
		} else if args[0] == "week" {
			startTime = time.Now().AddDate(0, -6, 0)
		} else if args[0] == "month" {
			startTime = time.Now().AddDate(0, -30, 0)
		}

		entries := dl(startTime)
		ok := make([]task, 0)
		for _, e := range entries {
			var t task
			t.Date = strings.Split(e.ID, "_")[0]
			t.AppName = e.AppName
			duration := e.Duration.Round(time.Second)
			h := duration / time.Hour
			duration -= h * time.Hour
			m := duration / time.Minute
			duration -= m * time.Minute
			s := duration / time.Second
			t.Duration = fmt.Sprintf("%02d:%02d:%02d", h, m, s)
			ok = append(ok, t)
		}
		t, err := template.New("data").Parse(htmlCode)
		if err != nil {
			println("error occured:", err)
			return
		}

		var result bytes.Buffer

		err = t.Execute(&result, struct{ Entries []task }{Entries: ok})
		if err != nil {
			return
		}

		fileName := args[1]
		if !strings.Contains(fileName, ".html") {
			fileName += ".html"
		}

		f, err := os.Create(fileName)
		if err != nil {
			println("error occured:", err)
		}

		defer f.Close()

		_, err = f.WriteString(result.String())
		if err != nil {
			println("error occured:", err)
		}

		fmt.Println("Data saved to", fileName)
	},
}

func dl(start time.Time) []data.Entry {
	entries := make([]data.Entry, 0)
	_, err := exec.Command("mkdir", "-p", os.Getenv("HOME")+"/.hourglass").Output()
	if err != nil {
		log.Fatal(err)
		return entries
	}
	db, err := data.GetBadgerDB(os.Getenv("HOME")+"/.hourglass/data", data.BadgerDBUtilsDefault{})
	if err != nil {
		log.Print("db error:", err)
		return entries
	}
	for start.Format(tracker.EntryIDDateFormat) != time.Now().AddDate(0, 1, 0).Format(tracker.EntryIDDateFormat) {
		ok, err := db.ReadList(start.Format(tracker.EntryIDDateFormat))
		start = start.AddDate(0, 1, 0)
		if err != nil {
			println("error while downloading:", err)
			continue
		}
		entries = append(entries, ok...)
	}

	return entries
}

func init() {
	rootCmd.AddCommand(dlCmd)
}
