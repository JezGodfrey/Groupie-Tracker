# groupie-tracker

## Overview

This project is a web server written in Go that makes use of an API to display information on popular music artists. The data from the API is unmarshaled and collected into corresponding data types, and displayed onto the page with HTML and CSS.

## Usage: how to run

Simply run the program like so to start the server:

```sh
go run .
```

Then visit `localhost:8080` on your device. On the main page, you can either follow the links to the artists' individual pages (by ID), or search for the artist directly in the search bar (by name). Trying to search for an artist not in the API will return a 404-page-not-found error page. Below is an example of the main page, and an individual artist's page - Scorpions - detailing information about the artist.


![IndexPage](/IndexPage.PNG?raw=true "Index Page")
![ScorpionsPage](/Scorpions.PNG?raw=true "Scorpions Page")

## Authors

This program was written by Jez Godfrey as part of the 01 Founders fellowship
