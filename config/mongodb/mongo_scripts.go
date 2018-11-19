// Packages/imports
package mongo

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


// Create collection called "brypt_users"



// Create a collection called "brypt_nodes"


// Create a collection "brypt_users_networks" for many:many relationship???


// Create a collection called "brypt_networks"


// Add a network to the collection (insert data into collection)


// Remove a network from the collection


// Retrieve a network ID based on the network name (network name must be unique for this to work!)


// Update a network in the collection .... maybe wait on this



