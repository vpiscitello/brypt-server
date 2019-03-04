package database

import (
	"reflect"
	"fmt"
	"log"
	"strings"
	"time"
	"os"
	config "brypt-server/config"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/ftdc/bsonx"
	"github.com/mongodb/ftdc/bsonx/objectid"
)

var configuration = config.Configuration{}

type key string

/* **************************************************************************
** Managers
** *************************************************************************/
type Manager struct {
	Uid				string			`bson:"uid" json:"uid"`
	ManagerName		string			`bson:"manager_name" json:"manager_name"`
}

/* **************************************************************************
** Clusters
** *************************************************************************/
type Cluster struct {
	Uid					string			`bson:"uid" json:"uid"`
	ConnectionToken		string			`bson:"connection_token" json:"connection_token"`
	CoordIP				string			`bson:"coord_ip" json:"coord_ip"`
	CoordPort			string			`bson:"coord_port" json:"coord_port"`
	CommTech			string			`bson:"comm_tech" json:"comm_tech"`
}

/* **************************************************************************
** Networks
** *************************************************************************/
type Network struct {
	Uid					string			`bson:"uid" json:"uid"`
	NetworkName			string			`bson:"network_name" json:"network_name"`
	OwnerName			string			`bson:"owner_name" json:"owner_name"`
	Managers			[]string		`bson:"managers" json:"managers"`
	DirectPeers			int32			`bson:"direct_peers" json:"direct_peers"`
	TotalPeers			int32			`bson:"total_peers" json:"total_peers"`
	RootAP				string			`bson:"root_ap" json:"root_ap"`
	RootIP				string			`bson:"root_ip" json:"root_ip"`
	RootPort			int32			`bson:"root_port" json:"root_port"`
	ConnectionToken		string			`bson:"connection_token" json:"connection_token"`
	Clusters			[]string		`bson:"clusters" json:"clusters"`
	CreatedOn			time.Time		`bson:"created_on" json:"created_on"`
	LastAccessed		time.Time		`bson:"last_accessed" json:"last_accessed"`
}

/* **************************************************************************
** Users
** *************************************************************************/
type User struct {
	Uid				string			`bson:"uid" json:"uid"`
	Username		string			`bson:"username" json:"username"`
	FirstName		string			`bson:"first_name" json:"first_name"`
	LastName		string			`bson:"last_name" json:"last_name"`
	Email			string			`bson:"email" json:"email"`
	Organization	string			`bson:"organization" json:"organization"`
	Networks[]		string			`bson:"networks" json:"networks"`
	Birthdate		time.Time		`bson:"birthdate" json:"birthdate"`
	JoinDate		time.Time		`bson:"join_date" json:"join_date"`
	LastLogin		time.Time		`bson:"last_login" json:"last_login"`
	LoginAttempts	int32			`bson:"login_attempts" json:"login_attempts"`
	LoginToken		string			`bson:"login_token" json:"login_token"`
	Region			string			`bson:"region" json:"region"`
	Password		string			`bson:"password" json:"password"`
}

