package dashboard

import (
	"net/http"

    "brypt-server/internal/handlebars"
    "brypt-server/api/access"

    "github.com/go-chi/chi"
)

type Resources struct{}

/* **************************************************************************
** Function: Routes
** Description: Register the dashboard resources with chi router and returns the
** built router.
** *************************************************************************/
func (rs Resources) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get( "/", access.CheckAuth( rs.Index ) )	// Implemetation of base dashboard page

	return r
}

/* **************************************************************************
** Function: Index
** URI: dashboard.host (GET)
** Description: Handles compilation the access dashboard page
** Client: Displays the compiled page/
** *************************************************************************/
func (rs Resources) Index(w http.ResponseWriter, r *http.Request) {

	dashboardCTX := make( map[string]interface{} )

    dashboardCTX["title"] = "Brypt"

	page := handlebars.CompilePage( "dashboard", dashboardCTX )

	w.Header().Set( "Content-Type", "text/html" )
	w.Write( []byte( page ) )

}
