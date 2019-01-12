package access

import (
	"fmt"
	db "brypt-server/api/database"
	"net/http"
	"time"
    "brypt-server/internal/handlebars"
	"encoding/json"
	"io/ioutil"

	"github.com/go-chi/chi"
	"github.com/mongodb/ftdc/bsonx/objectid"
	// "github.com/aymerick/raymond"

	"golang.org/x/crypto/bcrypt"

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
	
	TestInsert()	// TODO: REMOVE WHEN FINISHED TESTING DB INSERT
//	TestUpdate()	// TODO: REMOVE WHEN FINSHED TESTING DB UPDATE, FIX
//	TestDelete()	// TODO: REMOVE WHEN FINISHED TESTING DB DELETE
//	TestFind()		// TODO: REMOVE WHEN FINISHED TESTING DB FIND, FIX

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
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Terrible errors in Login\n")
		return
	}
	bodyString := string(bodyBytes)
	w.Write([]byte(bodyString))
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
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Terrible errors in Register\n")
		return
	}
	regCTX := make( map[string]interface{} )
	if err := json.Unmarshal(bodyBytes, &regCTX); err != nil {
		fmt.Println("Sadness\n")
	}
	regCTX["time_registered"] = time.Now().Round(time.Millisecond)
	fmt.Println("This is regCTX:")
	fmt.Println(regCTX)

	regCTX["password"], err = bcrypt.GenerateFromPassword([]byte(regCTX["password"].(string)), 10)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte("Error registering"))
	}
	fmt.Println("This is regCTX after:")
	fmt.Println(regCTX)
	//regCTX := make( map[string]interface{} )
	//regCTX["username"] = dat["username"]
	id := db.Write("brypt_users", regCTX)
	print("\nnil id: ")
	fmt.Println(id)
	w.Write([]byte("Registered!"))
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
func TestInsert() {
	// db.Connect()	

	objID1 := objectid.New().Hex()
	objID2 := objectid.New().Hex()
	objID3 := objectid.New().Hex()
//	var login_attempts int32 = 4
	testCTX := make( map[string]interface{} )
	testCTX["username"] = "m@llory6"
	testCTX["first_name"] = "Mal"
	testCTX["last_name"] = "Allen"
	testCTX["region"] = "Wonderland"
	testCTX["birthdate"] = time.Now().Round(time.Millisecond)
	testCTX["login_attempts"] = 1
	testCTX["networks"] = []string{objID1, objID2, objID3}
	testCTX["password"] = "qwerty"
	id := db.Write("brypt_usrs", testCTX)	// Incorrect collection name (should return nilObjectID)
	print("\nnil id: ")
	fmt.Print(id)
	id = db.Write("brypt_users", testCTX)
	print("\nid: ")
	fmt.Print(id)
	id = db.Write("brypt_users", testCTX)
	print("\nid: ")
	fmt.Print(id)

/*	testCTX["username"] = "TotallyTom"
	testCTX["first_name"] = "Alice"
	testCTX["last_name"] = "Allen"
	testCTX["region"] = "Wonderland"
	testCTX["birthdate"] = time.Now().Round(time.Millisecond)
	testCTX["login_attempts"] = 4
	testCTX["networks"] = []objectid.ObjectID{objID1, objID2, objID3}
	id2 := db.Write(w, "brypt_users", testCTX)
	print("\nid2: ")
	fmt.Print(id2)*/
//	defer db.Disconnect()	// Causes an internal server error for some reason...
}

func TestDelete() {

	testCTX := make( map[string]interface{} )
//	testCTX["username"] = "AwesomeAlice"
//	testCTX["first_name"] = "Alice"
	testCTX["last_name"] = "Allen"
	err := db.DeleteOne("brypt_users", testCTX)
	print("\nDelete One error response: ")
	fmt.Print(err)
//	err = db.DeleteAll("brypt_users", testCTX)
//	print("\nDelete All error response: ")
//	fmt.Print(err)
}

func TestFind() {

	testCTX := make( map[string]interface{} )
	testCTX["username"] = "m@llory5"
//	testCTX["username"] = "TotallyTom"
//	testCTX["first_name"] = "Alice"
//	testCTX["last_name"] = "Allen"

	/**********FIND ALL TEST**************/
	retCTX, err := db.FindAll("brypt_users", testCTX)
	
	print("\nFind All results: \n")
	fmt.Printf("%+v\n", retCTX)
	
	print("\nFind All error response: ")
	fmt.Println(err)

	/**********FIND ONE TEST**************/
	testCTX["username"] = "notInDB"
	retCTX, err = db.FindOne("brypt_users", testCTX)
	
	print("\nFind One result:\n ")
	fmt.Printf("%+v\n", retCTX["ret"])
	
	print("\nFind One error response: ")
	fmt.Println(err)

}

func TestUpdate() {
	testCTX := make( map[string]interface{} )
	testCTX["username"] = "TotallyTom"
	testCTX["first_name"] = "Alice"
	testCTX["last_name"] = "Allen"
	
	updateFieldCTX := make( map[string]interface{} )
	updateFieldCTX["username"] = "Re@llyTom" 
	updateFieldCTX["first_name"] = "Tom"
	updateCTX := make( map[string]interface{} )
	updateCTX["$set"] = updateFieldCTX
	//	updateCTX["first_name"] = "Tom"
	
	err := db.UpdateOne("brypt_users", testCTX, updateCTX)
	print("\nUpdate One response: ")
	fmt.Print(err)
}
