package database

import (

	"reflect"	// For printing types
	"fmt"
	"log"
	"net/http"
//	"context"
//	"encoding/json"
	"time"
	"os"
	config "brypt-server/config"

//	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/mongodb/mongo-go-driver/mongo"
	// "github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/ftdc/bsonx"
	"github.com/mongodb/ftdc/bsonx/objectid"
//	"github.com/mongodb/ftdc/bsonx/bsontype"
	//"github.com/mongodb/ftdc/bsonx/elements"
//	"github.com/mongodb/mongo-go-driver/x/mongo/driver/uuid"
)

var configuration = config.Configuration{}

type key string

/* **************************************************************************
** Managers
** *************************************************************************/
type Manager struct {
	ID                objectid.ObjectID	`bson:"_id,omitempty" json:"_id,omitempty"`
	Manager_name      string            `bson:"manager_name" json:"manager_name"`
}

/* **************************************************************************
** Clusters
** *************************************************************************/
type Cluster struct {
	ID                objectid.ObjectID	`bson:"_id,omitempty" json:"_id,omitempty"`
	Connection_token  string            `bson:"connection_token" json:"connection_token"`
	Coord_ip          string            `bson:"coord_ip" json:"coord_ip"`
	Coord_port        string            `bson:"coord_port" json:"coord_port"`
	Comm_tech         string            `bson:"comm_tech" json:"comm_tech"`
}

/* **************************************************************************
** Networks
** *************************************************************************/
type Network struct {
	ID                objectid.ObjectID	`bson:"_id,omitempty" json:"_id,omitempty"`
	Network_name      string            `bson:"network_name" json:"network_name"`
	Owner_name        string            `bson:"owner_name" json:"owner_name"`
	Managers          []objectid.ObjectID         `bson:"managers" json:"managers"`
	Direct_peers      int32             `bson:"direct_peers" json:"direct_peers"`
	Total_peers       int32            `bson:"total_peers" json:"total_peers"`
	Ip_address        string            `bson:"ip_address" json:"ip_address"`
	Port              int32             `bson:"port" json:"port"`
	Connection_token  string            `bson:"connection_token" json:"connection_token"`
	Clusters          []objectid.ObjectID         `bson:"clusters" json:"clusters"`
	Created_on        time.Time         `bson:"created_on" json:"created_on"`
	Last_accessed     time.Time         `bson:"last_accessed" json:"last_accessed"`
}

/* **************************************************************************
** Users
** *************************************************************************/
type User struct {
	ID                objectid.ObjectID	`bson:"_id,omitempty" json:"_id,omitempty"`
	Username          string            `bson:"username" json:"username"`
	First_name        string            `bson:"first_name" json:"first_name"`
	Last_name         string            `bson:"last_name" json:"last_name"`
	Email             string            `bson:"email" json:"email"`
	Organization      string            `bson:"organization" json:"organization"`
	Networks          []objectid.ObjectID         `bson:"networks" json:"networks"`
	Age               time.Time         `bson:"age" json:"age"`
	Join_date         time.Time         `bson:"join_date" json:"join_date"`
	Last_login        time.Time         `bson:"last_login" json:"last_login"`
	Login_attempts    int32             `bson:"login_attempts" json:"login_attempts"`
	Login_token       string            `bson:"login_token" json:"login_token"`
	Region            string            `bson:"region" json:"region"`
}

/* **************************************************************************
** Nodes
** *************************************************************************/
type Node struct {
	ID                objectid.ObjectID	`bson:"_id,omitempty" json:"_id,omitempty"`
	Serial_number     string            `bson:"serial_number" json:"serial_number"`
	Type              string            `bson:"type" json:"type"`
	Created_on        time.Time         `bson:"created_on" json:"created_on"`
	Registered_on     time.Time         `bson:"registered_on" json:"registered_on"`
	Registered_to     string            `bson:"registered_to" json:"registered_to"`
	Connected_network string            `bson:"connected_network" json:"connected_network"`
}


const (
	hostKey = "123"
)

var (
	Client *mongo.Client
)

