package database

import (

//	"reflect"	// For printing types
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
	Uid               string						`bson:"uid" json:"uid"`
	Manager_name      string            `bson:"manager_name" json:"manager_name"`
}

/* **************************************************************************
** Clusters
** *************************************************************************/
type Cluster struct {
//	ID                objectid.ObjectID	`bson:"_id,omitempty" json:"_id,omitempty"`
	Uid               string						`bson:"uid" json:"uid"`
	Connection_token  string            `bson:"connection_token" json:"connection_token"`
	Coord_ip          string            `bson:"coord_ip" json:"coord_ip"`
	Coord_port        string            `bson:"coord_port" json:"coord_port"`
	Comm_tech         string            `bson:"comm_tech" json:"comm_tech"`
}

/* **************************************************************************
** Networks
** *************************************************************************/
type Network struct {
//	ID                objectid.ObjectID	`bson:"_id,omitempty" json:"_id,omitempty"`
	Uid               string						`bson:"uid" json:"uid"`
	Network_name      string            `bson:"network_name" json:"network_name"`
	Owner_name        string            `bson:"owner_name" json:"owner_name"`
	Managers          []string          `bson:"managers" json:"managers"`
	Direct_peers      int32             `bson:"direct_peers" json:"direct_peers"`
	Total_peers       int32             `bson:"total_peers" json:"total_peers"`
	Ip_address        string            `bson:"ip_address" json:"ip_address"`
	Port              int32             `bson:"port" json:"port"`
	Connection_token  string            `bson:"connection_token" json:"connection_token"`
	Clusters          []string          `bson:"clusters" json:"clusters"`
	Created_on        time.Time         `bson:"created_on" json:"created_on"`
	Last_accessed     time.Time         `bson:"last_accessed" json:"last_accessed"`
}

/* **************************************************************************
** Users
** *************************************************************************/
type User struct {
//	ID                objectid.ObjectID	`bson:"_id,omitempty" json:"_id,omitempty"`
	Uid               string						`bson:"uid" json:"uid"`
	Username          string            `bson:"username" json:"username"`
	First_name        string            `bson:"first_name" json:"first_name"`
	Last_name         string            `bson:"last_name" json:"last_name"`
	Email             string            `bson:"email" json:"email"`
	Organization      string            `bson:"organization" json:"organization"`
	Networks          []string          `bson:"networks" json:"networks"`
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
//	ID                objectid.ObjectID	`bson:"_id,omitempty" json:"_id,omitempty"`
	Uid               string						`bson:"uid" json:"uid"`
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

func Write(w http.ResponseWriter, collection string, dataCTX map[string]interface{}) string {
	
	print("In write request handler!\n")

	sterilizeCTXData(dataCTX)	// TODO: Need to implement this function
		
	var id string

	switch collection {
		case "brypt_users":
				print("Handling request to add new user\n")
				id = WriteUser(w, dataCTX)
				print("New user id: ")
				fmt.Print(id)
				print("\n")
			break
		case "brypt_nodes":
				print("Handling request to add new node\n")
				id = WriteNode(w, dataCTX)
			break
		case "brypt_networks":
				print("Handling request to add new network\n")
				id = WriteNetwork(w, dataCTX)
			break
		case "brypt_clusters":
				print("Handling request to add new cluster\n")
				id = WriteCluster(w, dataCTX)
			break
		case "brypt_managers":
				print("Handling request to add new manager\n")
				id = WriteManager(w, dataCTX)
			break
		default:
				print("ERROR: Invalid POST request\n")
				id = ""
			break
	}
	return id
}

func DeleteAll(w http.ResponseWriter, collection string, dataCTX map[string]interface{}) error {
	sterilizeCTXData(dataCTX)
	err := deleteMany(w, collection, dataCTX)
	return err
}

func DeleteOne(w http.ResponseWriter, collection string, dataCTX map[string]interface{}) error {
	sterilizeCTXData(dataCTX)
	err := deleteOne(w, collection, dataCTX)
	return err
}

func FindAll(w http.ResponseWriter, collection string, dataCTX map[string]interface{}) (map[string]interface{}, error) {
	sterilizeCTXData(dataCTX)
	retCTX, err := getAll(w, collection, dataCTX)
	return retCTX, err
}

func FindOne(w http.ResponseWriter, collection string, dataCTX map[string]interface{}) (map[string]interface{}, error) {
	sterilizeCTXData(dataCTX)
//	res, err := getOne(w, collection, dataCTX)
	retCTX, err := getOne(w, collection, dataCTX)
	return retCTX, err
}

func UpdateOne(w http.ResponseWriter, collection string, dataCTX map[string]interface{}, updateCTX map[string]interface{}) error {
	sterilizeCTXData(dataCTX)
	err := updateOne(w, collection, dataCTX, updateCTX)
	return err
}

/* **************************************************************************
** Function: sterilizeCTXData
** URI:
** Description:
** *************************************************************************/
func sterilizeCTXData(ctx map[string]interface{}) {
	// TODO: Loop through ctx and check that values don't contain invalid characters
	print("\nIn sterilizeCTXData\n")
}

/* **************************************************************************
** Function: insertValue
** URI:
** Description:
** *************************************************************************/
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
					valObjID, okObjID := ctx[key].([]string)
					if okObjID {
						arr := bsonx.NewArray()
						for i := range valObjID {
							arr.Append(bsonx.VC.String(valObjID[i]))
						}
						doc = bsonx.NewDocument(bsonx.EC.Array(key, arr))
					}	else {	// Value is not a string, int, int32, or time
						fmt.Print("\nFailed to insert value: ")
						print(key)
						print("\n")
						//fmt.Println(reflect.TypeOf(ctx[key]))
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

/* **************************************************************************
** Function: appendValue
** URI:
** Description:
** *************************************************************************/
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
				} else {	// Check if array of object ids (strings)
					valObjID, okObjID := ctx[key].([]string)
					if okObjID {
						arr := bsonx.NewArray()
						for i := range valObjID {	// Build array type *Array of object ids
							arr.Append(bsonx.VC.String(valObjID[i]))
						}
						doc.Append(bsonx.EC.Array(key, arr))	// Append the object id array to the BSON document
					} else {
						fmt.Print("\nFailed to append value!\n")
					}
				}
			}
		}
	}
}

