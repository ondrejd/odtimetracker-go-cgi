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

	// TODO Try to use `http.HandleFunc("/images", handler)` to serve all images at once!
	if (r.URL.String() == "/images/mstile-144x144.png") {
		fp := path.Join("images", "mstile-144x144.png")
		http.ServeFile(w, r, fp)
	} else if (r.URL.String() == "/images/favicon-16x16.png") {
		fp := path.Join("images", "favicon-16x16.png")
		http.ServeFile(w, r, fp)
	} else if (r.URL.String() == "/images/favicon-32x32.png") {
		fp := path.Join("images", "favicon-32x32.png")
		http.ServeFile(w, r, fp)
	} else if (r.URL.String() == "/images/favicon-96x96.png") {
		fp := path.Join("images", "favicon-96x96.png")
		http.ServeFile(w, r, fp)
	} else if (r.URL.String() == "/images/favicon-160x160.png") {
		fp := path.Join("images", "favicon-160x160.png")
		http.ServeFile(w, r, fp)
	} else if (r.URL.String() == "/images/favicon-192x192.png") {
		fp := path.Join("images", "favicon-192x192.png")
		http.ServeFile(w, r, fp)
	} else if (r.URL.String() == "/images/apple-touch-icon-57x57.png") {
		fp := path.Join("images", "apple-touch-icon-57x57.png")
		http.ServeFile(w, r, fp)
	} else if (r.URL.String() == "/images/apple-touch-icon-60x60.png") {
		fp := path.Join("images", "apple-touch-icon-60x60.png")
		http.ServeFile(w, r, fp)
	} else if (r.URL.String() == "/images/apple-touch-icon-72x72.png") {
		fp := path.Join("images", "apple-touch-icon-72x72.png")
		http.ServeFile(w, r, fp)
	} else if (r.URL.String() == "/images/apple-touch-icon-76x76.png") {
		fp := path.Join("images", "apple-touch-icon-76x76.png")
		http.ServeFile(w, r, fp)
	} else if (r.URL.String() == "/images/apple-touch-icon-114x114.png") {
		fp := path.Join("images", "apple-touch-icon-114x114.png")
		http.ServeFile(w, r, fp)
	} else if (r.URL.String() == "/images/apple-touch-icon-120x120.png") {
		fp := path.Join("images", "apple-touch-icon-120x120.png")
		http.ServeFile(w, r, fp)
	} else if (r.URL.String() == "/images/apple-touch-icon-144x144.png") {
		fp := path.Join("images", "apple-touch-icon-144x144.png")
		http.ServeFile(w, r, fp)
	} else if (r.URL.String() == "/images/apple-touch-icon-152x152.png") {
		fp := path.Join("images", "apple-touch-icon-152x152.png")
		http.ServeFile(w, r, fp)
	} else if (r.URL.String() == "/images/apple-touch-icon-180x180.png") {
		fp := path.Join("images", "apple-touch-icon-180x180.png")
		http.ServeFile(w, r, fp)
	} else if (r.URL.String() == "/browserconfig.xml") {
		log.Println("Rendering browserconfig.xml...")
		w.Header().Set("Content-Type", "application/xml;charset=utf-8")
		browserconfig := `<?xml version="1.0" encoding="utf-8"?>
<browserconfig>
  <msapplication>
    <tile>
      <square70x70logo src="images/mstile-70x70.png"/>
      <square150x150logo src="images/mstile-150x150.png"/>
      <square310x310logo src="images/mstile-310x310.png"/>
      <wide310x150logo src="images/mstile-310x150.png"/>
      <TileColor>#e7e7e7</TileColor>
    </tile>
  </msapplication>
</browserconfig>
`
		w.Write([]byte(browserconfig))
	} else if (r.URL.String() == "/GetRunningActivity") {
		getRunningActivity(w, r)
	} else if (r.URL.String() == "/StartActivity") {
		startActivity(w, r)
	} else if (r.URL.String() == "/StopActivity") {
		stopActivity(w, r)
	} else if (r.URL.String() == "/ListActivities") {
		listActivities(w, r)
	} else if (r.URL.String() == "/ListProjects") {
		listProjects(w, r)
	} else {
		mainPage(w, r)
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

// Start activity.
func startActivity(w http.ResponseWriter, r *http.Request) error {
	log.Println("TODO Start activity!")

	r.ParseForm()
	log.Println(r.Form())
	name := r.FormValue("name")
	project := r.FormValue("project")
	description := r.FormValue("description")
	tags := r.FormValue("tags")

	log.Println(name)
	log.Println(project)
	log.Println(description)
	log.Println(tags)

	// ...

	return nil
}

// Stop activity.
func stopActivity(w http.ResponseWriter, r *http.Request) error {
	log.Println("TODO Stop activity!")

	r.ParseForm()
	log.Println(r.Form())
	aid := r.FormValue("aid")

	log.Println(aid)

	// ...

	return nil
}

// Render JSON with details about currently running activity.
func getRunningActivity(w http.ResponseWriter, r *http.Request) error {
	log.Println("Rendering RunningActivity.json...")

	db, err := database.InitStorage(dbPath)
	checkError(err)
	defer db.Close()

	activity, _ := database.SelectActivityRunning(db)
	json, err := json.Marshal(activity)
	checkError(err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)

	return nil
}

// Render JSON with activities.
func listActivities(w http.ResponseWriter, r *http.Request) error {
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
func listProjects(w http.ResponseWriter, r *http.Request) error {
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

// Render main HTML page.
func mainPage(w http.ResponseWriter, r *http.Request) error {
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
