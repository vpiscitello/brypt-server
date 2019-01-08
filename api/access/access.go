package access

import (
	"fmt"
	db "brypt-server/api/database"
	"net/http"
	"time"

    "brypt-server/internal/handlebars"

	"github.com/go-chi/chi"
	"github.com/mongodb/ftdc/bsonx/objectid"
	"github.com/gorilla/securecookie"

	"golang.org/x/crypto/bcrypt"
	// "github.com/aymerick/raymond"

	// "brypt-server/api/users"
)

type Resources struct{}

var cookieName = "testcookie"

var hashKey = []byte( securecookie.GenerateRandomKey( 32 ) )
var blockKey = []byte( securecookie.GenerateRandomKey( 32 ) )

var sc = securecookie.New( hashKey, blockKey )

func checkAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if ReadCookieHandler(r) {
			fmt.Printf("Success\n")
			h.ServeHTTP(w, r)
			return
		} else {
			fmt.Printf("Not authorized %s\n", 401)
		}

		h.ServeHTTP(w, r)
	}
}

func ReadCookieHandler(r *http.Request) bool {
	for _, cookie := range r.Cookies() {
		fmt.Print(cookie.Name)
	}
	fmt.Printf("\n\nCookie name: %#s\n", cookieName)
	cookie, err := r.Cookie(cookieName)
	if err == nil {
		fmt.Printf("Inside if statement\n")
		value := make(map[string]string)
		if err = sc.Decode(cookieName, cookie.Value, &value); err == nil {
			fmt.Printf("Cookie good %#s\n", cookie)
			return true
		}
		fmt.Printf("Bad:\n")
		fmt.Printf("Cookie: %#s\n", cookie)
		return false
	}
	fmt.Printf("No Cookie %s\n", err)
	return false
}

func SetCookieHandler(w http.ResponseWriter, r *http.Request, id string) {
	value := make(map[string]string)
	//value := map[string]string{
	//	"id": id,
	//}

	if encoded, err := sc.Encode(cookieName, value); err == nil {
		fmt.Printf("encoded: %s\n", string(encoded))
		cookie := &http.Cookie{
			Name:  cookieName,
			Value: "temp",
			Path:  "/",
			Secure: true,
		}
		http.SetCookie(w, cookie)
		//w.Header().Set( "Set-Cookie", "testcookie=temp" )
		fmt.Fprint(w, cookie)
		fmt.Printf("Set the cookie\n")
		fmt.Printf("Cookie: %#s\n", cookie)
		//if err != nil {
		//	panic("Cookie error")
		//}
	}
}

func identifyUser(w http.ResponseWriter, username string, password string) (db.User, error) {
	testCTX := make( map[string]interface{} )
	testCTX["username"] = username

	retCTX, err := db.FindOne(w, "brypt_users", testCTX)
	if err != nil {
		fmt.Println(err)
	}

	du := retCTX["ret"].(db.User)
	fmt.Println(du)

	fmt.Printf("%+v\n", du.Uid)
	if err != nil {
		print("\nError finding user:")
		fmt.Println(err)
		return du, err
	} else {
		print("\nDoing plainTextPW\n")
		plainTextPW := []byte(fmt.Sprintf("%s%s", "salt", password))
		//plainTextPW := []byte(fmt.Sprintf("%s%s", du.Password_salt, password))

		print("Doing GenerateFromPassword\n")
		hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPW), 0)
		if err != nil {
			print("Error doing GenerateFromPassword\n")
			fmt.Println(err)
			return du, err
		}
		fmt.Println("Hash is: ")
		println(hash)

		err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(plainTextPW))
		if err == nil {
			fmt.Printf("User %s authorized\n", username)
			//fmt.Printf("User %s authorized\n", du.Username)
			return du, err
		} else {
			print("Error doing CompareHashAndPassword: ")
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

	r.Get( "/", rs.Index )
	r.Get( "/login", rs.Login )
	r.Get( "/register", rs.Register )
	r.Get( "/link", checkAuth(rs.Link) )

	//r.Get( "/", rs.Index )	// Implemetation of base access page which will support login and registration actions
	//r.Post( "/login", rs.Login )		// Post request for user login
	//r.Post( "/register", rs.Register )	// Post request for registering an account
	//r.Post( "/link", rs.Link )	// Post request for linking a device to a user account

	return r
}

func (rs Resources) Index2(w http.ResponseWriter, r *http.Request) {

	w.Write( []byte( "Hi" ) )
}

