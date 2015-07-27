// Copyright 2015 Ondřej Doněk. All rights reserved.
// See LICENSE file for more details about licensing.

// odTimeTracker is simple time-tracking tool.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/odTimeTracker/odtimetracker-go-lib"
	"github.com/odTimeTracker/odtimetracker-go-lib/database"
	"github.com/odTimeTracker/odtimetracker-go-cgi/jsonrpc"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path"
	"strings"
	"time"
)

var (
	appName                     = "odTimeTracker"
	appShortName                = "odtimetracker"
	appVersion                  = odtimetracker.Version{Major: 0, Minor: 1, Maintenance: 0}
	appInfo                     = appName + " " + appVersion.String()
	appDescription              = "Simple tool for time-tracking."
	templateType                = "bootstrap"       // Used template type (just "bootstrap" at the moment):
	dbPath                      = getDatabasePath() // Path to SQLite database
	ErrTemplateDoesNotExist     = errors.New("The template does not exist.")
	ErrDatabaseConnectionFailed = errors.New("Unable to connect database.")
)

// Define a type for the response
type httpHandler struct{}

// Let that type implement the ServeHTTP method (defined in interface http.Handler)
func (h httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)

	// TODO Try to use `http.HandleFunc("/images", handler)` to serve all images at once!
	if r.URL.String() == "/images/mstile-144x144.png" {
		fp := path.Join("images", "mstile-144x144.png")
		http.ServeFile(w, r, fp)
	} else if r.URL.String() == "/images/favicon-16x16.png" {
		fp := path.Join("images", "favicon-16x16.png")
		http.ServeFile(w, r, fp)
	} else if r.URL.String() == "/images/favicon-32x32.png" {
		fp := path.Join("images", "favicon-32x32.png")
		http.ServeFile(w, r, fp)
	} else if r.URL.String() == "/images/favicon-96x96.png" {
		fp := path.Join("images", "favicon-96x96.png")
		http.ServeFile(w, r, fp)
	} else if r.URL.String() == "/images/favicon-160x160.png" {
		fp := path.Join("images", "favicon-160x160.png")
		http.ServeFile(w, r, fp)
	} else if r.URL.String() == "/images/favicon-192x192.png" {
		fp := path.Join("images", "favicon-192x192.png")
		http.ServeFile(w, r, fp)
	} else if r.URL.String() == "/images/apple-touch-icon-57x57.png" {
		fp := path.Join("images", "apple-touch-icon-57x57.png")
		http.ServeFile(w, r, fp)
	} else if r.URL.String() == "/images/apple-touch-icon-60x60.png" {
		fp := path.Join("images", "apple-touch-icon-60x60.png")
		http.ServeFile(w, r, fp)
	} else if r.URL.String() == "/images/apple-touch-icon-72x72.png" {
		fp := path.Join("images", "apple-touch-icon-72x72.png")
		http.ServeFile(w, r, fp)
	} else if r.URL.String() == "/images/apple-touch-icon-76x76.png" {
		fp := path.Join("images", "apple-touch-icon-76x76.png")
		http.ServeFile(w, r, fp)
	} else if r.URL.String() == "/images/apple-touch-icon-114x114.png" {
		fp := path.Join("images", "apple-touch-icon-114x114.png")
		http.ServeFile(w, r, fp)
	} else if r.URL.String() == "/images/apple-touch-icon-120x120.png" {
		fp := path.Join("images", "apple-touch-icon-120x120.png")
		http.ServeFile(w, r, fp)
	} else if r.URL.String() == "/images/apple-touch-icon-144x144.png" {
		fp := path.Join("images", "apple-touch-icon-144x144.png")
		http.ServeFile(w, r, fp)
	} else if r.URL.String() == "/images/apple-touch-icon-152x152.png" {
		fp := path.Join("images", "apple-touch-icon-152x152.png")
		http.ServeFile(w, r, fp)
	} else if r.URL.String() == "/images/apple-touch-icon-180x180.png" {
		fp := path.Join("images", "apple-touch-icon-180x180.png")
		http.ServeFile(w, r, fp)
	} else if r.URL.String() == "/browserconfig.xml" {
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
	} else if r.URL.String() == "/ui/bootstrap/script.js" {
		w.Header().Set("Content-Type", "text/javascript;charset=utf-8")
		js, _ := ioutil.ReadFile("ui/bootstrap/script.js")
		w.Write([]byte(js))
	} else if r.URL.String() == "/ui/bootstrap/style.css" {
		w.Header().Set("Content-Type", "text/css;charset=utf-8")
		css, _ := ioutil.ReadFile("ui/bootstrap/style.css")
		w.Write([]byte(css))
	} else if r.URL.String() == "/ui/bootstrap/assets/css/bootstrap.min.css" {
		w.Header().Set("Content-Type", "text/css;charset=utf-8")
		css, _ := ioutil.ReadFile("ui/bootstrap/assets/css/bootstrap.min.css")
		w.Write([]byte(css))
	} else if r.URL.String() == "/ui/bootstrap/assets/css/bootstrap-theme.min.css" {
		w.Header().Set("Content-Type", "text/css;charset=utf-8")
		css, _ := ioutil.ReadFile("ui/bootstrap/assets/css/bootstrap-theme.min.css")
		w.Write([]byte(css))
	} else if r.URL.String() == "/ui/bootstrap/assets/css/jquery-ui.min.css" {
		w.Header().Set("Content-Type", "text/css;charset=utf-8")
		css, _ := ioutil.ReadFile("ui/bootstrap/assets/css/jquery-ui.min.css")
		w.Write([]byte(css))
	} else if r.URL.String() == "/ui/bootstrap/assets/css/jquery-ui.structure.min.css" {
		w.Header().Set("Content-Type", "text/css;charset=utf-8")
		css, _ := ioutil.ReadFile("ui/bootstrap/assets/css/jquery-ui.structure.min.css")
		w.Write([]byte(css))
	} else if r.URL.String() == "/ui/bootstrap/assets/js/bootstrap-3.3.5.min.js" {
		w.Header().Set("Content-Type", "text/javascript;charset=utf-8")
		css, _ := ioutil.ReadFile("ui/bootstrap/assets/js/bootstrap-3.3.5.min.js")
		w.Write([]byte(css))
	} else if r.URL.String() == "/ui/bootstrap/assets/js/jquery-1.11.3.min.js" {
		w.Header().Set("Content-Type", "text/javascript;charset=utf-8")
		css, _ := ioutil.ReadFile("ui/bootstrap/assets/js/jquery-1.11.3.min.js")
		w.Write([]byte(css))
	} else if r.URL.String() == "/ui/bootstrap/assets/js/jquery-ui-1.11.4.min.js" {
		w.Header().Set("Content-Type", "text/javascript;charset=utf-8")
		css, _ := ioutil.ReadFile("ui/bootstrap/assets/js/jquery-ui-1.11.4.min.js")
		w.Write([]byte(css))
	} else if r.URL.String() == "/ui/bootstrap/assets/fonts/glyphicons-halflings-regular.eot" {
		w.Header().Set("Content-Type", "application/vnd.ms-fontobject")
		css, _ := ioutil.ReadFile("ui/bootstrap/assets/fonts/glyphicons-halflings-regular.eot")
		w.Write([]byte(css))
	} else if r.URL.String() == "/ui/bootstrap/assets/fonts/glyphicons-halflings-regular.svg" {
		w.Header().Set("Content-Type", "image/svg+xml")
		css, _ := ioutil.ReadFile("ui/bootstrap/assets/fonts/glyphicons-halflings-regular.svg")
		w.Write([]byte(css))
	} else if r.URL.String() == "/ui/bootstrap/assets/fonts/glyphicons-halflings-regular.ttf" {
		w.Header().Set("Content-Type", "application/x-font-ttf")
		css, _ := ioutil.ReadFile("ui/bootstrap/assets/fonts/glyphicons-halflings-regular.ttf")
		w.Write([]byte(css))
	} else if r.URL.String() == "/ui/bootstrap/assets/fonts/glyphicons-halflings-regular.woff" {
		w.Header().Set("Content-Type", "application/font-woff")
		css, _ := ioutil.ReadFile("ui/bootstrap/assets/fonts/glyphicons-halflings-regular.woff")
		w.Write([]byte(css))
	} else if r.URL.String() == "/ui/bootstrap/assets/fonts/glyphicons-halflings-regular.woff2" {
		w.Header().Set("Content-Type", "application/font-woff2")
		css, _ := ioutil.ReadFile("ui/bootstrap/assets/fonts/glyphicons-halflings-regular.woff2")
		w.Write([]byte(css))
	} else if r.URL.String() == "/GetRunningActivity" {
		getRunningActivity(w, r)
	} else if r.URL.String() == "/StartActivity" {
		startActivity(w, r)
	} else if r.URL.String() == "/StopActivity" {
		stopActivity(w, r)
	} else if r.URL.String() == "/ListActivities" {
		listActivities(w, r)
	} else if r.URL.String() == "/ListProjects" {
		listProjects(w, r)
	} else if r.URL.String() == "/ProjectNameAutocomplete" {
		projectsAutocomplete(w, r)
	} else if r.URL.String() == "/CountOfProjectActivities" {
		countOfProjectActivities(w, r)
	} else if r.URL.String() == "/RemoveActivity" {
		removeActivity(w, r)
	} else if r.URL.String() == "/RemoveProject" {
		removeProject(w, r)
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
// TODO Using command-line arguments provide several UI types (just 'bootstrap' at the moment)
// TODO We need some security on requests/responses...
// TODO We need track requests/responses...
func main() {
	log.Println("Entering main function...")

	log.Println(os.Args);

	var tplType string
	if len(os.Args) > 2 {
		usage()
	} else if len(os.Args) == 2 {
		if os.Args[1] == "--help" {
			usage()
		} else if strings.HasPrefix(os.Args[1], "--type=") == true {
			tplType = strings.Replace(os.Args[1], "--type=", "", 1)
		}
	}

	// Ensure template type is correct
	if tplType == "angularjs" || tplType == "bootstrap" || tplType == "polymer" {
		templateType = tplType
	}

	// TODO Theme with Polymer used is not implemented yet!
	if templateType == "angularjs" || templateType == "polymer" {
		log.Println("TODO Theme with Polymer used is not implemented yet!")
		templateType = "bootstrap"
	}

	var h httpHandler
	// TODO User should have possibility of changing this (configuration file)!
	http.ListenAndServe("localhost:4000", h)
}

// Get database path
func getDatabasePath() string {
	usr, err := user.Current()
	if err != nil {
		log.Println("Unable to create correct database path. Using in-memory database!")
		return ":memory:"
	}

	// TODO User should have possibility of changing this (configuration file)!
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
// TODO Check if there is not any running activity already!
func startActivity(w http.ResponseWriter, r *http.Request) {
	db, err := database.InitStorage(dbPath)
	defer db.Close()
	if err != nil {
		outputJson(jsonrpc.NewErrorResponse(jsonrpc.InitStorageError, "id-XXX"), w)
		return
	}

	// Firstly we need to check if there is no running activity
	ra, err := database.SelectActivityRunning(db)
	if ra.ActivityId > 0 {
		outputJson(jsonrpc.NewErrorResponse(jsonrpc.AnotherRunningActivityError, "id"), w)
		return
	}

	r.ParseForm()
	var project database.Project
	projectName := r.FormValue("project")
	projects, _ := database.SelectProjectByName(db, projectName)
	// Note: We don't bother about error - in that case we just create new project.
	//checkError(err)

	if len(projects) >= 1 {
		project = projects[0]
	} else if len(projects) == 0 {
		p, err := database.InsertProject(db, projectName, "")
		if err != nil {
			outputJson(jsonrpc.NewErrorResponse(jsonrpc.NewProjectError, "id-XXX"), w)
			return
		}
		project = p
	}

	var a database.Activity
	a, err = database.InsertActivity(db, project.ProjectId, r.FormValue("name"),
		r.FormValue("description"), r.FormValue("tags"))
	if err != nil {
		outputJson(jsonrpc.NewErrorResponse(jsonrpc.NewActivityError, "id-XXX"), w)
		return
	}
	a.SetProject(project)

	var res = map[string]interface{}{
		"Message": "Activity was successfully started.",
		"Activity": a,
	}
	outputJson(jsonrpc.NewResponse(res, "id-XXX"), w)
}

// Stop activity.
func stopActivity(w http.ResponseWriter, r *http.Request) {
	db, err := database.InitStorage(dbPath)
	defer db.Close()
	if err != nil {
		outputJson(jsonrpc.NewErrorResponse(jsonrpc.InitStorageError, "id-XXX"), w)
		return
	}

	ra, err := database.SelectActivityRunning(db)
	if err != nil {
		outputJson(jsonrpc.NewErrorResponse(jsonrpc.NoRunningActivityError, "id-XXX"), w)
		return
	}

	ra.Stopped = time.Now().Format(time.RFC3339)
	_, err = database.UpdateActivity(db, ra)
	if err != nil {
		outputJson(jsonrpc.NewErrorResponse(jsonrpc.UpdateActivityError, "id-XXX"), w)
	}

	var msg = map[string]string{
		"Message": "Activity was successfully stopped.",
	}
	outputJson(jsonrpc.NewResponse(msg, "id-XXX"), w)
}

// Render JSON with details about currently running activity.
func getRunningActivity(w http.ResponseWriter, r *http.Request) {
	db, err := database.InitStorage(dbPath)
	defer db.Close()
	if err != nil {
		outputJson(jsonrpc.NewErrorResponse(jsonrpc.InitStorageError, "id"), w)
		return
	}

	activity, _ := database.SelectActivityRunning(db)
	json, err := json.Marshal(activity)
	checkError(err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// Render JSON with activities.
// TODO We need from and to for pagination!
func listActivities(w http.ResponseWriter, r *http.Request) {
	db, err := database.InitStorage(dbPath)
	defer db.Close()
	if err != nil {
		outputJson(jsonrpc.NewErrorResponse(jsonrpc.InitStorageError, "id-XXX"), w)
		return
	}

	activities, err := database.SelectActivities(db, 20)
	checkError(err)

	json, err := json.Marshal(activities)
	checkError(err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// Render JSON with projects.
func listProjects(w http.ResponseWriter, r *http.Request) {
	db, err := database.InitStorage(dbPath)
	defer db.Close()
	if err != nil {
		outputJson(jsonrpc.NewErrorResponse(jsonrpc.InitStorageError, "id-XXX"), w)
		return
	}

	projects, err := database.SelectProjects(db, -1)
	if err != nil {
		outputJson(jsonrpc.NewErrorResponse(jsonrpc.InitStorageError, "id-XXX"), w)
	}

	json, err := json.Marshal(projects)
	checkError(err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// Provide auto-complete for projects.
func projectsAutocomplete(w http.ResponseWriter, r *http.Request) {
	db, err := database.InitStorage(dbPath)
	defer db.Close()
	if err != nil {
		outputJson(jsonrpc.NewErrorResponse(jsonrpc.InitStorageError, "id-XXX"), w)
		return
	}

	r.ParseForm()
	term := r.FormValue("term")

	projects, err := database.SelectProjectWithTerm(db, term)
	if err != nil {
		outputJson(jsonrpc.NewErrorResponse(jsonrpc.InitStorageError, "id-XXX"), w)
	}

	json, err := json.Marshal(projects)
	checkError(err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// Return count of activities of given project.
func countOfProjectActivities(w http.ResponseWriter, r *http.Request) {
	// TODO Implement this
	data := struct { ActivitiesCount int64 }{ 0, }
	json, err := json.Marshal(data)
	checkError(err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// Remove activity.
func removeActivity(w http.ResponseWriter, r *http.Request) {
	outputJson(jsonrpc.NewErrorResponse(jsonrpc.NotImplementedError, "id-XXX"), w)
}

// Remove project.
func removeProject(w http.ResponseWriter, r *http.Request) {
	outputJson(jsonrpc.NewErrorResponse(jsonrpc.NotImplementedError, "id-XXX"), w)
}

// Render main HTML page.
func mainPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"Name": "odTimeTracker"}

	p := "ui/" + templateType + "/"
	tpl, err := template.ParseFiles(p +  "main.tmpl", p + "header.tmpl", p + "footer.tmpl")
	checkError(err)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tpl.Execute(w, data)
	checkError(err)
}

// Print usage information.
func usage() {
	fmt.Println(appInfo)
	fmt.Printf("\n%s\n\n", appDescription)
	fmt.Printf("Usage:\n\n")
	fmt.Printf("%s --help         Print this help\n", appShortName)
	fmt.Printf("%s --type=[TYPE]  Use template of given type\n\n", appShortName)
	fmt.Printf("Available template types are: bootstrap\n\n")
	os.Exit(0)
}

// Helper function for printing Json.
func outputJson(data interface{}, w http.ResponseWriter) {
	log.Println(data)

	json, err := json.Marshal(data)
	checkError(err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
