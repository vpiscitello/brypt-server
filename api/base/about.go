package base

import (
	"os"
	"net/http"
	"path/filepath"

	"brypt-server/internal/handlebars"
)

var workingDir, _ = os.Getwd()

// about Brypt
func RenderAbout(w http.ResponseWriter, r *http.Request) {

	tmDataPath := filepath.Join( workingDir, "/web/data/team.json" )

	aboutCTX := make(map[string]interface{})

	aboutCTX["title"] = "Brypt"
	aboutCTX["teamMember"], _ = os.Open( tmDataPath )

	page := handlebars.CompilePage( "about", aboutCTX )

	w.Header().Set( "Content-Type", "text/html" )
	w.Write( []byte( page ) )

}

// contact
func RenderContact(w http.ResponseWriter, r *http.Request) {

	contactCTX := make(map[string]interface{})

	contactCTX["title"] = "Brypt"

	page := handlebars.CompilePage( "contact", contactCTX )

	w.Header().Set( "Content-Type", "text/html" )
	w.Write( []byte( page ) )

}
