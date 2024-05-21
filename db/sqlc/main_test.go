package db

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testStore Store

const (
	dbDriver="postgres"
	DBSource="postgres://postgres:postgres@localhost:5432/simple_bank"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	/*config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testStore = NewStore(connPool)
	os.Exit(m.Run())*/

	/*conn, err :=sql.Open(dbDriver, DBSource)
	if err !=nil{
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())*/

	//conn, err := pgx.Connect(context.Background(), DBSource)
	conn, err := pgxpool.New(context.Background(), DBSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	//defer conn.Close(context.Background())

	//testQueries = New(conn)
	testStore = NewStore(conn)
	os.Exit(m.Run())
}