package config

// Setup Configuration
// Setup MySql Connection
// Setup MongoDB Connection

import (
    // "fmt"
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

func Setup() {
    err := gonfig.GetConf( getConfigFile(), &configuration )
    if err != nil {
        os.Exit( 500 )
    }
}

func GetConfig() Configuration {
    return configuration
}

func getConfigFile() string {
	env := os.Getenv( "ENV" )

	if len( env ) == 0 {
		env = "development"
	}

	filename := []string{ "config.", env, ".json" }
	_, dirname, _, _ := runtime.Caller( 0 )
	filePath := path.Join( filepath.Dir( dirname ), strings.Join( filename, "") )

	return filePath
}
