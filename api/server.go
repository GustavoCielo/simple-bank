package api

import (
	"fmt"

	db "github.com/GustavoCielo/simple-bank/db/sqlc"
	"github.com/GustavoCielo/simple-bank/token"
	"github.com/GustavoCielo/simple-bank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server servers HTTP requests for our banking service.
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	// tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey) also works perfectly
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	// Binds the validation package to Gin server engine.
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	// add routes to router.
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	AuthRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	AuthRoutes.POST("/accounts", server.createAccount)
	AuthRoutes.GET("/accounts/:id", server.getAccount)
	AuthRoutes.GET("/accounts", server.listAccount)

	AuthRoutes.POST("/transfers", server.createTransfer)

	// No Auth
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
