package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/katatrina/simplebank/util"
	_ "github.com/lib/pq"
)

var (
	testQueries  *Queries
	testConnPool *sql.DB
)

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../app.env")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	testConnPool, err = sql.Open(config.DriverName, config.DataSourceName)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer testConnPool.Close()

	testQueries = New(testConnPool)

	os.Exit(m.Run())
}
