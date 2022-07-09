package db

import (
	"bank-mvp/util"
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	var config util.Config
	config, err = util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Cannot read config:", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalln("Cannot connect to DB:", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}