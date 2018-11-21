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

// Switch to database (or create it)




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
	Networks					array							`bson:"networks" json:"networks"`
	Age								date							`bson:"age" json:"age"`							
	Join_date					date							`bson:"join_date" json:"join_date"`
	Last_login				date							`bson:"last_login" json:"last_login"`
	Login_attempts		int								`bson:"login_attempts" json:"login_attempts"`	
	Login_token				string						`bson:"login_token" json:"login_token"`
	Region						string						`bson:"region" json:"region"`
}


// Create collection called "brypt_users"

c := client.Database("brypt_server").Collection("brypt_users")

// Users data handler function

usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "GET":
			c := client.Database("brypt_server").Collection("brypt_users")
			// TODO: Implement! See: http://github.com/compose-grandtour/golang/blob/master/example-mongodb-go-driver/example-mongodb.go/
			
			sort, err := mongo.Opt.Sort(bson.NewDocument(bson.EC.Int32("username", 1)))		// Sort by username?

			if err != nil {	// Error handler for sort
				log.Fatal("Error in usersHandler() sorting: ", err)
			}

			return

		case "PUT":
			r.ParseForm()
			c := client.Database("brypt_server").Collection("brypt_users")

			// TODO: Implement! See: http://github.com/compose-grandtour/golang/blob/master/example-mongodb-go-driver/example-mongodb.go/
			return
	}
	return
}

// Handler Setup

fs := http.FileServer(http.Dir("public"))
http.Handle("/", fs)
http.HandleFunc("/access", usersHandler)	// TODO: Not sure what path to look for
http.HandleFunc("/bridge", usersHandler)	// TODO: Not sure what path to look for
http.HandleFunc("/users", usersHandler)	// TODO: Not sure what path to look for
fmt.Println("Listening on localhost:8080")	// TODO: Probably change port number??
http.ListenAndServe(":8080", nil)





//////////////////////////////////////
// Nodes
//////////////////////////////////////

// Create a collection called "brypt_nodes"
c := client.Database("brypt_server").Collection("brypt_nodes")



// TODO: Handler Setup



//////////////////////////////////////
// Networks
//////////////////////////////////////

// Create a collection called "brypt_networks"
c := client.Database("brypt_server").Collection("brypt_networks")

// TODO: Handler Setup


// Add a network to the collection (insert data into collection)

// Remove a network from the collection

// Retrieve a network ID based on the network name (network name must be unique for this to work!)

// Update a network in the collection .... maybe wait on this



