package handlebars

import (
    "os"
    "net/http"
    "path/filepath"

    "io/ioutil"

    "github.com/aymerick/raymond"
)

var workingDir, _ = os.Getwd()
var headerPath = filepath.Join( workingDir, "/web/views/partials/header.hbs" )
var footerPath = filepath.Join( workingDir, "/web/views/partials/footer.hbs" )
var layoutPath = filepath.Join( workingDir, "/web/views/layouts/main.hbs" )

/* **************************************************************************
** Function: Setup
** Description: Initially registers the partials for use with the webpage rendering
** *************************************************************************/
func Setup() {

    dat, err := ioutil.ReadFile( headerPath )
    if err != nil {
        panic( "Something went wrong reading the header partial!" )
    }
    raymond.RegisterPartial( "header", string(dat) )

    dat, err = ioutil.ReadFile( footerPath )
    if err != nil {
        panic( "Something went wrong reading the footer partial!" )
    }
    raymond.RegisterPartial( "footer", string(dat) )

}

/* **************************************************************************
** Function: RenderPage
** Description: Handles the serving and rendering of the webpages
** Client: Displays a specified page
** *************************************************************************/
func RenderPage(page string, bodyCTX map[string]string) http.HandlerFunc {

    bodyPath := ""
    switch page {
        case "index":
            bodyPath = filepath.Join( workingDir, "/web/views/pages/index.hbs" )
        case "access":
            bodyPath = filepath.Join( workingDir, "/web/views/pages/access.hbs" )
    }

    return http.HandlerFunc( func(w http.ResponseWriter, r *http.Request) {

        bodyTmpl, err := raymond.ParseFile( bodyPath )
        if err != nil {
            panic( "Something went wrong parsing the body!" )
        }

        body, err := bodyTmpl.Exec( bodyCTX )
        if err != nil {
            panic( err )
        }

        pageCTX := map[string]string {
            "title": "Brypt",
            "pagestyle": page,
            "body": body,
        }

        layoutTmpl, err := raymond.ParseFile( layoutPath )
        if err != nil {
            panic( "Something went wrong parsing the full!" )
        }

        page, err := layoutTmpl.Exec( pageCTX )
        if err != nil {
            panic( err )
        }

        w.Header().Set( "Content-Type", "text/html" )
        w.Write( []byte( page ) )
    } )
}

