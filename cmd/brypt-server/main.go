package main

// Register Routes
// Run Server

import (
   "fmt"
   "os"
   "net/http"
   "net/url"
   "path/filepath"
   "strings"
   "strconv"

   config "brypt-server/config"

   db "brypt-server/api/database"

   "brypt-server/api/access"
   "brypt-server/api/base"
   // "brypt-server/api/bridge"
   "brypt-server/api/dashboard"
   // "brypt-server/api/users"

   "brypt-server/internal/handlebars"

   "github.com/go-chi/chi"
   "github.com/go-chi/chi/middleware"

   "github.com/go-chi/hostrouter"

)

var workingDir, _ = os.Getwd()

var configuration = config.Configuration{}

/* **************************************************************************
** Function:
** URI:
** Description:
** Client:
** *************************************************************************/
func redirectToHTTPS( w http.ResponseWriter, r *http.Request )  {
    // Build the HTTPS target URL using URL builder
    target := url.URL{
        Scheme: "https",
        Host: strings.Split(r.Host, ":")[0] + ":" + strconv.Itoa( configuration.Server.HTTPSPort ), // Pop off the HTTP port and add the proper HTTPS port
        Path: r.URL.Path,
        RawQuery: r.URL.RawQuery,
    }
    http.Redirect( w, r, target.String(), http.StatusTemporaryRedirect )    // Redirect requests to the HTTPS equiv.
}

/* **************************************************************************
** Function:
** URI:
** Description:
** Client:
** *************************************************************************/
func main()  {
    fmt.Println( "Startup...\n" )

    config.Setup()  // Setup the Server Configuration
    configuration = config.GetConfig()  // Get the Configuration Settings

    fmt.Println( "Got Configuration\n" )

    db.Setup()
		db.Connect()
		handlebars.Setup()

    fmt.Println( "DB and Handlebars Ready\n" )

    HTTPPortString := strconv.Itoa( configuration.Server.HTTPPort )
    HTTPSPortString := strconv.Itoa( configuration.Server.HTTPSPort )

    go http.ListenAndServe( ":" + HTTPPortString, http.HandlerFunc( redirectToHTTPS ) )  // Start the Server

    fmt.Println( "Listening for HTTP\n" )

    router := chi.NewRouter()

    router.Use( middleware.RequestID )
    router.Use( middleware.RealIP )
    router.Use( middleware.Logger )
    router.Use( middleware.Recoverer )

    hr := hostrouter.New()

    hr.Map( configuration.Server.AccessDomain, access.Resources{}.Routes() ) // Handle access.host routing requests

    hr.Map( configuration.Server.BridgeDomain, bridgeRouter() ) // Handle bridge.host routing requests

    hr.Map( configuration.Server.DashboardDomain, dashboard.Resources{}.Routes() )   // Handle dashboard.host routing requests

    hr.Map( "*", baseRouter() ) // Handle everything else

    fmt.Println( "Host Router Ready\n" )

    router.Mount( "/", hr )

    fmt.Println( "Domain: " + configuration.Server.Domain + "\tPort: " + HTTPSPortString + "\n" )

    http.ListenAndServeTLS( ":" + HTTPSPortString, "./config/ssl/cert.pem", "./config/ssl/key.pem", router )  // Start the Server

}


func bridgeRouter() chi.Router {
    router := chi.NewRouter()

    router.Get( "/", renderBridge )

    return router
}

func renderBridge(w http.ResponseWriter, r *http.Request) {
    w.Write( []byte( "Bridge!\n" ) )
}

func dashboardRouter() chi.Router {
    router := chi.NewRouter()

    router.Get( "/", renderDashboard )

    return router
}

func renderDashboard(w http.ResponseWriter, r *http.Request) {
    w.Write( []byte( "Dashboard!\n" ) )
}

/* **************************************************************************
** Function:
** URI:
** Description:
** Client:
** *************************************************************************/
func buildWildRedirectURI( subdomain string, URI string ) string {
    return "/" + strings.Replace( URI, "/" + subdomain + "/", "", 1 )
}

