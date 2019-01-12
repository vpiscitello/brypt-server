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

	"golang.org/x/crypto/bcrypt"

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

	du = retCTX["ret"].(db.User)

	if err != nil {
		fmt.Println(err)
		return du, err
	} else {
		plainTextPW := []byte(fmt.Sprintf("%s%s", "salt", password)) //TODO Use a salt
		//plainTextPW := []byte(fmt.Sprintf("%s%s", du.Password_salt, password))

		// TODO Actually get their hashed password from the database to compare, rather than generating it
		hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPW), 0)
		if err != nil {
			fmt.Println(err)
			return du, err
		}

		err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(plainTextPW))
		if err == nil {
			fmt.Printf("User %s authorized\n", username)
			//fmt.Printf("User %s authorized\n", du.Username)
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

	//r.Get( "/", rs.Index )
	//r.Get( "/login", rs.Login )
	//r.Get( "/register", rs.Register )
	//r.Get( "/link", CheckAuth(rs.Link) )

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

//	TestInsert()	// TODO: REMOVE WHEN FINISHED TESTING DB INSERT
	TestUpdate()	// TODO: REMOVE WHEN FINSHED TESTING DB UPDATE, FIX
	TestDelete()	// TODO: REMOVE WHEN FINISHED TESTING DB DELETE
	TestFind()		// TODO: REMOVE WHEN FINISHED TESTING DB FIND, FIX

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
	//finish parsing this

	username := "m@llory5" //Hard-coded for temporary use, replace with parameters
	password := "pswd" //Hard-coded for temporary use, replace with parameters
	du, err := identifyUser(w, username, password)
	if err != nil {
		w.Write( []byte( "Could not login...\n" ) )
		return
	}

	// On success, add a cookie
	SetCookieHandler(w, r, du.Uid)
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

func TestInsert2(w http.ResponseWriter) string {
	// db.Connect()	

	objID1 := objectid.New().Hex()
	objID2 := objectid.New().Hex()
	objID3 := objectid.New().Hex()
//	var login_attempts int32 = 4
	testCTX := make( map[string]interface{} )
	testCTX["username"] = "m@llory5"
	testCTX["first_name"] = "Mallory"
	testCTX["last_name"] = "Allen"
	testCTX["region"] = "Wonderland"
	testCTX["age"] = time.Now().Round(time.Millisecond)
	testCTX["login_attempts"] = 1 
	testCTX["networks"] = []string{objID1, objID2, objID3}
	id := db.Write("brypt_users", testCTX)
	print("\nid: ")
	fmt.Print(id)
	return id
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

func TestFind2() map[string]interface{} {

	testCTX := make( map[string]interface{} )
	testCTX["username"] = "m@llory5"
//	testCTX["username"] = "TotallyTom"
//	testCTX["first_name"] = "Alice"
//	testCTX["last_name"] = "Allen"

	/**********FIND ALL TEST**************/
	retCTX, err := db.FindOne("brypt_users", testCTX)
	
	print("\nFind All results: \n")
	fmt.Printf("%+v\n\n\n", retCTX)
	fmt.Printf("%+v\n", retCTX["ret"])
	
	print("\nFind All error response: ")
	fmt.Println(err)
	return retCTX

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
