// Copyright 2015 Ondřej Doněk. All rights reserved.
// See LICENSE file for more details about licensing.

// odTimeTracker is simple time-tracking tool.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"github.com/odTimeTracker/odtimetracker-go-lib"
	"github.com/odTimeTracker/odtimetracker-go-lib/database"
	"os"
	"os/user"
	"path"
	"log"
	"net/http"
	"strings"
)

var (
	appName = "odTimeTracker"
	appShortName = "odtimetracker"
	appVersion = odtimetracker.Version{ Major: 0, Minor: 1, Maintenance: 0, }
	appInfo = appName + " " + appVersion.String()
	appDescription = "Simple tool for time-tracking."
	// Used template type ("bootstrap", "dojo", "polymer"):
	templateType = "bootstrap"
	// Database related:
	dbPath = getDatabasePath()
	// Errors:
	ErrTemplateDoesNotExist = errors.New("The template does not exist.")
	ErrDatabaseConnectionFailed = errors.New("Unable to connect database.")
)

// define a type for the response
type httpHandler struct{}

// let that type implement the ServeHTTP method (defined in interface http.Handler)
func (h httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)

	if (r.URL.String() == "/favicon.ico") {
		renderFavicon(w)
	} else if (r.URL.String() == "/RunningActivity.json") {
		renderRunningActivityJson(w)
	} else if (r.URL.String() == "/Activities.json") {
		renderActivitiesJson(w)
	} else if (r.URL.String() == "/Projects.json") {
		renderProjectsJson(w)
	} else {
		renderMainPage(w)
	}
}

// Init function.
func init() {
	log.Println("Entering init function...")
	log.Printf("Database path: %s\n", dbPath)
}

// Main (entry) function.
// TODO Using command-line arguments provide several UI types ('bootstrap', 'dojo', 'polymer')
func main() {
	log.Println("Entering main function...")

	if len(os.Args) > 1 {
		usage()
	}

	var tplType string
	if os.Args[0] == "--help" {
		usage()
	} else if strings.HasPrefix(os.Args[0], "--type=") == true {
		tplType = strings.Replace(os.Args[0], "--type=", "", 1)
	}

	// Ensure template type is correct
	if (tplType == "bootstrap" || tplType == "dojo" || tplType == "polymer") {
		templateType = tplType
	}

	// TODO Theme with Polymer used is not implemented yet!
	if (templateType == "polymer") {
		log.Println("TODO Theme with Polymer used is not implemented yet!")
		templateType = "bootstrap"
	}

	var h httpHandler
	http.ListenAndServe("localhost:4000", h)
}

// Get database path
func getDatabasePath() string {
	usr, err := user.Current()
	if err != nil {
		log.Println("Unable to create correct database path. Using in-memory database!")
		return ":memory:"
	}

	return path.Join(usr.HomeDir, ".odtimetracker.sqlite")
}

// Helper method for dealing with errors.
func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

// Render JSON with details about currently running activity.
func renderRunningActivityJson(w http.ResponseWriter) error {
	log.Println("Rendering RunningActivity.json...")

	db, err := database.InitStorage(dbPath)
	checkError(err)
	defer db.Close()

	activity, err := database.SelectActivityRunning(db)
	checkError(err)

	json, err := json.Marshal(activity)
	checkError(err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)

	return nil
}

// Render JSON with activities.
func renderActivitiesJson(w http.ResponseWriter) error {
	log.Println("Rendering Activities.json...")

	db, err := database.InitStorage(dbPath)
	checkError(err)
	defer db.Close()

	activities, err := database.SelectActivities(db, -1)
	checkError(err)

	json, err := json.Marshal(activities)
	checkError(err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)

	return nil
}

// Render JSON with projects.
func renderProjectsJson(w http.ResponseWriter) error {
	log.Println("Rendering Projects.json...")

	db, err := database.InitStorage(dbPath)
	checkError(err)
	defer db.Close()

	projects, err := database.SelectProjects(db, -1)
	checkError(err)

	json, err := json.Marshal(projects)
	checkError(err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)

	return nil
}

// Render favicon.
func renderFavicon(w http.ResponseWriter) error {
	log.Println("TODO Render favicon!")
	// ...
	return nil
}

// Render main HTML page.
func renderMainPage(w http.ResponseWriter) error {
	log.Println("Rendering main page...")

	data := map[string]string{"Name": "odTimeTracker",}

	p := "ui/" + templateType + "/"
	tpl, err := template.ParseFiles(p + "main.tmpl", p + "header.tmpl", p + "footer.tmpl")
	checkError(err)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tpl.Execute(w, data)
	checkError(err)

	return nil
}

// Print usage information.
func usage() {
	fmt.Println(appInfo)
	fmt.Printf("\n%s\n\n", appDescription)
	fmt.Printf("Usage:\n\n")
	fmt.Printf("%s --help         Print this help\n", appShortName)
	fmt.Printf("%s --type=[TYPE]  Use template of given type\n\n", appShortName)
	fmt.Printf("Available template types are: bootstrap,dojo,polymer\n\n")
	os.Exit(0)
}
