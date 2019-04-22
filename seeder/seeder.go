package main

import (
	"os"
	"github.com/spf13/viper"
	"database/sql"
)

type DBConnInfo struct {
	username, password, host, port, database, charset string
}

var configFolder string
var databaseFolder string
var dbEnvConfig string
var connInfo *DBConnInfo

func main () {
	databaseFolder = os.Getenv("GOPATH") + "/database"
	configFolder = os.Getenv("GOPATH") + "/config"
	dbEnvConfig = configFolder + "/database.env.json"

	viper.SetConfigType("json")
	viper.SetConfigName("database")
	viper.AddConfigPath(configFolder)
	viper.ReadInConfig()

	if FileExists(dbEnvConfig) {
		viper.SetConfigName("database.env")
		viper.MergeInConfig()
	}

	mk := viper.GetString("default")
	connInfo = &DBConnInfo {
		username: viper.GetString(mk+".username"),
		password: viper.GetString(mk+".password"),
		host: viper.GetString(mk+".host"),
		port: viper.GetString(mk+".port"),
		database: viper.GetString(mk+".database"),
		charset: viper.GetString(mk+".charset"),
	}

	db, err := DbConn(connInfo)
	if err != nil {
		panic(err)
	}

	err = WithTransaction(db, func(tx *sql.Tx) (err error) {
		// 
		return
	})
	handlerError(err)
}

func handlerError (err error) {
	if err != nil {
		panic(err)
	}
}