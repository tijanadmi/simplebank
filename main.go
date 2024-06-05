package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	_ "github.com/tijanadmi/simplebank/doc/statik"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tijanadmi/simplebank/api"
	db "github.com/tijanadmi/simplebank/db/sqlc"
	"github.com/tijanadmi/simplebank/gapi"
	"github.com/tijanadmi/simplebank/pb"
	"github.com/tijanadmi/simplebank/util"
)


func main() {
	

	config, err := util.LoadConfig(".")
	if err != nil{
		log.Fatal().Err(err).Msg("cannot load config")
	}

	
	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	//conn, err := pgx.Connect(context.Background(), DBSource)
	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	//defer conn.Close(context.Background())

	//testQueries = New(conn)
	store := db.NewStore(conn)
	// server, err:=api.NewServer(config,store)
	// if err != nil{
	// 	log.Fatal().Err(err).Msg("cannot create server:")
	// }

	// err=server.Start(config.HTTPServerAddress)
	// if err != nil{
	// 	log.Fatal().Err(err).Msg("cannot create server:")
	// }
	go runGatewayServer(config,store)
	runGrpcServer(config, store)
}

func runGrpcServer(
	/*ctx context.Context,
	waitGroup *errgroup.Group,*/
	config util.Config,
	store db.Store,
	//taskDistributor worker.TaskDistributor,
) {
	
	server, err := gapi.NewServer(config, store/*, taskDistributor*/)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	gprcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(gprcLogger)
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	log.Info().Msgf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err!=nil{
		log.Fatal().Err(err).Msg("cannot start gRPC server")
	}
		// if err != nil {
		// 	if errors.Is(err, grpc.ErrServerStopped) {
		// 		return nil
		// 	}
		// 	log.Error().Err(err).Msg("gRPC server failed to serve")
		// 	return err
		// }

	/*waitGroup.Go(func() error {
		log.Info().Msgf("start gRPC server at %s", listener.Addr().String())

		err = grpcServer.Serve(listener)
		if err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				return nil
			}
			log.Error().Err(err).Msg("gRPC server failed to serve")
			return err
		}

		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown gRPC server")

		grpcServer.GracefulStop()
		log.Fatal("gRPC server is stopped")

		return nil
	})*/
}

func runGatewayServer(
	/*ctx context.Context,
	waitGroup *errgroup.Group,*/
	config util.Config,
	store db.Store,
	/*taskDistributor worker.TaskDistributor,*/
) {
	server, err := gapi.NewServer(config, store/*, taskDistributor*/)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create statik fs")
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	log.Info().Msgf("start HTTP gateway server at %s", listener.Addr().String())
	handler:=gapi.HttpLogger(mux)
	err = http.Serve(listener,handler)
	if err!=nil{
		log.Fatal().Err(err).Msg("cannot start HTTP gateway server")
	}

	// statikFS, err := fs.New()
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("cannot create statik fs")
	// }

	// swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	// mux.Handle("/swagger/", swaggerHandler)

	// httpServer := &http.Server{
	// 	Handler: gapi.HttpLogger(mux),
	// 	Addr:    config.HTTPServerAddress,
	// }

	// waitGroup.Go(func() error {
	// 	log.Info().Msgf("start HTTP gateway server at %s", httpServer.Addr)
	// 	err = httpServer.ListenAndServe()
	// 	if err != nil {
	// 		if errors.Is(err, http.ErrServerClosed) {
	// 			return nil
	// 		}
	// 		log.Error().Err(err).Msg("HTTP gateway server failed to serve")
	// 		return err
	// 	}
	// 	return nil
	// })

	// waitGroup.Go(func() error {
	// 	<-ctx.Done()
	// 	log.Info().Msg("graceful shutdown HTTP gateway server")

	// 	err := httpServer.Shutdown(context.Background())
	// 	if err != nil {
	// 		log.Error().Err(err).Msg("failed to shutdown HTTP gateway server")
	// 		return err
	// 	}

	// 	log.Info().Msg("HTTP gateway server is stopped")
	// 	return nil
	// })
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start server")
	}
}