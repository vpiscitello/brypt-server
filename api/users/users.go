package users

import (
	"net/http"

	"github.com/go-chi/chi"
)

type Resources struct{}

/* **************************************************************************
** Function: Routes
** Description: Register the user resources with chi router and returns the
** built router.
** *************************************************************************/
func (rs Resources) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", rs.List)    // GET /todos - read a list of todos
	r.Post("/", rs.Create) // POST /todos - create a new todo and persist it
	r.Put("/", rs.Delete)
	r.Delete("/", rs.Delete) // DELETE /todos/{id} - delete a single todo by :id

	return r
}

/* **************************************************************************
** Function: Index
** URI: access.host (GET)
** Description: Handles compilation the access index page which displays login
** and registration forms.
** Client: Displays the compiled page/
** *************************************************************************/
func (rs Resources) List(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("list of stuff.."))
}

/* **************************************************************************
** Function: Index
** URI: access.host (GET)
** Description: Handles compilation the access index page which displays login
** and registration forms.
** Client: Displays the compiled page/
** *************************************************************************/
func (rs Resources) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create"))
}

/* **************************************************************************
** Function: Index
** URI: access.host (GET)
** Description: Handles compilation the access index page which displays login
** and registration forms.
** Client: Displays the compiled page/
** *************************************************************************/
func (rs Resources) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get"))
}

/* **************************************************************************
** Function: Index
** URI: access.host (GET)
** Description: Handles compilation the access index page which displays login
** and registration forms.
** Client: Displays the compiled page/
** *************************************************************************/
func (rs Resources) Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("update"))
}

func (rs Resources) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("delete"))
}
