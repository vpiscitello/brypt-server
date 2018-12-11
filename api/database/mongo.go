package database

import(
    config "brypt-server/config"

    "github.com/mongodb/mongo-go-driver/mongo"
)

var configuration = config.Configuration{}

type key string

func Setup() {
    configuration = config.getConfigFile()


}
