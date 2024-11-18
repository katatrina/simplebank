package api

import (
	"fmt"
	
	"github.com/gin-gonic/gin"
	db "github.com/katatrina/simplebank/db/sqlc"
	"github.com/katatrina/simplebank/mail"
	"github.com/katatrina/simplebank/token"
	"github.com/katatrina/simplebank/util"
	"github.com/katatrina/simplebank/worker"
)

const (
	EnvironmentProduction = "production"
	EnvironmentDevelop    = "development"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	router          *gin.Engine
	store           db.Store
	tokenMaker      token.Maker
	config          util.Config
	taskDistributor worker.TaskDistributor
	mailer          mail.EmailSender
}

func NewServer(store db.Store, config util.Config, taskDistributor worker.TaskDistributor, mailer mail.EmailSender) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSecretKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker <= %w", err)
	}
	
	server := &Server{
		store:           store,
		tokenMaker:      tokenMaker,
		config:          config,
		taskDistributor: taskDistributor,
		mailer:          mailer,
	}
	
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	switch server.config.Environment {
	case EnvironmentDevelop:
		gin.ForceConsoleColor()
	case EnvironmentProduction:
		gin.SetMode(gin.ReleaseMode)
	}
	
	router := gin.Default()
	v1 := router.Group("/v1")
	
	v1.GET("/health", server.healthCheck)
	
	v1.POST("/users", server.createUser)
	v1.PATCH("/users", authMiddleware(server.tokenMaker), server.updateUser)
	v1.POST("/users/login", server.loginUser)
	v1.GET("/users/verify_email", server.verifyUserEmail)
	
	v1.POST("/tokens/renew_access", server.renewAccessToken)
	
	authRoutes := v1.Group("/", authMiddleware(server.tokenMaker))
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)
	
	authRoutes.POST("/transfers", server.createTransfer)
	
	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"message": err.Error()}
}
