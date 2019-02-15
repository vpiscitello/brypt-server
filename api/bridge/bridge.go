package bridge

import (
	"fmt"
	db "brypt-server/api/database"
	"net/http"
	// "time"

    "brypt-server/internal/handlebars"
	"encoding/json"
	// "io/ioutil"

	"github.com/go-chi/chi"
	"github.com/mongodb/mongo-go-driver/bson"

    "brypt-server/api/access"
)

type Resources struct{}

/* **************************************************************************
** Function: Routes
** Description: Register the birdge resources with chi router and returns the
** built router.
** *************************************************************************/
func (rs Resources) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get( "/", access.CheckAuth( rs.Index ) )	// Implemetation of base bridge page
	r.Post( "/node", access.CheckAuth( rs.RegisterNode ) )	// Register a node in the user's network
	r.Get( "/network", access.CheckAuth( rs.GetNodes ) )	// Get the nodes within a user's network

	return r
}

/* **************************************************************************
** Function: Index
** URI: birdge.host (GET)
** Description: Handles compilation the access birdge page
** Client: Displays the compiled page/
** *************************************************************************/
func (rs Resources) Index(w http.ResponseWriter, r *http.Request) {

	bridgeCTX := make( map[string]interface{} )

    bridgeCTX["title"] = "Brypt"

	page := handlebars.CompilePage( "bridge", bridgeCTX )

	w.Header().Set( "Content-Type", "text/html" )
	w.Write( []byte( page ) )

}

/* **************************************************************************
** Function: RegisterNode
** URI: birdge.host (Post)
** Description: Handles registering a node within a user's network
** Client: Provides node registration data
** *************************************************************************/
func (rs Resources) RegisterNode(w http.ResponseWriter, r *http.Request) {

	w.Header().Set( "Content-Type", "text/html" )
	w.Write( []byte( "Registered Node" ) )

}

/* **************************************************************************
** Function: GetNodes
** URI: birdge.host (GET)
** Description: Handles getting the nodes within a user's network
** Client: Handles returned JSON data
** *************************************************************************/
func (rs Resources) GetNodes(w http.ResponseWriter, r *http.Request) {

	// Parse which network based on user cookie

	networkSearchCTX := make( map[string]interface{} )
	networkSearchCTX["managers"] = bson.D{{"$all", bson.A{"5c60b34fe25f5a42f00c4569"}}}

	networkObject := db.Network{}
	// Find user's network based on their user uid
	networkRet, err := db.FindOne("brypt_networks", networkSearchCTX)

	networkObject = networkRet["ret"].(db.Network)

	if err != nil {
		fmt.Println(err)
		w.Header().Set( "Content-Type", "text/html" )
		w.Write( []byte( "Error Occured" ) )
	} else {
		networkObject = networkRet["ret"].(db.Network)
		fmt.Println(networkObject)

		// Find all the nodes within that network
		nodesSearchCTX := make( map[string]interface{} )
		nodesSearchCTX["network"] = networkObject.Uid

		nodesObject := db.Node{}
		// Find user's network based on their user uid
		nodesRet, err := db.FindAll("brypt_nodes", nodesSearchCTX)

		if err != nil {
			fmt.Println(err)
			w.Header().Set( "Content-Type", "text/html" )
			w.Write( []byte( "Error Occured" ) )
		} else {
			nodesObject = nodesRet["ret"].(db.Node)
			fmt.Println(nodesObject)

			nodesJSON, err := json.Marshal(nodesObject)
			if err != nil {
			  http.Error(w, err.Error(), http.StatusInternalServerError)
			  return
			}

			w.Header().Set( "Content-Type", "application/json" )
			w.Write( nodesJSON )

		}
	}

}