/* **************************************************************************
** Function: Setup
** URI:
** Description:
** *************************************************************************/
func Setup() {

	var err error

	configuration = config.GetConfig()

	connectionURL := configuration.Database.MongoURI
	if connectionURL == "" {
		panic( "Connection variable is not set!" )
	}

	Client, err = mongo.NewClient( connectionURL )
	if err != nil {
		panic( err )
	}

	fmt.Print( Client )

}

func ReqHandler(w http.ResponseWriter, r *http.Request, collection string, dataCTX map[string]interface{}) {
	
	print("In users handler!\n")

	sterlizeCTXData(dataCTX)	// TODO: Need to implement this function
		
	switch collection {
		case "users":
				WriteUser(w, dataCTX)
				break
		case "nodes":
				WriteNode(w, dataCTX)
				break
		case "networks":
				WriteNetwork(w, dataCTX)
				break
		case "clusters":
				WriteCluster(w, dataCTX)
				break
		case "managers":
				WriteManager(w, dataCTX)
				break
	}

	/*switch r.Method {
	case "GET":
		users_collection := Client.Database("heroku_ckmt3tbl").Collection("brypt_users")

		var sort *options.FindOptions
		var err error

		//sort, err := users_collection.find({}, {"username":1, _id:0}).sort({"username":1})
		// TODO: implement SetSort
		// sort, err := mongo.Opt.Sort(bsonx.NewDocument(bsonx.EC.Int32("username", 1)))   // Sort by username?

		if err != nil { // Error handler for sort
			log.Fatal("Error in usersHandler() sorting: ", err)
		}

		cursor, err := users_collection.Find(nil, nil, sort)  // Get a cursor to the start of the collection?
		if err != nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatal("Internal server error")
			return
		}

		defer cursor.Close(context.Background())

		var users []User  // Create an array to store data from users table

		for cursor.Next(nil) {  // Iterate through users_collection
			user := User{}
			err := cursor.Decode(&user) // Catch any errors while decoding the user object
			if err != nil { // Log the error caught
				log.Fatal("Users collection decode error: ", err)
			}
			users = append(users, user) // Append the stored object to our users array
		}

		if err := cursor.Err(); err != nil {  // Check for a cursor error
			log.Fatal("Cursor error: ", err)
		}

		jsonstr, err := json.Marshal(users)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatal("Internal server error")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonstr)  // Returns the entire json users collection
		return

	case "PUT":
		r.ParseForm()
		//	users_collection := Client.Database("heroku_ckmt3tbl").Collection("brypt_users")

		// TODO: Perform error checking
		newUser := bsonx.NewDocument(bsonx.EC.String("username", r.Form.Get("username")),
		bsonx.EC.String("first_name", r.Form.Get("first_name")),
		bsonx.EC.String("last_name", r.Form.Get("last_name")))
		WriteUsers(newUser, w)
		return
	}*/
	return
}

func sterlizeCTXData(ctx map[string]interface{}) {
	// TODO: Loop through ctx and check that values don't contain invalid characters
	for k := range ctx {
		print(k)
	}
}

// TODO: Change to createBSONDocument(ctx, keys) and return a document
func getStringValues(ctx map[string]interface{}, keys []string) map[string]interface{} {
	tempCTX := make( map[string]interface{} )

	for i := range keys {	
		tempCTX[keys[i]] = ""	// Initialize all key value pairs to empty
	}

	for k := range ctx {
		for j := range keys {
			if k == keys[j] {	// Store value if k matches a key in the users collection
				tempCTX[keys[j]] = ctx[k]
			}
		}
	}

	return tempCTX
}

