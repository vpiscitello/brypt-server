// Packages/imports
package mongo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

)

//////////////////////////////////////
// Global variables
//////////////////////////////////////
var client *mongo.Client


//////////////////////////////////////
// Create client
//////////////////////////////////////
func createClient() {

	var err error

	connection_url, url_exists := os.LookupEnv("COMPOSE_MONGODB_URL")
	if !url_exists) {
		log.Fatal("COMPOSE_MONGODB_URL environmental variable is not set. This needs to be set to ...")
	}

	cert_path, cert_exists := os.LookupEnv("MONGODB_CERT_PATH")

	if cert_exists {	// If user has certification, create a new client with cert info 
		client, err = mongo.NewClientWithOptions(connection_url, mongo.ClientOpt.SSLCaFile(cert_path))
	} else {  // Else create a new client without cert info
		client, err = mongo.NewClient(connection_url)
	}

	if err != nil {
		log.Fatal(err)	// Log any errors which come up during client connection
	}

}


//////////////////////////////////////
// Create connection
//////////////////////////////////////
func connect() {
	var err error

	err = client.Connect(nil)

	if err != nil {
		log.Fatal(err)	// Log any errors thrown during connection
	}

}



//////////////////////////////////////
// Disconnect
//////////////////////////////////////
defer client.Disconnect(nil)	// Disconnect client


//////////////////////////////////////
// Users
//////////////////////////////////////

type user struct {
	ID								objectid.ObjectID	`bson:"_id,omitempty" json:"_id,omitempty"`
	Username					string						`bson:"username" json:"username"`
	First_name				string						`bson:"first_name" json:"first_name"`
	Last_name					string						`bson:"last_name" json:"last_name"`
	Email							string						`bson:"email" json:"email"`
	Organization			string						`bson:"organization" json:"organization"`
	Networks					[]network					`bson:"networks" json:"networks"`
	Age								date							`bson:"age" json:"age"`							
	Join_date					date							`bson:"join_date" json:"join_date"`
	Last_login				date							`bson:"last_login" json:"last_login"`
	Login_attempts		int								`bson:"login_attempts" json:"login_attempts"`	
	Login_token				string						`bson:"login_token" json:"login_token"`
	Region						string						`bson:"region" json:"region"`
}

//////////////////////////////////////
// Nodes
//////////////////////////////////////

type node struct {
	ID								objectid.ObjectID	`bson:"_id,omitempty" json:"_id,omitempty"`
	Serial_number			string						`bson:"serial_number" json:"serial_number"`
	Type							string						`bson:"type" json:"type"`
	Created_on				date							`bson:"created_on" json:"created_on"`							
	Registered_on			date							`bson:"registered_on" json:"registered_on"`
	Registered_to			string						`bson:"registered_to" json:"registered_to"`
	Connected_network	string						`bson:"connected_network" json:"connected_network"`
}

//////////////////////////////////////
// Networks
//////////////////////////////////////

type network struct {
	ID								objectid.ObjectID	`bson:"_id,omitempty" json:"_id,omitempty"`
	Network_name			string						`bson:"network_name" json:"network_name"`
	Owner_name				string						`bson:"owner_name" json:"owner_name"`
	Managers					[]manager					`bson:"managers" json:"managers"`
	Direct_peers			int								`bson:"direct_peers" json:"direct_peers"`	
	Total_peers				int								`bson:"total_peers" json:"total_peers"`	
	Ip_address				string						`bson:"ip_address" json:"ip_address"`
	Port							int								`bson:"port" json:"port"`	
	Connection_token	string						`bson:"connection_token" json:"connection_token"`
	Clusters					[]cluster					`bson:"clusters" json:"clusters"`
	Created_on				date							`bson:"created_on" json:"created_on"`							
	Last_accessed			date							`bson:"last_accessed" json:"last_accessed"`
}

//////////////////////////////////////
// Managers
//////////////////////////////////////

type manager struct {
	ID								objectid.ObjectID	`bson:"_id,omitempty" json:"_id,omitempty"`
	Manager_name			string						`bson:"manager_name" json:"manager_name"`
}

//////////////////////////////////////
// Clusters
//////////////////////////////////////

type cluster struct {
	ID								objectid.ObjectID	`bson:"_id,omitempty" json:"_id,omitempty"`
	Connection_token	string						`bson:"connection_token" json:"connection_token"`
	Coord_ip					string						`bson:"coord_ip" json:"coord_ip"`
	Coord_port				string						`bson:"coord_port" json:"coord_port"`
	Comm_tech					string						`bson:"comm_tech" json:"comm_tech"`
}



//////////////////////////////////////////
// Create collection called "brypt_users"
//////////////////////////////////////////

users_collection := client.Database("brypt_server").Collection("brypt_users")

//////////////////////////////////////////
// Create collection called "brypt_nodes"
//////////////////////////////////////////

nodes_collection := client.Database("brypt_server").Collection("brypt_nodes")

