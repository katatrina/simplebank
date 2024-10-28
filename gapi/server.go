package gapi

import (
	"fmt"
	
	db "github.com/katatrina/simplebank/db/sqlc"
	"github.com/katatrina/simplebank/pb"
	"github.com/katatrina/simplebank/token"
	"github.com/katatrina/simplebank/util"
)

// Server serves gRPC requests for our banking service.
type Server struct {
	pb.UnimplementedSimpleBankServer
	store      db.Store
	tokenMaker token.Maker
	config     util.Config
}

// NewServer creates a new gRPC server.
func NewServer(store db.Store, config util.Config) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSecretKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker <= %w", err)
	}
	
	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}
	
	return server, nil
}
