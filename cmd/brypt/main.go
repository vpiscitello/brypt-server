//package main

// Setup MySql Connection
// Setup MongoDB Connection
// Register Routes
// Run Server

//import {
//    "fmt"
//    "log"
//    "time"
//    "net/http"
//    "encoding/json"
//
//    "github.com/go-chi/chi"
//    "github.com/go-chi/chi/middleware"
//}
//
//func main()  {
//    r := chi.NewRouter()
//
//    r.Use( middleware.requestID )
//    r.Use( middleware.realIP )
//    r.Use( middleware.Logger )
//    r.Use( middleware.Recoverer )
//
//    router.Get( "/", func( w http.ResponseWriter, r *http.requestID )  {
//        w.Write( []byte( "." ) )
//    } )
//
//    // router.Mount( "/route", routeResources{}.Routes() )
//
//    http.ListenAndServe( ":3005", r )
//
//}
//How to use SubDomains in Golang , Subdomains With Go , http://codepodu.com/subdomains-with-golang/
// subdomains.go
//
// Please read http://codepodu.com/subdomains-with-golang/
// It's just copy and paste :smile:
//
//
// URLs :
// http://admin.localhost:8080/admin/pathone
// http://admin.localhost:8080/admin/pathtwo
// http://analytics.localhost:8080/analytics/pathone
// http://analytics.localhost:8080/analytics/pathtwo
//

package main

import (
	"fmt"
	"net/http"
	"strings"
)

type Subdomains map[string]http.Handler

func (subdomains Subdomains) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	domainParts := strings.Split(r.Host, ".")

	if mux := subdomains[domainParts[0]]; mux != nil {
		// Let the appropriate mux serve the request
		mux.ServeHTTP(w, r)
	} else {
		// Handle 404
		http.Error(w, "Not found", 404)
	}
}

type Mux struct {
	http.Handler
}

func (mux Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux.ServeHTTP(w, r)
}

func adminHandlerOne(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "It's adminHandlerOne , Hello, %q", r.URL.Path[1:])
}

func adminHandlerTwo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "It's adminHandlerTwo , Hello, %q", r.URL.Path[1:])
}

func analyticsHandlerOne(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "It's analyticsHandlerOne , Hello, %q", r.URL.Path[1:])
}

func analyticsHandlerTwo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "It's analyticsHandlerTwo , Hello, %q", r.URL.Path[1:])
}

func main() {
	adminMux := http.NewServeMux()
	adminMux.HandleFunc("/admin/pathone", adminHandlerOne)
	adminMux.HandleFunc("/admin/pathtwo", adminHandlerTwo)

	analyticsMux := http.NewServeMux()
	analyticsMux.HandleFunc("/analytics/pathone", analyticsHandlerOne)
	analyticsMux.HandleFunc("/analytics/pathtwo", analyticsHandlerTwo)

	subdomains := make(Subdomains)
	subdomains["admin"] = adminMux
	subdomains["analytics"] = analyticsMux

	http.ListenAndServe(":8080", subdomains)
}