//////////////////////////////////////////
// Create collection called "brypt_networks"
//////////////////////////////////////////

networks_collection := client.Database("brypt_server").Collection("brypt_networks")

//////////////////////////////////////////
// Handle requests for users database
//////////////////////////////////////////

usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "GET":
			users_collection := client.Database("brypt_server").Collection("brypt_users")
			// TODO: Implement! See: http://github.com/compose-grandtour/golang/blob/master/example-mongodb-go-driver/example-mongodb.go/
			
			sort, err := mongo.Opt.Sort(bson.NewDocument(bson.EC.Int32("username", 1)))		// Sort by username?

			if err != nil {	// Error handler for sort
				log.Fatal("Error in usersHandler() sorting: ", err)
			}

			cursor, err := users_collection.Find(nil, nil, sort)	// Get a cursor to the start of the collection?
			if err != nil{
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Fatal("Internal server error")	
				return
			}

			defer cursor.Close(context.Background())

			var users []user	// Create an array to store data from users table

			for cursor.Next(nil) {	// Iterate through users_collection
				user := user{}
				err := cursor.Decode(&user)	// Catch any errors while decoding the user object
				if err != nil {	// Log the error caught
					log.Fatal("Users collection decode error: ", err)
				}
				users = append(users, user)	// Append the stored object to our users array
			}

			if err := cursor.Err(); err != nil {	// Check for a cursor error
				log.Fatal("Cursor error: ", err)
			}

			jsonstr, err := json.Marshal(users)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Fatal("Internal server error")	
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonstr)	// Returns the entire json users collection
			return

		case "PUT":
			r.ParseForm()
			users_collection := client.Database("brypt_server").Collection("brypt_users")

			newUser := bson.NewDocument(bson.EC.String("username", r.Form.Get("username")),
																	bson.EC.String("first_name", r.Form.Get("first_name")),
																	bson.EC.String("last_name", r.Form.Get("last_name")))

			_, err := users_collection.InsertOne(nil, newUser)
			if err != nil {
							log.Println("Error inserting new user: ", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusAccepted)
			return
	}
	return
}

//////////////////////////////////////////
// Handle requests for nodes database
//////////////////////////////////////////

nodesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "GET":
			nodes_collection := client.Database("brypt_server").Collection("brypt_nodes")
			// TODO: Implement! See: http://github.com/compose-grandtour/golang/blob/master/example-mongodb-go-driver/example-mongodb.go/
			
			sort, err := mongo.Opt.Sort(bson.NewDocument(bson.EC.Int32("serial_number", 1)))		// Sort by username?

			if err != nil {	// Error handler for sort
				log.Fatal("Error in nodesHandler() sorting: ", err)
			}

			return

		case "PUT":
			r.ParseForm()
			nodes_collection := client.Database("brypt_server").Collection("brypt_nodes")

			// TODO: Implement! See: http://github.com/compose-grandtour/golang/blob/master/example-mongodb-go-driver/example-mongodb.go/
			return
	}
	return
}

//////////////////////////////////////////
// Handle requests for networks database
//////////////////////////////////////////

networksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "GET":
			networks_collection := client.Database("brypt_server").Collection("brypt_networks")
			// TODO: Implement! See: http://github.com/compose-grandtour/golang/blob/master/example-mongodb-go-driver/example-mongodb.go/
			
			sort, err := mongo.Opt.Sort(bson.NewDocument(bson.EC.Int32("network_name", 1)))		// Sort by username?

			if err != nil {	// Error handler for sort
				log.Fatal("Error in networksHandler() sorting: ", err)
			}

			return

		case "PUT":
			r.ParseForm()
			networks_collection := client.Database("brypt_server").Collection("brypt_networks")

			// TODO: Implement! See: http://github.com/compose-grandtour/golang/blob/master/example-mongodb-go-driver/example-mongodb.go/
			return
	}
	return
}


//////////////////////////////////////////
// Users Handler Setup
//////////////////////////////////////////

fs := http.FileServer(http.Dir("public"))
http.Handle("/", fs)
http.HandleFunc("/users", usersHandler)	// TODO: Not sure what path to look for
fmt.Println("Listening on localhost:8080")	// TODO: Probably change port number??
http.ListenAndServe(":8080", nil)


//////////////////////////////////////////
// Nodes Handler Setup
//////////////////////////////////////////

fs := http.FileServer(http.Dir("public"))
http.Handle("/", fs)
http.HandleFunc("/access", nodesHandler)	// TODO: Not sure what path to look for
fmt.Println("Listening on localhost:8080")	// TODO: Probably change port number??
http.ListenAndServe(":8080", nil)


//////////////////////////////////////////
// Networks Handler Setup
//////////////////////////////////////////

fs := http.FileServer(http.Dir("public"))
http.Handle("/", fs)
http.HandleFunc("/bridge", networksHandler)	// TODO: Not sure what path to look for
fmt.Println("Listening on localhost:8080")	// TODO: Probably change port number??
http.ListenAndServe(":8080", nil)








