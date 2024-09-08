package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	driverName     = "postgres"
	dataSourceName = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var (
	testQueries  *Queries
	testConnPool *sql.DB
)

func TestMain(m *testing.M) {
	var err error
	
	testConnPool, err = sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer testConnPool.Close()

	testQueries = New(testConnPool)

	os.Exit(m.Run())
}