func insertValue(ctx map[string]interface{}, key string) *bsonx.Document {
	valStr, okStr := ctx[key].(string)	// Check if the type is a string
	doc := bsonx.NewDocument(bsonx.EC.String("fail", "fail"))	// TODO: Return an error of some sort
	if okStr {	
			doc = bsonx.NewDocument(bsonx.EC.String(key, valStr))
	} else {	// Check if int
			valInt, okInt := ctx[key].(int)
		if okInt {
			valInt32 := int32(valInt)
			doc = bsonx.NewDocument(bsonx.EC.Int32(key, valInt32))
		} else {	// Check if int 32
			valInt32_c, okInt32 := ctx[key].(int32)
			if okInt32 {
				doc = bsonx.NewDocument(bsonx.EC.Int32(key, valInt32_c))
			} else {	// Check if time
				valTime, okTime := ctx[key].(time.Time)
				if okTime {
					doc = bsonx.NewDocument(bsonx.EC.Time(key, valTime))
				}	else {
					valObjID, okObjID := ctx[key].([]objectid.ObjectID)
					if okObjID {
						arr := bsonx.NewArray()
						for i := range valObjID {
							arr.Append(bsonx.VC.ObjectID(valObjID[i]))
						}
						doc = bsonx.NewDocument(bsonx.EC.Array(key, arr))
					}	else {	// Value is not a string, int, int32, or time
						fmt.Print("\nFailed to insert value: ")
						print(key)
						print("\n")
						fmt.Println(reflect.TypeOf(ctx[key]))
						print("\n")
					}
				}
			}
			//doc := bsonx.NewDocument(bsonx.EC.String("fail", "fail"))	// TODO: Return an error of some sort
		}
	}

	print("\ninserted!\n")
	return doc
}

func appendValue(doc *bsonx.Document, ctx map[string]interface{}, key string) {
	valStr, okStr := ctx[key].(string)	// Check if the type is a string
	if okStr {	
		doc.Append(bsonx.EC.String(key, valStr))
	} else {	// Check if int
		valInt, okInt := ctx[key].(int)
		if okInt {
			valInt32 := int32(valInt)
			doc.Append(bsonx.EC.Int32(key, valInt32))
		} else {	// Check if int32
			valInt32_c, okInt32 := ctx[key].(int32)
			if okInt32 {
				doc.Append(bsonx.EC.Int32(key, valInt32_c))
			} else {	// Check if time
				valTime, okTime := ctx[key].(time.Time)
				if okTime {
					doc.Append(bsonx.EC.Time(key, valTime))
				} else {	// Check if array of object ids
					valObjID, okObjID := ctx[key].([]objectid.ObjectID)
					if okObjID {
						arr := bsonx.NewArray()
						for i := range valObjID {	// Build array type *Array of object ids
							arr.Append(bsonx.VC.ObjectID(valObjID[i]))
						}
						doc.Append(bsonx.EC.Array(key, arr))	// Append the object id array to the BSON document
					} else {
						fmt.Print("\nFailed to append value!\n")
					}
				}
			}
		}
	}

	print("\nappended!\n")
}

func createBSONDocument(ctx map[string]interface{}, keys []string) *bsonx.Document {
	firstPass := true	// Used to know when to start appending to the new document
	var NewUser *bsonx.Document
	//	tempCTX := make( map[string]interface{} )

/*	for i := range keys {	
		tempCTX[keys[i]] = ""	// Initialize all key value pairs to empty
	}
*/
	for k := range ctx {
		for j := range keys {
			if k == keys[j] {	// Store value if k matches a key in the users collection
				if firstPass {
					NewUser = insertValue(ctx, keys[j])
					firstPass = false
				} else {
					appendValue(NewUser, ctx, keys[j])
				}
				//	tempCTX[keys[j]] = ctx[k]
			}
		}
	}

//	NewUser = bsonx.NewDocument(bsonx.EC.String("username", tempCTX["username"].(string)))
//	print("\n\nBefore append...\n\n")
//	fmt.Print(NewUser)	
//	NewUser.Append(bsonx.EC.String("last_name", tempCTX["last_name"].(string)))
//	bsonx.NewUser.Append(bsonx.EC.String("last_name", tempCTX["last_name"].(string)))
						//									 bsonx.EC.String("email", stringCTX["email"].(string)))
	print("\n\nAfter append...\n\n")
	fmt.Print(NewUser)
	print("\n\n")
	return NewUser 
}

