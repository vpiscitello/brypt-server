package dashboard

import (
	"net/http"

    "brypt-server/internal/handlebars"

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

	return r
}
