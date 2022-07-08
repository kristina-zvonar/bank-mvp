package main

import (
	"bank-mvp/api"
	db "bank-mvp/db/sqlc"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var config map[string]string = map[string]string{
	"dbDriver": "postgres",
	"dbSource": "postgresql://root:root@localhost:5432/bank_mvp?sslmode=disable",
	"serverAddress": "0.0.0.0:8080",
}

func main() {
	conn, err := sql.Open(config["dbDriver"], config["dbSource"])
	if err != nil {
		log.Fatalln("Cannot connect to DB:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(config["serverAddress"])
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}