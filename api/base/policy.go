package base

import (
	"net/http"

	"brypt-server/internal/handlebars"
)

// privacy policy
func RenderPolicy(w http.ResponseWriter, r *http.Request) {

	policyCTX := make(map[string]interface{})

	policyCTX["title"] = "Brypt"

	page := handlebars.CompilePage( "policy", policyCTX )

	w.Header().Set( "Content-Type", "text/html" )
	w.Write( []byte( page ) )

}

// terms of service
