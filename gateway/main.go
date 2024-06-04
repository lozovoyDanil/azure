package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

const (
	moviesServiceAddr = "http://fastapiproject:5000"
	actorsServiceAddr = "http://fastapiproject2:5001"
	authServiceAddr   = "http://prj:5003"
)

var backendServices = map[string]string{
	//Movies endpoints
	"/api/movies/":              moviesServiceAddr,
	"/api/movies":               moviesServiceAddr,
	"/api/movies/filterByActor": moviesServiceAddr,
	"/api/movies/searchByTitle": moviesServiceAddr,
	"/api/admin/movies/add":     moviesServiceAddr,
	"/api/admin/movies/":        moviesServiceAddr,
	//Actors endpoints
	"/api/actors":                   actorsServiceAddr,
	"/api/actors/":                  actorsServiceAddr,
	"/api/actors/searchByLastName":  actorsServiceAddr,
	"/api/actors/searchByFirstName": actorsServiceAddr,
	"/api/actors/searchByFullName/": actorsServiceAddr,
	"/api/actors/movie/":            actorsServiceAddr,
	"/api/admin/actors/add":         actorsServiceAddr,
	"/api/admin/actors/":            actorsServiceAddr,
	//Auth endpoints
	"/api/favorites/": authServiceAddr,
	"/api/favorites":  authServiceAddr,
	"/api/sign-up":    authServiceAddr,
	"/api/sign-in":    authServiceAddr,
	"/api/healthz":    authServiceAddr,
	"/api/identity":   authServiceAddr,
}

func main() {
	http.HandleFunc("/", handleRequest)
	if err := http.ListenAndServe(":5050", nil); err != nil {
		log.Fatal(err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	for prefix, backendURL := range backendServices {
		if strings.HasPrefix(r.URL.Path, prefix) {
			remote, err := url.Parse(backendURL)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			proxy := httputil.NewSingleHostReverseProxy(remote)
			proxy.ServeHTTP(w, r)
			return
		}
	}

	http.Error(w, "Not found", http.StatusNotFound)
}
