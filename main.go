package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tijanadmi/simplebank/api"
	db "github.com/tijanadmi/simplebank/db/sqlc"
	"github.com/tijanadmi/simplebank/util"
)


func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
		//Err(err).Msg("cannot load config")
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
	server, err:=api.NewServer(config,store)
	if err != nil{
		log.Fatal("cannot create server:", err)
	}

	err=server.Start(config.HTTPServerAddress)
	if err != nil{
		log.Fatal("cannot start server:", err)
	}
}