/* **************************************************************************
** Function: createBSONDocument
** URI:
** Description:
** Returns: Object ID and BSON Document
** *************************************************************************/
func createBSONDocument(ctx map[string]interface{}, keys []string) (string, *bsonx.Document) {
	var NewDoc *bsonx.Document
	objID := objectid.New().Hex()	// Create and store new object id

	NewDoc = bsonx.NewDocument(bsonx.EC.String("uid", objID))
	
	for k := range ctx {
		for j := range keys {
			if k == keys[j] {	// Store value if k matches a key in the users collection
				appendValue(NewDoc, ctx, keys[j])
			}
		}
	}

	return objID, NewDoc 
}

/* **************************************************************************
** Function: WriteUser
** URI:
** Description:
** *************************************************************************/

//Combine (w, ctx, keys, collectionName)
//Return an error if an error
func WriteUser(w http.ResponseWriter, userCTX map[string]interface{}) string {
//	users_collection := Client.Database("heroku_ckmt3tbl").Collection("brypt_users")
	var keys = []string {"username","first_name","last_name","email", "organization", "networks", "age", "join_date", "last_login", "login_attempts", "login_token", "region"}

	objID, newUser := createBSONDocument(userCTX, keys)
	print("\n\n In Write User...\n\n")
	fmt.Print(newUser)
	
	collection := Client.Database("heroku_ckmt3tbl").Collection("brypt_users")

	_, err := collection.InsertOne(nil, newUser)
	if err != nil {
		log.Println("Error inserting new user: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return objID
	}

	w.WriteHeader(http.StatusAccepted)
	return objID
}

/* **************************************************************************
** Function: WriteNetwork
** URI:
** Description:
** *************************************************************************/
func WriteNetwork(w http.ResponseWriter, networkCTX map[string]interface{}) string {
	var keys = []string {"network_name", "owner_name", "managers", "direct_peers", "total_peers", "ip_address", "port", "connection_token", "clusters", "created_on", "last_accessed"}
	objID, newNetwork := createBSONDocument(networkCTX, keys)
	print("\n\n In Write Network...\n\n")
	fmt.Print(newNetwork)
	
	collection := Client.Database("heroku_ckmt3tbl").Collection("brypt_networks")

	_, err := collection.InsertOne(nil, newNetwork)
	if err != nil {
		log.Println("Error inserting new network: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return objID
	}

	w.WriteHeader(http.StatusAccepted)
	return objID
}

/* **************************************************************************
** Function: WriteNode
** URI:
** Description:
** *************************************************************************/
func WriteNode(w http.ResponseWriter, nodeCTX map[string]interface{}) string {
	var keys = []string {"serial_number", "type", "created_on", "registered_on", "registered_to", "connected_network"}
	objID, newNode := createBSONDocument(nodeCTX, keys)
	print("\n\n In Write Node...\n\n")
	fmt.Print(newNode)
	
	collection := Client.Database("heroku_ckmt3tbl").Collection("brypt_nodes")

	_, err := collection.InsertOne(nil, newNode)
	if err != nil {
		log.Println("Error inserting new node: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return objID
	}

	w.WriteHeader(http.StatusAccepted)
	return objID
}

/* **************************************************************************
** Function: WriteCluster
** URI:
** Description:
** *************************************************************************/
func WriteCluster(w http.ResponseWriter, clusterCTX map[string]interface{}) string {
	var keys = []string {"connection_token", "coord_ip", "coord_port", "comm_tech"}
	objID, newCluster := createBSONDocument(clusterCTX, keys)
	print("\n\n In Write Cluster...\n\n")
	fmt.Print(newCluster)
	
	collection := Client.Database("heroku_ckmt3tbl").Collection("brypt_clusters")

	_, err := collection.InsertOne(nil, newCluster)
	if err != nil {
		log.Println("Error inserting new cluster: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return objID
	}

	w.WriteHeader(http.StatusAccepted)
	return objID
}

/* **************************************************************************
** Function: WriteManager 
** URI:
** Description:
** *************************************************************************/
func WriteManager(w http.ResponseWriter, managerCTX map[string]interface{}) string {
	var keys = []string {"manager_name"}
	objID, newManager := createBSONDocument(managerCTX, keys)
	print("\n\n In Write Manager...\n\n")
	fmt.Print(newManager)
	
	collection := Client.Database("heroku_ckmt3tbl").Collection("brypt_managers")

	_, err := collection.InsertOne(nil, newManager)
	if err != nil {
		log.Println("Error inserting new manager: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return objID
	}

	w.WriteHeader(http.StatusAccepted)
	return objID
}

/* **************************************************************************
** Function: DeleteMany 
** URI:
** Description:
** *************************************************************************/
func deleteMany(w http.ResponseWriter, col string, filterCTX map[string]interface{}) error {
		
	collection := Client.Database("heroku_ckmt3tbl").Collection(col)
	_, err := collection.DeleteMany(nil, filterCTX)

	if err != nil {
		log.Println("Error inserting new manager: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusAccepted)
	return err
}

/* **************************************************************************
** Function: DeleteOne 
** URI:
** Description:
** *************************************************************************/
func deleteOne(w http.ResponseWriter, col string, filterCTX map[string]interface{}) error {

	collection := Client.Database("heroku_ckmt3tbl").Collection(col)
	_, err := collection.DeleteOne(nil, filterCTX)

	if err != nil {
		log.Println("Error inserting new manager: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusAccepted)
	return err
}

/* **************************************************************************
** Function: getAll
** URI:
** Description:
** *************************************************************************/
func getAll(w http.ResponseWriter, col string, filterCTX map[string]interface{}) (map[string]interface{}, error) {

 	retCTX := make( map[string]interface{} )
	
	collection := Client.Database("heroku_ckmt3tbl").Collection(col)
	cursor, err := collection.Find(nil, filterCTX)

	if err != nil {
		log.Println("Error inserting new manager: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	var users []User
	for cursor.Next(nil) {
		user := User{}
		e := cursor.Decode(&user)
		if e != nil {
			print("Decode err...")
			fmt.Println(e)
		}
		users = append(users, user)
	//	fmt.Println(user)
	}

	retCTX["ret"] = users

	w.WriteHeader(http.StatusAccepted)	// TODO: Remove??
	return retCTX, err
}

/* **************************************************************************
** Function: getOne
** URI:
** Description:
** *************************************************************************/
func getOne(w http.ResponseWriter, col string, filterCTX map[string]interface{}) (map[string]interface{}, error) {
 	retCTX := make( map[string]interface{} )
	var err error
	collection := Client.Database("heroku_ckmt3tbl").Collection(col)
	res := collection.FindOne(nil, filterCTX)

	switch col {
		case "brypt_users":
				var user User
				err = res.Decode(&user)
				retCTX["ret"] = user
			break
		case "brypt_nodes":
				var node Node
				err = res.Decode(&node)
				retCTX["ret"] = node
			break
		case "brypt_networks":
				var network Network
				err = res.Decode(&network)
				retCTX["ret"] = network
			break
		case "brypt_managers":
				var mngr Manager
				err = res.Decode(&mngr)
				retCTX["ret"] = mngr
			break
		case "brypt_clusters":
				var cluster Cluster
				err = res.Decode(&cluster)
				retCTX["ret"] = cluster
			break
		default:
				err = nil	// TODO: Change to "Unknown collection"
				retCTX["ret"] = nil
			break
	}

/*	fmt.Printf("%+v\n", res)
	
	err := res.Decode(&usr)
*/
	if err != nil {
		fmt.Println(err)
	}
	
//	fmt.Printf("%+v\n", usr)
//	retCTX["ret"] = usr
	w.WriteHeader(http.StatusAccepted)	// TODO: Remove??

	return retCTX, err
}

/* **************************************************************************
** Function: updateOne
** URI:
** Description:
** *************************************************************************/
func updateOne(w http.ResponseWriter, col string, filterCTX map[string]interface{}, updateCTX map[string]interface{}) error {

	collection := Client.Database("heroku_ckmt3tbl").Collection(col)
	_, err := collection.UpdateOne(nil, filterCTX, updateCTX)

	if err != nil {
		log.Println("Error inserting new manager: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusAccepted)
	return err	// TODO: Fix what's broken >:(
}

/* **************************************************************************
** Function: CreateClient 
** URI:
** Description: Creates and configures a new client
** *************************************************************************/
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

/* **************************************************************************
** Function: Connect
** URI:
** Description: Creates client connection
** *************************************************************************/
func Connect() {
	err := Client.Connect(nil)

	if err != nil {
		log.Fatal(err)  // Log any errors thrown during connect
	}

}

//TODO: Figure out how to disconnect client without causing internal server error
/* **************************************************************************
** Function: Disconnect
** URI:
** Description: Disconnects client
** *************************************************************************/
/*func Disconnect() {
	err := Client.Disconnect(nil)	// Disconnection client
	
	if err != nil {
		log.Fatal(err)  // Log any errors thrown during disconnect
	}
}*/
