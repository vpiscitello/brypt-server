package config

import (
    "os"
    "path"
    "strings"
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
    err := gonfig.GetConf( getConfigFile(), &configuration )
    if err != nil {
        os.Exit( 500 )
    }
}

func GetConfig() Configuration {
    return configuration
}

/* **************************************************************************
** Function:
** URI:
** Description:h
** *************************************************************************/

func getConfigFile() string {
	env := os.Getenv( "ENV" )

	if len( env ) == 0 {
		env = "development"
	}

	filename := []string{ "config.", strings.TrimSpace(env), ".json" }
    filePath := path.Join( "/app/config/", strings.Join( filename, "") )

	return filePath
}
