package access

import (
	"fmt"
	"net/http"
)

// authenticate node

/* **************************************************************************
** Function: CheckAuth
** Description: This function reads the cookie and determines whether or not
the user is authenticated. Returns a http.HandlerFunc.
** *************************************************************************/
func CheckAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if ReadCookieHandler(r) {
			fmt.Printf("Authorization Success!\n")
			h.ServeHTTP(w, r)
			return
		} else {
			fmt.Printf("Not authorized %d\n", 401)
		}

		http.Redirect(w, r, "https://access.localhost:3006", http.StatusUnauthorized)
	}
}
