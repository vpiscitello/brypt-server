package main

// Register Routes
// Run Server

import (
   "fmt"
   "strings"
   "strconv"
   "net/http"
   "net/url"

   config "brypt-server/config"

   "github.com/go-chi/chi"
   "github.com/go-chi/chi/middleware"

   "github.com/go-chi/hostrouter"
)

var configuration = config.Configuration{}

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

func main()  {
    config.Setup()  // Setup the Server Configuration
    configuration = config.GetConfig()  // Get the Configuration Settings

    HTTPPortString := strconv.Itoa( configuration.Server.HTTPPort )
    HTTPSPortString := strconv.Itoa( configuration.Server.HTTPSPort )

    go http.ListenAndServe( ":" + HTTPPortString, http.HandlerFunc( redirectToHTTPS ) )  // Start the Server

    router := chi.NewRouter()

    router.Use( middleware.RequestID )
    router.Use( middleware.RealIP )
    router.Use( middleware.Logger )
    router.Use( middleware.Recoverer )

    hr := hostrouter.New()

    hr.Map( configuration.Server.AccessDomain, accessRouter() ) // Handle access.host routing requests

    hr.Map( configuration.Server.BridgeDomain, bridgeRouter() ) // Handle bridge.host routing requests

    hr.Map( configuration.Server.DashboardDomain, dashboardRouter() )   // Handle dashboard.host routing requests

    hr.Map( "*", baseRouter() ) // Handle everything else

    router.Mount("/", hr)

    fmt.Println( "Domain: " + configuration.Server.Domain + "\tPort: " + HTTPPortString + "\n" )

    http.ListenAndServeTLS( ":" + HTTPSPortString, "/config/ssl/cert.pem", "/config/ssl/key.pem", router )  // Start the Server

}

func accessRouter() chi.Router {
    router := chi.NewRouter()

    router.Get( "/", renderAccess )

    return router
}

func renderAccess(w http.ResponseWriter, r *http.Request) {
    w.Write( []byte( "Access!\n" ) )
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

func baseRouter() chi.Router {
    router := chi.NewRouter()


    // Redirect requests to host/access to access.host
    router.Get( "/access", func ( w http.ResponseWriter, r *http.Request ) {
        http.Redirect( w, r, "https://" + configuration.Server.AccessDomain, http.StatusMovedPermanently )
    })

    // Redirect requests to host/access/* to access.host/*
    router.Get( "/access/*", func ( w http.ResponseWriter, r *http.Request ) {
        http.Redirect( w, r, "https://" + configuration.Server.AccessDomain + r.RequestURI, http.StatusMovedPermanently )
    })

    // Redirect requests to host/bridge to bridge.host
    router.Get( "/bridge", func ( w http.ResponseWriter, r *http.Request ) {
        http.Redirect( w, r, "https://" + configuration.Server.BridgeDomain, http.StatusMovedPermanently )
    })

    // Redirect requests to host/dashboard to dashboard.host
    router.Get( "/dashboard", func ( w http.ResponseWriter, r *http.Request ) {
        http.Redirect( w, r, "https://" + configuration.Server.DashboardDomain, http.StatusMovedPermanently )
    })

    //router.Get( "/", renderIndex )
    router.Mount( "/", http.FileServer( http.Dir( "../../docs/sketchpad/UI Elements" ) ) )


    return router
}

func renderIndex(w http.ResponseWriter, r *http.Request) {
    w.Write( []byte( "Hello World!\n" ) )
}
