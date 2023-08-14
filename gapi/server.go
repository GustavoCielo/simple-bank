package gapi

import (
	"fmt"

	db "github.com/GustavoCielo/simple-bank/db/sqlc"
	"github.com/GustavoCielo/simple-bank/pb"
	"github.com/GustavoCielo/simple-bank/token"
	"github.com/GustavoCielo/simple-bank/util"
	"github.com/GustavoCielo/simple-bank/worker"
)

// Server servers gRPC requests for our banking service.
type Server struct {
	pb.UnimplementedSimpleBankServer
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

// NewServer creates a new gRPC server.
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	// tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey) also works perfectly
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
