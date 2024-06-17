package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
)

// Data types for unmarshaling JSON data
type Artist struct {
	Id           int
	Image        string
	Name         string
	Members      []string
	CreationDate float64
	FirstAlbum   string
	Locations    []string
	ConcertDates [][]string
}

type Locs struct {
	Id        int
	Locations []string
}

type Dates struct {
	Id             int
	DatesLocations map[string][]string
}

var ls []Locs
var ds []Dates
var artists []Artist

// Obtaining location data from the locations API
func getLocations() {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// This JSON is formatted beginning with "index"? Not sure what it means but this gets around it
	respData = respData[9 : len(respData)-2]

	json.Unmarshal(respData, &ls)
}

// Obtaining date data from the relations API - easier to do so from here than the dates API as the dates here are already mapped to locations
func getDates() {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// This JSON is formatted beginning with "index"? Not sure what it means but this gets around it
	respData = respData[9 : len(respData)-2]

	json.Unmarshal(respData, &ds)
}

// Obtaining artist data from the artists API
func getArtists() {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(respData, &artists)
}

// General handler for index and non-specified pages
func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		r.URL.Path = "/index.html"
	}

	// If page not found, display respective error page
	p, err := template.ParseFiles("gt" + r.URL.Path)
	if err != nil {
		errorHandler(w, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	p.ExecuteTemplate(w, "index.html", artists)
}

// Handler for displaying artist pages
func artistHandler(w http.ResponseWriter, r *http.Request) {
	// Handling if the request is made via the search bar
	if r.FormValue("myArtist") != "" {
		for i, a := range artists {
			if strings.ToLower(r.FormValue("myArtist")) == strings.ToLower(a.Name) {
				// If the artist web page is missing, display Internal Server Error
				p, err := template.ParseFiles("gt/artist.html")
				if err != nil {
					errorHandler(w, http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusOK)
				p.ExecuteTemplate(w, "artist.html", artists[i])
				return
			}
		}

		errorHandler(w, http.StatusNotFound)
		return
	}

	if r.URL.Path == "/artist" {
		errorHandler(w, http.StatusBadRequest)
		return
	}

	// Handling if the request is made via clicking the links on the main page
	id, err := strconv.Atoi(r.URL.Path[8:])
	if err != nil {
		errorHandler(w, http.StatusBadRequest)
		return
	}

	if id > artists[len(artists)-1].Id || id < 1 {
		errorHandler(w, http.StatusNotFound)
		return
	}

	p, err := template.ParseFiles("gt/artist.html")
	if err != nil {
		errorHandler(w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	p.ExecuteTemplate(w, "artist.html", artists[id-1])
}

// Handling errors
func errorHandler(w http.ResponseWriter, status int) {
	w.WriteHeader(status)

	switch status {
	case 400:
		p, err := template.ParseFiles("gt/400.html")
		if err != nil {
			errorHandler(w, http.StatusInternalServerError)
			return
		}
		p.Execute(w, "400.html")
	case 404:
		p, err := template.ParseFiles("gt/404.html")
		if err != nil {
			errorHandler(w, http.StatusInternalServerError)
			return
		}
		p.Execute(w, "404.html")
	case 500:
		http.Error(w, "500 - Internal Server Error", http.StatusInternalServerError)
	default:
		http.Error(w, "%v Error", status)
	}
}

func main() {
	// Extract data from api
	getArtists()
	getLocations()
	getDates()

	// Putting locations and respective concert dates into respective artist variables
	for i := range artists {
		artists[i].Locations = ls[i].Locations
		for _, l := range ls[i].Locations {
			artists[i].ConcertDates = append(artists[i].ConcertDates, ds[i].DatesLocations[l])
		}
	}

	// For implementing CSS
	fs := http.FileServer(http.Dir("./stylesheets"))
	http.Handle("/stylesheets/", http.StripPrefix("/stylesheets/", fs))

	// Handlers for HTML pages
	http.HandleFunc("/", handler)
	http.HandleFunc("/artist", artistHandler)
	http.HandleFunc("/artist/", artistHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
