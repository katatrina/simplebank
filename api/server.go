package api

import (
	"fmt"
	
	"github.com/gofiber/fiber/v2"
	db "github.com/katatrina/simplebank/db/sqlc"
	"github.com/katatrina/simplebank/token"
	"github.com/katatrina/simplebank/util"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	router     *fiber.App
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
	}
	
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	app := fiber.New()
	
	// if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	// 	v.RegisterValidation("currency", validCurrency)
	// }
	
	app.Post("/users", server.createUser)
	app.Post("/users/login", server.loginUser)
	
	app.Post("/accounts", server.createAccount)
	app.Get("/accounts/:id", server.getAccount)
	app.Get("/accounts", server.listAccounts)
	
	app.Post("/transfers", server.createTransfer)
	
	server.router = app
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Listen(address)
}

func errorResponse(err error) fiber.Map {
	return fiber.Map{"error": err.Error()}
}
