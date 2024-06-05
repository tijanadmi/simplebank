package gapi

import (
	"fmt"

	db "github.com/tijanadmi/simplebank/db/sqlc"
	"github.com/tijanadmi/simplebank/pb"
	"github.com/tijanadmi/simplebank/token"
	"github.com/tijanadmi/simplebank/util"
)

// Server serves gRPC requests for our banking service.
type Server struct {
	pb.UnimplementedSimpleBankServer
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	//taskDistributor worker.TaskDistributor
}

// NewServer creates a new gRPC server.
func NewServer(config util.Config, store db.Store/*, taskDistributor worker.TaskDistributor*/) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		//taskDistributor: taskDistributor,
	}

	return server, nil
}