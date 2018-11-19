package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

type accessResources struct{}

func (rs accessResources) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get( "/", rs.Index )
	r.Get( "/register", rs.Register )
	r.Get( "/login", rs.Login )
	r.Get( "/link", rs.Link )


	return r
}

func (rs accessResources) Index(w http.ResponseWriter, r *http.Request) {
    w.Write( []byte( "Access!\n" ) )
}

// Route handler for user registration
func (rs accessResources) Register(w http.ResponseWriter, r *http.Request) {
	w.Write( []byte( "Register..." ) )
}

// Route handler for user login
func (rs accessResources) Register(w http.ResponseWriter, r *http.Request) {
	w.Write( []byte( "Register..." ) )
}

// Route handler for linking a device to a user's account
func (rs accessResources) Link(w http.ResponseWriter, r *http.Request) {
	w.Write( []byte( "Linking..." ) )
}

