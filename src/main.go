package main

import (
	"canonical-REST/api"
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func NewPostgres(username, password, host, port, dbName string) (*pgx.Conn, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, dbName)
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func main() {

	conn, err := NewPostgres("postgres", os.Getenv("DATABASE_PASSWORD"), "localhost", "5432", "Book Management System")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	// var greeting string
	// err = conn.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	// 	os.Exit(1)
	// }

	// fmt.Println(greeting)
	api.NewRouter()
	defer conn.Close(context.Background())
}
