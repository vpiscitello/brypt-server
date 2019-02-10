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
	"github.com/gorilla/securecookie"

	"golang.org/x/crypto/bcrypt"
	// "github.com/aymerick/raymond"

	// "brypt-server/api/users"
)

type Resources struct{}

var cookieName = "user_authentication"

var hashKey = []byte( securecookie.GenerateRandomKey( 32 ) )
var blockKey = []byte( securecookie.GenerateRandomKey( 32 ) )

var sc = securecookie.New( hashKey, blockKey )

/* **************************************************************************
** Function: ReadCookieHandler
** Description: Reads the cookie from a http request and validates that cookie.
** *************************************************************************/
func ReadCookieHandler(r *http.Request) bool {
	cookie, err := r.Cookie(cookieName)
	if err == nil {
		value := make(map[string]string)
		if err = sc.Decode(cookieName, cookie.Value, &value); err == nil {
			fmt.Printf("Decoded Cookie: %#s\n", cookie)
			return true
		}
		return false
	}
	fmt.Printf("No matching cookie: %s\n", err)
	return false
}

/* **************************************************************************
** Function: SetCookieHandler
** Description: Generates an encoded value for the cookie, and sets the cookie
in the http response header.
** *************************************************************************/
func SetCookieHandler(w http.ResponseWriter, r *http.Request, id string) {
	value := map[string]string{
		"id": id,
	}

	if encoded, err := sc.Encode(cookieName, value); err == nil {
		expiration := time.Now().AddDate(0, 0, 1)
		cookie := &http.Cookie{
			Name:  cookieName,
			Value: encoded,
			Path:  "/",
			Secure: true,
			Expires: expiration,
		}
		http.SetCookie(w, cookie)
		fmt.Printf("Set Cookie: %#s\n", cookie)
	}
}

/* **************************************************************************
** Function: identifyUser
** Description: Finds a user in the database and validates the password provided
against the hashed value stored in the database
** *************************************************************************/
func identifyUser(w http.ResponseWriter, username string, password string) (db.User, error) {
	testCTX := make( map[string]interface{} )
	testCTX["username"] = username

	du := db.User{}

	retCTX, err := db.FindOne("brypt_users", testCTX)
	if err != nil {
		fmt.Println(err)
		return du, err
	}

	fmt.Print("RETCTX: ")
	fmt.Println(retCTX)

	du = retCTX["ret"].(db.User)

	if err != nil {
		fmt.Println(err)
		return du, err
	} else {
		plainTextPW := []byte(password)

		err = bcrypt.CompareHashAndPassword([]byte(du.Password), []byte(plainTextPW))
		if err == nil {
			fmt.Printf("User %s authorized\n", username)
			return du, err
		} else {
			fmt.Println(err)
			return du, err
		}
	}
	return du, err
}

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

	//TestInsert()	// TODO: REMOVE WHEN FINISHED TESTING DB INSERT
	//TestUpdate()	// TODO: REMOVE WHEN FINSHED TESTING DB UPDATE, FIX
	//TestDelete()	// TODO: REMOVE WHEN FINISHED TESTING DB DELETE
	//TestFind()		// TODO: REMOVE WHEN FINISHED TESTING DB FIND, FIX

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

	fmt.Print("\n\n")
	fmt.Print(bodyBytes)
	fmt.Print("\n\n")

	loginCTX := make( map[string]interface{} )
	if err := json.Unmarshal(bodyBytes, &loginCTX); err != nil {
		fmt.Print("Sadness: ")
		fmt.Println(err)
	}
	fmt.Print("Register: ")
	fmt.Println(loginCTX)

	du, err := identifyUser(w, loginCTX["username"].(string), loginCTX["password"].(string))
	if err != nil {
		w.Write( []byte( "Could not login...\n" ) )
		return
	}

	// On success, add a cookie
	SetCookieHandler(w, r, du.Uid)
	w.Write([]byte("Logged in"))

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
	fmt.Print("Register: ")
	fmt.Println(regCTX)
	regCTX["time_registered"] = time.Now().Round(time.Millisecond)

	if checkUserRegistration(regCTX["username"].(string)) {
		fmt.Println("Username already registered")
		w.Write([]byte("Cannot register"))
		return
	}

	userCTX := make(map[string]interface{})
	userCTX["username"] = regCTX["username"]
	userCTX["region"] = regCTX["Region"]
	userCTX["first_name"] = regCTX["first_name"]
	userCTX["last_name"] = regCTX["last_name"]
	//userCTX["birthdate"] = regCTX["Birthday"]
	userCTX["email"] = regCTX["email"]
	//userCTX["join_date"] = regCTX["time_registered"]

	hash, err := bcrypt.GenerateFromPassword([]byte(regCTX["password"].(string)), 0)
	if err != nil {
		fmt.Print("Error hashing password: ")
		fmt.Println(err)
		w.Write([]byte("Error registering"))
	}
	userCTX["password"] = string(hash)
	fmt.Println("This is userCTX:")
	fmt.Println(userCTX)

	db.Write("brypt_users", userCTX)

	w.Write([]byte("Registered!"))
}

/* **************************************************************************
** Function: checkUserRegistration
** Description: Check if there is already a user registered with a certain username
** Returns: True if there is already a user registered under a username
** *************************************************************************/
func checkUserRegistration(username string) bool {

	checkUsrCTX := make( map[string]interface{} )
	checkUsrCTX["username"] = username
	retCTX, err := db.FindAll("brypt_users", checkUsrCTX)
	if err != nil {
		fmt.Print("ERROR: ")
		fmt.Println(err)
	}
	du := retCTX["ret"].([]db.User)
	fmt.Print("User: ")
	fmt.Println(du)
	fmt.Print("Length: ")
	fmt.Println(len(du))
	if len(du) > 0 {
		return true
	}
	return false

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
