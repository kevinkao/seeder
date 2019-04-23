package main

import (
	"os"
	"github.com/spf13/viper"
	"database/sql"
	"github.com/spf13/cobra"
	// "github.com/manifoldco/promptui"
	"encoding/json"
	"io/ioutil"
	"fmt"
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

	var cmdRun =&cobra.Command{
		Use: "run",
		Short: "Run all the seed file in db.json",
		Run: func (cmd *cobra.Command, args []string) {
			seedFolder := databaseFolder + "/seeds"
			seedFile := seedFolder + "/db.json"
			if !FileExists(seedFile) {
				panic("No such db.json in seeds folder")
			}

			file, _ := ioutil.ReadFile(seedFile)
			var data map[string][]string
			err := json.Unmarshal(file, &data)
			handlerError(err)

			err = WithTransaction(db, func(tx *sql.Tx) (err error) {
				for i := 0; i < len(data["run"]); i++ {
					sql := data["run"][i]
					path := seedFolder + fmt.Sprintf("/%s", sql)
					fmt.Println(path)

					if !FileExists(path) {
						panic(fmt.Sprintf("No such sql file: %s", path))
					}

					content, err := ioutil.ReadFile(path)
					if err != nil {
						panic(err)
					}

					result, err := tx.Exec(string(content))
					if err != nil {
						panic(err)
					}

					if _, err := result.RowsAffected(); err != nil {
						panic(err)
					}
				}
				return
			})
			handlerError(err)
			
			return
		},
	}

	var cmdSql =&cobra.Command{
		Use: "sql [SQL_NAME]",
		Short: "Specify the sql file in seed folder",
		Run: func (cmd *cobra.Command, args []string) {
			err = WithTransaction(db, func(tx *sql.Tx) (err error) {
				// 
				return
			})
			handlerError(err)
		},
	}

	var rootCmd = &cobra.Command{Use: "seeder"}
	rootCmd.AddCommand(cmdRun, cmdSql)
	rootCmd.Execute()
}

func handlerError (err error) {
	if err != nil {
		panic(err)
	}
}