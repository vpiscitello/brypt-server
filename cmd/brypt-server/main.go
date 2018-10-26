package main

// Register Routes
// Run Server

import (
   "fmt"
   "strconv"
   "net/http"

   config "brypt-server/config"

   "github.com/go-chi/chi"
   "github.com/go-chi/chi/middleware"

   "github.com/go-chi/hostrouter"
)

var configuration = config.Configuration{}

func main()  {
    config.Setup()  // Setup the Server Configuration
    configuration = config.GetConfig()  // Get the Configuration Settings

    portString := strconv.Itoa( configuration.Server.Port )

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

    fmt.Println( "Domain: " + configuration.Server.Domain + "\tPort: " + portString + "\n" )

    http.ListenAndServe( ":" + strconv.Itoa( configuration.Server.Port ), router )  // Start the Server

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
        http.Redirect( w, r, "http://" + configuration.Server.AccessDomain, http.StatusMovedPermanently )
    })

    // Redirect requests to host/bridge to bridge.host
    router.Get( "/bridge", func ( w http.ResponseWriter, r *http.Request ) {
        http.Redirect( w, r, "http://" + configuration.Server.BridgeDomain, http.StatusMovedPermanently )
    })

    // Redirect requests to host/dashboard to dashboard.host
    router.Get( "/dashboard", func ( w http.ResponseWriter, r *http.Request ) {
        http.Redirect( w, r, "http://" + configuration.Server.DashboardDomain, http.StatusMovedPermanently )
    })

    router.Get( "/", renderIndex )


    return router
}

func renderIndex(w http.ResponseWriter, r *http.Request) {
    w.Write( []byte( "Hello World!\n" ) )
}
