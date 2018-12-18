package access

import (
	// "fmt"
	db "brypt-server/api/database"
	"net/http"
	"time"
    "brypt-server/internal/handlebars"

	"github.com/go-chi/chi"
	"github.com/mongodb/ftdc/bsonx/objectid"
	// "github.com/aymerick/raymond"

	// "brypt-server/api/users"
)

type Resources struct{}

/* **************************************************************************
** Function: Routes
** Description: Register the access resources with chi router and returns the
** built router.
** *************************************************************************/
func (rs Resources) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get( "/", rs.Index )	// Implemetation of base access page which will support login and registration actions
	r.Post( "/login", rs.Login )		// Post request for user login
	r.Post( "/register", rs.Register )	// Post request for registering an account
	r.Post( "/link", rs.Link )	// Post request for linking a device to a user account

	return r
}

/* **************************************************************************
** Function: Index
** URI: access.host (GET)
** Description: Handles compilation the access index page which displays login
** and registration forms.
** Client: Displays the compiled page/
** *************************************************************************/
func (rs Resources) Index(w http.ResponseWriter, r *http.Request) {
	
	TestInsert(w, r)	// TODO: REMOVE WHEN FINISHED TESTING DB INSERT
	TestDelete(w, r)	// TODO: REMOVE WHEN FINISHED TESTING DB DELETE
	
	action := r.URL.Query().Get( "action" )
	accessCTX := make( map[string]interface{} )

	switch action {
		default:
			accessCTX["login"] = ""
			accessCTX["register"] = "bck"
			accessCTX["active"] = "log"
			accessCTX["inactive_text"] = "Register"
		case "register":
			accessCTX["login"] = "bck"
			accessCTX["register"] = ""
			accessCTX["active"] = "reg"
			accessCTX["inactive_text"] = "Login"
	}

	page := handlebars.CompilePage( "access", accessCTX )

	w.Header().Set( "Content-Type", "text/html" )
	w.Write( []byte( page ) )

}

/* **************************************************************************
** Function: Login
** URI: access.host/login (POST)
** Description: Handles login post request. Expects a user's login information.
** If successful an authenticated session cookie will be returned otherwise
** a error will be returned to a client.
** Client: If success the client should redirect to the dashboard page otherwise
** an error message should be displayed.
** *************************************************************************/
func (rs Resources) Login(w http.ResponseWriter, r *http.Request) {
    w.Write( []byte( "Login...\n" ) )
}

/* **************************************************************************
** Function: Register
** URI: access.host/register (POST)
** Description: Handles register post request. Expects a user's registration
** information. If successful a success message will be sent otherwise an
** error message.
** Client: Displays registration statsus message.
** *************************************************************************/
func (rs Resources) Register(w http.ResponseWriter, r *http.Request) {
	w.Write( []byte( "Register...\n" ) )
}

/* **************************************************************************
** Function: Link
** URI: access.host/link (POST)
** Description: Handles link post request. Expects user session information and
** device information. If successful a success message will be sent otherwise an
** error message.
** Client: Displays link statsus message.
** *************************************************************************/
func (rs Resources) Link(w http.ResponseWriter, r *http.Request) {
	w.Write( []byte( "Linking...\n" ) )
}

/* **************************************************************************
** Function: TestInsert
** Description: Just a test function to demonstrate db insert functionality
**	TODO: Remove when finished testing db insert
** *************************************************************************/
func TestInsert(w http.ResponseWriter, r *http.Request) {
	// db.Connect()	

	objID1 := objectid.New()
	objID2 := objectid.New()
	objID3 := objectid.New()
//	var login_attempts int32 = 4
	testCTX := make( map[string]interface{} )
	testCTX["username"] = "AwesomeAlice"
	testCTX["first_name"] = "Alice"
	testCTX["last_name"] = "Allen"
	testCTX["region"] = "Wonderland"
	testCTX["age"] = time.Now().Round(time.Millisecond)
	testCTX["login_attempts"] = 4
	testCTX["networks"] = []objectid.ObjectID{objID1, objID2, objID3}
	db.ReqHandler(w, r, "PUT", "brypt_users", testCTX)

	testCTX["username"] = "TotallyTom"
	testCTX["first_name"] = "Alice"
	testCTX["last_name"] = "Allen"
	testCTX["region"] = "Wonderland"
	testCTX["age"] = time.Now().Round(time.Millisecond)
	testCTX["login_attempts"] = 4
	testCTX["networks"] = []objectid.ObjectID{objID1, objID2, objID3}
	db.ReqHandler(w, r, "PUT", "brypt_users", testCTX)
//	defer db.Disconnect()	// Causes an internal server error for some reason...
}

func TestDelete(w http.ResponseWriter, r *http.Request) {

	testCTX := make( map[string]interface{} )
	testCTX["username"] = "AwesomeAlice"
	testCTX["first_name"] = "Alice"
	testCTX["last_name"] = "Allen"
	db.ReqHandler(w, r, "DELETE", "brypt_users", testCTX)
}
