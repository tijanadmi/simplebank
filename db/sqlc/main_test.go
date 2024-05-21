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


func TestMain(m *testing.M) {
	

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