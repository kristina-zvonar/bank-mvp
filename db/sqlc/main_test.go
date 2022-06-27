package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var config map[string]string = map[string]string{
	"dbDriver": "postgres",
	"dbSource": "postgresql://root:root@localhost:5432/bank_mvp?sslmode=disable",
}

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open(config["dbDriver"], config["dbSource"])
	if err != nil {
		log.Fatalln("Cannot connect to DB:", err)
	}

	testQueries = New(conn)
	os.Exit(m.Run())
}