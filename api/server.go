package api

import (
	"fmt"
	
	"github.com/gin-gonic/gin"
	db "github.com/katatrina/simplebank/db/sqlc"
	"github.com/katatrina/simplebank/token"
	"github.com/katatrina/simplebank/util"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	router     *gin.Engine
	store      db.Store
	tokenMaker token.Maker
	config     util.Config
}

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
	
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/tokens/renew_access", server.renewAccessToken)
	
	authRoutes := router.Group("/", authMiddleware(server.tokenMaker))
	authRoutes.POST("accounts", server.createAccount)
	authRoutes.GET("accounts/:id", server.getAccount)
	authRoutes.GET("accounts", server.listAccounts)
	
	authRoutes.POST("transfers", server.createTransfer)
	
	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
