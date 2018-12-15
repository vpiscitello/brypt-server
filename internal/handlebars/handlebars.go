package handlebars

import (
    "os"
    // "fmt"
    "net/http"
    "path/filepath"

    "encoding/json"
    "io/ioutil"

    "github.com/aymerick/raymond"
)

var workingDir, _ = os.Getwd()

var layoutPath = filepath.Join( workingDir, "/web/views/layouts/main.hbs" )

var headerPath = filepath.Join( workingDir, "/web/views/partials/header.hbs" )
var footerPath = filepath.Join( workingDir, "/web/views/partials/footer.hbs" )
var qlCardPath = filepath.Join( workingDir, "/web/views/partials/link-card.hbs" )
var tmCardPath = filepath.Join( workingDir, "/web/views/partials/team-card.hbs" )

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

    dat, err = ioutil.ReadFile( qlCardPath )
    if err != nil {
        panic( "Something went wrong reading the qlCard partial!" )
    }
    raymond.RegisterPartial( "qlCard", string(dat) )

    dat, err = ioutil.ReadFile( tmCardPath )
    if err != nil {
        panic( "Something went wrong reading the tmCard partial!" )
    }
    raymond.RegisterPartial( "tmCard", string(dat) )

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

/* **************************************************************************
** Function: CompilePage
** Description: Compiles a provided page source using handlers.
** Client: NA
** *************************************************************************/
func CompilePage(page string, bodyCTX map[string]interface{}) string {

    bodyPath := filepath.Join( workingDir, "web/views/pages/" + page + ".hbs" )

    parseContext( bodyCTX )

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

    source, err := layoutTmpl.Exec( pageCTX )
    if err != nil {
        panic( err )
    }

    return source

}

/* **************************************************************************
** Function: parseContext
** Description: Parses the supplied context and expands any keys if needed.
** Client: NA
** *************************************************************************/
func parseContext(ctx map[string]interface{}) {
    
    for key, data := range ctx {
        switch data.(type) {
            default:
                break
            case *os.File:
                expandFileContext( key, ctx )
        }
    }
    
}

/* **************************************************************************
** Function: expandFileContext
** Description: Loads a file into the context
** Client: NA
** *************************************************************************/
func expandFileContext(key string, ctx map[string]interface{}) {

    file := ctx[key].(*os.File)

    data, err := ioutil.ReadAll( file )
    if err != nil {
        panic( "Unable to expand file!" )
    }

    var arr []interface{}
    err = json.Unmarshal( data, &arr )
    if err != nil {
        panic( "Could not unmarshal data!" )
    }

    _ = file.Close()
    delete( ctx, key )

    ctx[key] = arr

}
