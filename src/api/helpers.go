package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func newPostgres(username, password, host, port, dbName string) (*pgx.Conn, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, dbName)
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func processor(c *gin.Context, incomingObj any) bool {
	if err := c.ShouldBind(incomingObj); err != nil {
		c.String(http.StatusBadRequest, "bad request: %v", err)
		return false
	}
	return true
}
func scanRow(row pgx.Row) (Book, error) {
	var book Book
	err := row.Scan(&book.ISBN, &book.Title, &book.Author, &book.Date, &book.Edition, &book.Genre, &book.Description)
	return book, err
}

func errorHandler(err error, c *gin.Context, message string) bool {
	if err != nil {
		log.Print(message)
		c.String(http.StatusBadRequest, "bad request: %v", err)
		return true
	}
	return false
}
