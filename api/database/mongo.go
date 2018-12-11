package database

import(
				"github.com/mongodb/mongo-go-driver/mongo"
				"context"
				"encoding/json"
				"fmt"
				"log"
				"net/http"
				"os"
)

type key string

const (
		hostKey =
)

var(
		client *mongo.Client  
)


/* **************************************************************************
** Function: CreateClient
** URI:
** Description: Creates a database client
** *************************************************************************/

func CreateClient() {

				var err error

				connection_url, url_exists := os.LookupEnv("COMPOSE_MONGODB_URL")
				if !url_exists) {
								log.Fatal("COMPOSE_MONGODB_URL environmental variable is not set. This needs to be set to ...")
				}

				cert_path, cert_exists := os.LookupEnv("MONGODB_CERT_PATH")

				if cert_exists {  // If user has certification, create a new client with cert info
								client, err = mongo.NewClientWithOptions(connection_url, mongo.ClientOpt.SSLCaFile(cert_path))
				} else {  // Else create a new client without cert info
								client, err = mongo.NewClient(connection_url)
				}

				if err != nil {
								log.Fatal(err)  // Log any errors which come up during client connection
				}

}

/* **************************************************************************
** Function: Connect
** URI:
** Description: Creates client connection
** *************************************************************************/
func Connect() {
				var err error

				err = client.Connect(nil)

				if err != nil {
								log.Fatal(err)  // Log any errors thrown during connection
				}

}


/* **************************************************************************
** Function: Disconnect
** URI:
** Description: Disconnects client
** *************************************************************************/
func Disconnect() {
	defer client.Disconnect(nil)	// Disconnection client
}






