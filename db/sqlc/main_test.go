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
	testQueries *Queries
)

func TestMain(m *testing.M) {
	connPool, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer connPool.Close()

	testQueries = New(connPool)

	os.Exit(m.Run())
}
