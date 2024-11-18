package db

import (
	"context"
	"os"
	"testing"
	
	"github.com/rs/zerolog/log"
	
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/katatrina/simplebank/util"
)

var testStore Store

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../app.env")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config file")
	}
	
	connPool, err := pgxpool.New(context.Background(), config.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create connection pool")
	}
	defer connPool.Close()
	
	testStore = NewStore(connPool)
	
	os.Exit(m.Run())
}
