package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var (
	Conn    *sql.DB
	cfgFile string
)

func init() {
	// connect to config toml
	viper.SetConfigFile("toml")
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("./configs")
		viper.SetConfigName("config")
	}
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("Not Load Viper Config")
	}
	// log.Println("Using Config File: ", viper.ConfigFileUsed())
}

func Connect() *sql.DB {
	host := viper.GetString("configDB.host")
	port := viper.GetString("configDB.port")
	user := viper.GetString("configDB.user")
	password := viper.GetString("configDB.password")
	dbname := viper.GetString("configDB.dbname")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	result, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		Conn = nil
		log.Println("Error Connection : ", err)
		return nil
	}
	Conn = result
	return result
}

func InsertChat(name string, chat string) error {
	createdAt := time.Now().Format("2006-01-02 15:04:05")
	sql := fmt.Sprintf("INSERT INTO chat (name_people,chat, created_at) VALUES ('%s', '%s', '%s'); ",
		name, chat, createdAt)
	_, errs := Conn.Query(sql)
	if errs != nil {
		return errs
	}
	return nil
}