/* **************************************************************************
** Function: Index
** URI: access.host (GET)
** Description: Handles compilation the access index page which displays login
** and registration forms.
** Client: Displays the compiled page/
** *************************************************************************/
func (rs Resources) Index(w http.ResponseWriter, r *http.Request) {

	//TestInsert(w)	// TODO: REMOVE WHEN FINISHED TESTING DB INSERT
	//TestUpdate(w)	// TODO: REMOVE WHEN FINSHED TESTING DB UPDATE, FIX
	//TestDelete(w)	// TODO: REMOVE WHEN FINISHED TESTING DB DELETE
	//TestFind(w)		// TODO: REMOVE WHEN FINISHED TESTING DB FIND, FIX

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
	//username := "m@llory5"
	//password := "pswd"

	//Commenting this out fixes it
	//du, err := identifyUser(w, username, password)
	//fmt.Printf("Uid: %s\n", du.Uid)
	//if err != nil {
	//	//w.Write( []byte( "Could not login...\n" ) )
	//	return
	//}

	// On success, add a cookie
	//SetCookieHandler(w, r, du.Uid)
	cookie := &http.Cookie{
		Name:  cookieName,
		Value: "temp",
		Path:  "/",
		Secure: true,
	}
	//http.SetCookie(w, cookie)
	fmt.Printf("Cookie: %#s\n", cookie)
	w.Header().Set( "Content-Type", "text/html" )
	w.Header().Set( "Set-Cookie", "testcookie=temp; Path=/; Secure" )
	fmt.Println(w.Header())
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
	TestInsert2(w)	// TODO: REMOVE WHEN FINISHED TESTING DB INSERT
	fmt.Println("")
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
	id := db.Write(w, "brypt_users", testCTX)
	print("\nid: ")
	fmt.Print(id)
	return id
	//id = db.Write(w, "brypt_users", testCTX)
	//print("\nid: ")
	//fmt.Print(id)
}


/* **************************************************************************
** Function: TestInsert
** Description: Just a test function to demonstrate db insert functionality
**	TODO: Remove when finished testing db insert
** *************************************************************************/
func TestInsert(w http.ResponseWriter) {
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
	id := db.Write(w, "brypt_usrs", testCTX)	// Incorrect collection name (should return nilObjectID)
	print("\nnil id: ")
	fmt.Print(id)
	id = db.Write(w, "brypt_users", testCTX)
	print("\nid: ")
	fmt.Print(id)
	id = db.Write(w, "brypt_users", testCTX)
	print("\nid: ")
	fmt.Print(id)

/*	testCTX["username"] = "TotallyTom"
	testCTX["first_name"] = "Alice"
	testCTX["last_name"] = "Allen"
	testCTX["region"] = "Wonderland"
	testCTX["age"] = time.Now().Round(time.Millisecond)
	testCTX["login_attempts"] = 4
	testCTX["networks"] = []objectid.ObjectID{objID1, objID2, objID3}
	id2 := db.Write(w, "brypt_users", testCTX)
	print("\nid2: ")
	fmt.Print(id2)*/
//	defer db.Disconnect()	// Causes an internal server error for some reason...
}

func TestDelete(w http.ResponseWriter) {

	testCTX := make( map[string]interface{} )
//	testCTX["username"] = "AwesomeAlice"
//	testCTX["first_name"] = "Alice"
	testCTX["last_name"] = "Allen"
	err := db.DeleteOne(w, "brypt_users", testCTX)
	print("\nDelete One error response: ")
	fmt.Print(err)
//	err = db.DeleteAll(w, "brypt_users", testCTX)
//	print("\nDelete All error response: ")
//	fmt.Print(err)
}

func TestFind2(w http.ResponseWriter) map[string]interface{} {

	testCTX := make( map[string]interface{} )
	testCTX["username"] = "m@llory5"
//	testCTX["username"] = "TotallyTom"
//	testCTX["first_name"] = "Alice"
//	testCTX["last_name"] = "Allen"

	/**********FIND ALL TEST**************/
	retCTX, err := db.FindOne(w, "brypt_users", testCTX)
	
	print("\nFind All results: \n")
	fmt.Printf("%+v\n\n\n", retCTX)
	fmt.Printf("%+v\n", retCTX["ret"])
	
	print("\nFind All error response: ")
	fmt.Println(err)
	return retCTX

}

func TestFind(w http.ResponseWriter) {

	testCTX := make( map[string]interface{} )
	testCTX["username"] = "m@llory5"
//	testCTX["username"] = "TotallyTom"
//	testCTX["first_name"] = "Alice"
//	testCTX["last_name"] = "Allen"

	/**********FIND ALL TEST**************/
	retCTX, err := db.FindAll(w, "brypt_users", testCTX)
	
	print("\nFind All results: \n")
	fmt.Printf("%+v\n", retCTX)
	
	print("\nFind All error response: ")
	fmt.Println(err)

	/**********FIND ONE TEST**************/
	testCTX["username"] = "notInDB"
	retCTX, err = db.FindOne(w, "brypt_users", testCTX)
	
	print("\nFind One result:\n ")
	fmt.Printf("%+v\n", retCTX["ret"])
	
	print("\nFind One error response: ")
	fmt.Println(err)

}

func TestUpdate(w http.ResponseWriter) {
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
	
	err := db.UpdateOne(w, "brypt_users", testCTX, updateCTX)
	print("\nUpdate One response: ")
	fmt.Print(err)
}
