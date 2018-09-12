package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

//Compile templates on start
var templates = template.Must(template.ParseFiles(
	"templates/notFound.html",
	"templates/header.html",
	"templates/footer.html",
	"templates/index.html"))

func fire() {
	currentTime := time.Now()
	fmt.Println("{" + currentTime.Format(time.RFC1123) + "} FIRE!!!")
}

func saveToken(token string) {
	currentTime := time.Now()
	fmt.Println("{" + currentTime.Format(time.RFC1123) + "} saving token: " + token)
}

//Display the named template
func display(w http.ResponseWriter, tmpl string, data interface{}) {
	templates.ExecuteTemplate(w, tmpl, data)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := Page{
		PageTitle: "Home",
	}
	display(w, "index", data)
}

func randomPageHandler(w http.ResponseWriter, r *http.Request) {
	// If empty show the home page
	// If static page show the static package
	// Else show the 404 page
	if r.URL.Path == "/" {
		homeHandler(w, r)
	} else if r.URL.Path == "/fire" {
		fireHandler(w, r)
	} else if r.URL.Path == "/token" {
		tokenHandler(w, r)
	} else if strings.HasSuffix(r.URL.Path[1:], ".html") {
		http.ServeFile(w, r, "static/html/"+r.URL.Path[1:])
	} else {
		fmt.Println("Sorry but it seems this page does not exist...")
		errorHandler(w, r, http.StatusNotFound)
	}
}

func fireHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fire()
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		saveToken(r.FormValue("streamlabs"))
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		display(w, "404", &Page{PageTitle: "404"})
	}
}

// Gets a random value from the low to high values. This will include the low and high values.
func randomValue(low int, high int) int {
	var scaledInt = high - low + 1 // The +1 is to offset the values so it can be the high value.
	return rand.Intn(scaledInt) + low
}

func getPort() string {
	if value, ok := os.LookupEnv("PORT"); ok {
		return ":" + value
	}
	return ":8080"
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", randomPageHandler)
	var port = getPort()
	fmt.Println("Now listening to port " + port)
	log.Fatal(http.ListenAndServe(port, nil))
}

//A Page structure
type Page struct {
	PageTitle string
}