/* **************************************************************************
** Function:
** URI:
** Description:
** Client:
** *************************************************************************/
func baseRouter() chi.Router {
    router := chi.NewRouter()

    // TODO: Append parameters

    // Redirect requests to host/access to access.host
    router.Get( "/access", func ( w http.ResponseWriter, r *http.Request ) {
        http.Redirect( w, r, "https://" + configuration.Server.AccessDomain, http.StatusMovedPermanently )
    })

    // Redirect requests to host/access/* to access.host/*
    router.Get( "/access/*", func ( w http.ResponseWriter, r *http.Request ) {
        redirectURI := buildWildRedirectURI( "access", r.RequestURI )
        fmt.Println( redirectURI )
        http.Redirect( w, r, "https://" + configuration.Server.AccessDomain + redirectURI, http.StatusMovedPermanently )
    })

    // Redirect requests to host/bridge to bridge.host
    router.Get( "/bridge", func ( w http.ResponseWriter, r *http.Request ) {
        http.Redirect( w, r, "https://" + configuration.Server.BridgeDomain, http.StatusMovedPermanently )
    })

    // Redirect requests to host/bridge/* to bridge.host/*
    router.Get( "/bridge/*", func ( w http.ResponseWriter, r *http.Request ) {
        redirectURI := buildWildRedirectURI( "bridge", r.RequestURI )
        http.Redirect( w, r, "https://" + configuration.Server.BridgeDomain + redirectURI, http.StatusMovedPermanently )
    })

    // Redirect requests to host/dashboard to dashboard.host
    router.Get( "/dashboard", func ( w http.ResponseWriter, r *http.Request ) {
        http.Redirect( w, r, "https://" + configuration.Server.DashboardDomain, http.StatusMovedPermanently )
    })

    // Redirect requests to host/dashboard/* to dashboard.host/*
    router.Get( "/dashboard/*", func ( w http.ResponseWriter, r *http.Request ) {
        redirectURI := buildWildRedirectURI( "dashboard", r.RequestURI )
        http.Redirect( w, r, "https://" + configuration.Server.DashboardDomain + redirectURI, http.StatusMovedPermanently )
    })

    router.Get( "/about", base.RenderAbout )

    router.Get( "/contact", base.RenderContact )

    router.Get( "/policy", base.RenderPolicy )

    router.Get( "/", renderIndex )

    cssDir := filepath.Join( workingDir, "/web/public/css" )    // Build the path to the CSS files
    scriptsDir := filepath.Join( workingDir, "/web/public/js" ) // Build the path to the JS files
    assetsDir := filepath.Join( workingDir, "/web/public/assets" )  // Build the path to the asset files

    // Setup the static file serving
    AddFileServer( router, "/css/", http.Dir( cssDir ) )
    AddFileServer( router, "/js/", http.Dir( scriptsDir ) )
    AddFileServer( router, "/assets/", http.Dir( assetsDir ) )

    return router
}

/* **************************************************************************
** Function: AddFileServer
** URI: /<path>/*
** Description: Adds the files from the supplied path to be served staticly
** *************************************************************************/
func AddFileServer(router chi.Router, path string, root http.FileSystem) {

    fs := http.StripPrefix( path, http.FileServer( root ) )

    trailingSlash := len( path ) - 1
    router.Get( path[ :trailingSlash ], http.RedirectHandler( path, 301 ).ServeHTTP )   // Redirect requests to /<path> to /<path>/

    // Serve files at <path>
    router.Get( path + "*", http.HandlerFunc( func(w http.ResponseWriter, r *http.Request) {
        fs.ServeHTTP( w, r )
    } ) )

}

/* **************************************************************************
** Function: renderIndex
** URI: /
** Description: Renders and servers the index page
** *************************************************************************/
func renderIndex(w http.ResponseWriter, r *http.Request) {

    qlDataPath := filepath.Join( workingDir, "/web/data/quick-links.json" )    // Build the path to quick links data
    tmDataPath := filepath.Join( workingDir, "/web/data/team.json" ) // Build the path to team member data

    indexCTX := make(map[string]interface{})

    indexCTX["title"] = "Brypt"
    indexCTX["quickLink"], _ = os.Open( qlDataPath )
    indexCTX["teamMember"], _ = os.Open( tmDataPath )

    page := handlebars.CompilePage( "index", indexCTX )

    w.Header().Set( "Content-Type", "text/html" )
    w.Write( []byte( page ) )

}
