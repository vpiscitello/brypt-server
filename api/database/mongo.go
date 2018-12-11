package database

import (

    "fmt"

    config "brypt-server/config"

    "github.com/mongodb/mongo-go-driver/mongo"
    
)

var configuration = config.Configuration{}

type key string

/* **************************************************************************
** Users
** *************************************************************************/
type User struct {
				ID                objectid.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
				Username          string            `bson:"username" json:"username"`
				First_name        string            `bson:"first_name" json:"first_name"`
				Last_name         string            `bson:"last_name" json:"last_name"`
				Email             string            `bson:"email" json:"email"`
				Organization      string            `bson:"organization" json:"organization"`
				Networks          []network         `bson:"networks" json:"networks"`
				Age               date              `bson:"age" json:"age"`
				Join_date         date              `bson:"join_date" json:"join_date"`
				Last_login        date              `bson:"last_login" json:"last_login"`
				Login_attempts    int               `bson:"login_attempts" json:"login_attempts"`
				Login_token       string            `bson:"login_token" json:"login_token"`
				Region            string            `bson:"region" json:"region"`
}

/* **************************************************************************
** Nodes
** *************************************************************************/
type Node struct {
				ID                objectid.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
				Serial_number     string            `bson:"serial_number" json:"serial_number"`
				Type              string            `bson:"type" json:"type"`
				Created_on        date              `bson:"created_on" json:"created_on"`
				Registered_on     date              `bson:"registered_on" json:"registered_on"`
				Registered_to     string            `bson:"registered_to" json:"registered_to"`
				Connected_network string            `bson:"connected_network" json:"connected_network"`
}

/* **************************************************************************
** Networks
** *************************************************************************/
type Network struct {
				ID                objectid.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
				Network_name      string            `bson:"network_name" json:"network_name"`
				Owner_name        string            `bson:"owner_name" json:"owner_name"`
				Managers          []manager         `bson:"managers" json:"managers"`
				Direct_peers      int               `bson:"direct_peers" json:"direct_peers"`
				Total_peers       int               `bson:"total_peers" json:"total_peers"`
				Ip_address        string            `bson:"ip_address" json:"ip_address"`
				Port              int               `bson:"port" json:"port"`
				Connection_token  string            `bson:"connection_token" json:"connection_token"`
				Clusters          []cluster         `bson:"clusters" json:"clusters"`
				Created_on        date              `bson:"created_on" json:"created_on"`
				Last_accessed     date              `bson:"last_accessed" json:"last_accessed"`
}

/* **************************************************************************
** Managers
** *************************************************************************/
type Manager struct {
				ID                objectid.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
				Manager_name      string            `bson:"manager_name" json:"manager_name"`
}

/* **************************************************************************
** Clusters
** *************************************************************************/
type Cluster struct {
				ID                objectid.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
				Connection_token  string            `bson:"connection_token" json:"connection_token"`
				Coord_ip          string            `bson:"coord_ip" json:"coord_ip"`
				Coord_port        string            `bson:"coord_port" json:"coord_port"`
				Comm_tech         string            `bson:"comm_tech" json:"comm_tech"`
}



const (
				hostKey =
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

func WriteUsers(struct User *u){


}

func WriteNetworks(struct Network *n) {
				// TODO
}

func WriteNodes(struct Node *n) {
				// TODO
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
// 								client, err = mongo.NewClientWithOptions(connection_url, mongo.ClientOpt.SSLCaFile(cert_path))
// 				} else {  // Else create a new client without cert info
// 								client, err = mongo.NewClient(connection_url)
// 				}
//
// 				if err != nil {
// 								log.Fatal(err)  // Log any errors which come up during client connection
// 				}
//
// }
//
// /* **************************************************************************
// ** Function: Connect
// ** URI:
// ** Description: Creates client connection
// ** *************************************************************************/
// func Connect() {
// 				var err error
//
// 				err = client.Connect(nil)
//
// 				if err != nil {
// 								log.Fatal(err)  // Log any errors thrown during connection
// 				}
//
// }
//
//
// /* **************************************************************************
// ** Function: Disconnect
// ** URI:
// ** Description: Disconnects client
// ** *************************************************************************/
// func Disconnect() {
// 	defer client.Disconnect(nil)	// Disconnection client
// }