func WriteUser(w http.ResponseWriter, userCTX map[string]interface{}){
//	users_collection := Client.Database("heroku_ckmt3tbl").Collection("brypt_users")
	var keys = []string {"username","first_name","last_name","email", "region", "login_attempts", "age", "objids"}

	newUser := createBSONDocument(userCTX, keys)
	print("\n\n In Write User...\n\n")
	fmt.Print(newUser)
	/*	tempCTX := getStringValues(userCTX, keys)
	NewUser := bsonx.NewDocument(bsonx.EC.String("username", tempCTX["username"].(string)),															 							 bsonx.EC.String("first_name", tempCTX["first_name"].(string)))

	print("\n\nBefore append...\n\n")
	fmt.Print(NewUser)	
	NewUser.Append(bsonx.EC.String("last_name", tempCTX["last_name"].(string)))
//	bsonx.NewUser.Append(bsonx.EC.String("last_name", tempCTX["last_name"].(string)))
						//									 bsonx.EC.String("email", stringCTX["email"].(string)))
	print("\n\nAfter append...\n\n")
	fmt.Print(NewUser)
	print("\n\n")
	_, err := users_collection.InsertOne(nil, newUser)
	if err != nil {
		log.Println("Error inserting new user: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
*/	return
}

func WriteNetwork(w http.ResponseWriter, networkCTX map[string]interface{}) {
	// TODO
}

func WriteNode(w http.ResponseWriter, nodeCTX map[string]interface{}) {
	// TODO
}

func WriteCluster(w http.ResponseWriter, clusterCTX map[string]interface{}) {
	// TODO
}

func WriteManager(w http.ResponseWriter, managerCTX map[string]interface{}){
	newManager := bsonx.NewDocument(bsonx.EC.String("manager_name", "testname"))

	m_collection := Client.Database("heroku_ckmt3tbl").Collection("brypt_managers")

	_, err := m_collection.InsertOne(nil, newManager)
	if err != nil {
		log.Println("Error inserting new manager: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	return
}

func CreateClient() {

	var err error
	print("\n\nSetting up client...\n")

	configuration = config.GetConfig()

	connectionURL := configuration.Database.MongoURI
	if connectionURL == "" {
		panic( "Connection variable is not set!" )
	}

	cert_path, cert_exists := os.LookupEnv("MONGODB_CERT_PATH")

	if cert_exists {  // If user has certification, create a new client with cert info
		print(cert_path)
		Client, err = mongo.NewClient(connectionURL)
		//Client, err = mongo.NewClientWithOptions(connectionURL, mongo.ClientOpt.SSLCaFile(cert_path))
	} else {  // Else create a new client without cert info
		Client, err = mongo.NewClient(connectionURL)
	}

	if err != nil {
		panic( err )
	}

	print("\nFinished setting up client...\n\n")
	fmt.Print( Client )
	return
}

// /* **************************************************************************
// ** Function: CreateClient
// ** URI:
// ** Description: Creates a database client
// ** *************************************************************************/
//
// func CreateClient() {
//
// 				var err error
//
// 				connection_url, url_exists := os.LookupEnv("COMPOSE_MONGODB_URL")
// 				if !url_exists) {
// 								log.Fatal("COMPOSE_MONGODB_URL environmental variable is not set. This needs to be set to ...")
// 				}
//
// 				cert_path, cert_exists := os.LookupEnv("MONGODB_CERT_PATH")
//
// 				if cert_exists {  // If user has certification, create a new client with cert info
// 								Client, err = mongo.NewClientWithOptions(connection_url, mongo.ClientOpt.SSLCaFile(cert_path))
// 				} else {  // Else create a new client without cert info
// 								Client, err = mongo.NewClient(connection_url)
// 				}
//
// 				if err != nil {
// 								log.Fatal(err)  // Log any errors which come up during client connection
// 				}
//
// }

/* **************************************************************************
** Function: Connect
** URI:
** Description: Creates client connection
** *************************************************************************/
func Connect() {
	var err error

	err = Client.Connect(nil)

	if err != nil {
		log.Fatal(err)  // Log any errors thrown during connection
	}

 }


// /* **************************************************************************
// ** Function: Disconnect
// ** URI:
// ** Description: Disconnects client
// ** *************************************************************************/
// func Disconnect() {
// 	defer Client.Disconnect(nil)	// Disconnection client
// }
