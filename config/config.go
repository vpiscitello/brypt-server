package config

// Setup Configuration
// Setup MySql Connection
// Setup MongoDB Connection

import (
    "fmt"
    "os"
    "path"
    "path/filepath"
    "runtime"
    "strings"

   // "encoding/json"

   "github.com/tkanos/gonfig"
)

type Configuration struct {
    Server struct {
        Domain string
        AccessDomain string
        BridgeDomain string
        DashboardDomain string
        HTTPPort int
        HTTPSPort int
    }
    Database struct {
        MongoURI string
    }
    Session struct {
        Options struct {
            Path string
            Domain string
            MaxAge int
            Secure bool
            UseHTTPS bool
            HTTPSOnly bool
        }
    }
    Template struct {
        Root string
        Children []string
    }
    View struct {
        BaseURI string
        Extension string
        Folder string
        Name string
        Caching bool
    }
}

var configuration = Configuration{}

/* **************************************************************************
** Function:
** URI:
** Description:
** *************************************************************************/

func Setup() {
    fmt.Println( "Configuration Setup\n" )
    err := gonfig.GetConf( getConfigFile(), &configuration )
    if err != nil {
        os.Exit( 500 )
    }
}

func GetConfig() Configuration {
    fmt.Println( "Get Configuration File\n" )
    return configuration
}

/* **************************************************************************
** Function:
** URI:
** Description:
** *************************************************************************/

func getConfigFile() string {
    fmt.Println( "Get Configuration FilePath\n" )
	env := os.Getenv( "ENV" )

	if len( env ) == 0 {
		env = "development"
	}

	filename := []string{ "config.", env, ".json" }
	_, dirname, _, _ := runtime.Caller( 0 )
	filePath := path.Join( filepath.Dir( dirname ), strings.Join( filename, "") )

	return filePath
}