/* **************************************************************************
** Nodes
** *************************************************************************/
type Node struct {
	Uid						string			  `bson:"uid" json:"uid"`
	SerialNumber			string            `bson:"serial_number" json:"serial_number"`
	Type					string            `bson:"type" json:"type"`
	RegisteredTimestamp		time.Time         `bson:"registered_timestamp" json:"registered_timestamp"`
	RegisteredTo			string            `bson:"registered_to" json:"registered_to"`
	ConnectedNetwork		string            `bson:"connected_network" json:"connected_network"`
	DataToken				string            `bson:"data_token" json:"data_token"`
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

func Write(collection string, dataCTX map[string]interface{}) string {

	print("In write request handler!\n")

	sterilizeCTXData(dataCTX)	// TODO: Need to implement this function

	var id string

	switch collection {
		case "brypt_users":
				print("Handling request to add new user\n")
				//id = WriteUser(w, dataCTX)
				user := User{}
				id = writeObject(dataCTX, user, collection)
				print("New user id: ")
				fmt.Print(id)
				print("\n")
			break
		case "brypt_nodes":
				print("Handling request to add new node\n")
				//id = WriteNode(w, dataCTX)
				node := Node{}
				id = writeObject(dataCTX, node, collection)
			break
		case "brypt_networks":
				print("Handling request to add new network\n")
			//	id = WriteNetwork(w, dataCTX)
				network := Network{}
				id = writeObject(dataCTX, network, collection)
			break
		case "brypt_clusters":
				print("Handling request to add new cluster\n")
			//	id = WriteCluster(w, dataCTX)
			cluster := Cluster{}
				id = writeObject(dataCTX, cluster, collection)
			break
		case "brypt_managers":
				print("Handling request to add new manager\n")
				//id = WriteManager(w, dataCTX)
				mngr := Manager{}
				id = writeObject(dataCTX, mngr, collection)
			break
		default:
				print("ERROR: Invalid POST request\n")
				id = ""
			break
	}
	return id
}

func DeleteAll(collection string, dataCTX map[string]interface{}) error {
	sterilizeCTXData(dataCTX)
	err := deleteMany(collection, dataCTX)
	return err
}

func DeleteOne(collection string, dataCTX map[string]interface{}) error {
	sterilizeCTXData(dataCTX)
	err := deleteOne(collection, dataCTX)
	return err
}

func FindAll(collection string, dataCTX map[string]interface{}) (map[string]interface{}, error) {
	sterilizeCTXData(dataCTX)
	retCTX, err := getAll(collection, dataCTX)
	return retCTX, err
}

func FindOne(collection string, dataCTX map[string]interface{}) (map[string]interface{}, error) {
	sterilizeCTXData(dataCTX)
//	res, err := getOne(w, collection, dataCTX)
	retCTX, err := getOne(collection, dataCTX)
	return retCTX, err
}

func UpdateOne(collection string, dataCTX map[string]interface{}, updateCTX map[string]interface{}) error {
	sterilizeCTXData(dataCTX)
	err := updateOne(collection, dataCTX, updateCTX)
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
** Function: WriteObject
** URI:
** Description:
** *************************************************************************/
func writeObject(objCTX map[string]interface{}, obj interface{}, collectionName string) string {
//	users_collection := Client.Database("heroku_ckmt3tbl").Collection("brypt_users")
//	var keys = []string {"username","first_name","last_name","email", "organization", "networks", "birthdate", "join_date", "last_login", "login_attempts", "login_token", "region"}
	var keys []string
	var k string
	user := User{}
	fmt.Printf("%v\n", obj)
	elem := reflect.ValueOf(&user).Elem()

	for i := 0; i < elem.NumField(); i++ {
		//keys[i] = elem.Type().Field(i).Name
		k = elem.Type().Field(i).Name
		k = strings.ToLower(k)
		if k != "uid" {
			keys = append(keys, k)
			fmt.Printf("%v\n", k)
		}
	}

	objID, newObj := createBSONDocument(objCTX, keys)
	print("\n\n In WriteObject...\n\n")
	fmt.Print(newObj)

	collection := Client.Database("heroku_ckmt3tbl").Collection(collectionName)

	_, err := collection.InsertOne(nil, newObj)
	if err != nil {
		log.Println("Error inserting new object: ", err)
		return objID
	}

	return objID
}

/* **************************************************************************
** Function: DeleteMany
** URI:
** Description:
** *************************************************************************/
func deleteMany(col string, filterCTX map[string]interface{}) error {

	collection := Client.Database("heroku_ckmt3tbl").Collection(col)
	_, err := collection.DeleteMany(nil, filterCTX)

	if err != nil {
		log.Println("Error inserting new manager: ", err)
		return err
	}

	return err
}

/* **************************************************************************
** Function: DeleteOne
** URI:
** Description:
** *************************************************************************/
func deleteOne(col string, filterCTX map[string]interface{}) error {

	collection := Client.Database("heroku_ckmt3tbl").Collection(col)
	_, err := collection.DeleteOne(nil, filterCTX)

	if err != nil {
		log.Println("Error inserting new manager: ", err)
		return err
	}

	return err
}

/* **************************************************************************
** Function: getAll
** URI:
** Description:
** *************************************************************************/
func getAll(col string, filterCTX map[string]interface{}) (map[string]interface{}, error) {

 	retCTX := make( map[string]interface{} )

	collection := Client.Database("heroku_ckmt3tbl").Collection(col)
	cursor, err := collection.Find(nil, filterCTX)
	if err != nil {
		log.Println("Error finding all: ", err)
		return nil, err
	}

	switch col {
		case "brypt_users":
				var users []User
				for cursor.Next(nil) {
					user := User{}
					e := cursor.Decode(&user)
					if e != nil {
						print("Decode err...")
						fmt.Println(e)
					}
					users = append(users, user)
				}
				retCTX["ret"] = users
			break
		case "brypt_nodes":
				var nodes []Node
				for cursor.Next(nil) {
					node := Node{}
					e := cursor.Decode(&node)
					if e != nil {
						print("Decode err...")
						fmt.Println(e)
					}
					nodes = append(nodes, node)
				}
				retCTX["ret"] = nodes
			break
		case "brypt_networks":
				var networks []Network
				for cursor.Next(nil) {
					network := Network{}
					e := cursor.Decode(&network)
					if e != nil {
						print("Decode err...")
						fmt.Println(e)
					}
					networks = append(networks, network)
				}
				retCTX["ret"] = networks
			break
		case "brypt_managers":
				var mngrs []Manager
				for cursor.Next(nil) {
					mngr := Manager{}
					e := cursor.Decode(&mngr)
					if e != nil {
						print("Decode err...")
						fmt.Println(e)
					}
					mngrs = append(mngrs, mngr)
				}
				retCTX["ret"] = mngrs
			break
		case "brypt_clusters":
				var clusters []Cluster
				for cursor.Next(nil) {
					cluster := Cluster{}
					e := cursor.Decode(&cluster)
					if e != nil {
						print("Decode err...")
						fmt.Println(e)
					}
					clusters = append(clusters, cluster)
				}
				retCTX["ret"] = clusters
			break
		default:
				err = nil	// TODO: Change to "Unknown collection"
				retCTX["ret"] = nil
			break
	}

	return retCTX, err
}

/* **************************************************************************
** Function: getOne
** URI:
** Description:
** *************************************************************************/
func getOne(col string, filterCTX map[string]interface{}) (map[string]interface{}, error) {
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

	if err != nil {
		fmt.Println(err)
	}

	return retCTX, err
}

/* **************************************************************************
** Function: updateOne
** URI:
** Description:
** *************************************************************************/
func updateOne(col string, filterCTX map[string]interface{}, updateCTX map[string]interface{}) error {

	collection := Client.Database("heroku_ckmt3tbl").Collection(col)
	_, err := collection.UpdateOne(nil, filterCTX, updateCTX)

	if err != nil {
		log.Println("Error inserting new manager: ", err)
		return err
	}

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
