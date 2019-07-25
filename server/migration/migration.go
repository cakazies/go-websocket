package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-websocket/server/models"
	"github.com/go-websocket/server/utils"
)

var (
	Conn *sql.DB
)

func main() {
	Conn = models.Connect()
	MigrationRooms()
}

func MigrationRooms() {
	tableName := "chat"
	queryCreate := fmt.Sprintf(`
					CREATE TABLE public.%s
					(
						id SERIAL NOT NULL,
						name_people character varying(200) COLLATE pg_catalog."default" NOT NULL,
						chat character varying(2000) COLLATE pg_catalog."default" NOT NULL,
						created_at timestamp without time zone NOT NULL,
						CONSTRAINT %s_pk PRIMARY KEY (id)
					);`, tableName, tableName)
	stmt, err := Conn.Prepare(queryCreate)
	utils.FailError(err, fmt.Sprintf("Error Create Table %s ", tableName))
	_, err = stmt.Exec()
	utils.FailError(err, fmt.Sprintf("Error Create Table %s ", tableName))
	log.Println(fmt.Sprintf("Import Table %s Succesfull", tableName))
}
