package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Book struct {
	ISBN        string      `json:"isbn" binding:"required"`
	Title       string      `json:"title" binding:"required"`
	Author      string      `json:"author" binding:"required"`
	Date        pgtype.Date `json:"published_date" binding:"required"`
	Edition     float64     `json:"edition"`
	Genre       string      `json:"genre"`
	Description string      `json:"description"`
}

type Collection struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type PutBody struct {
	Name     string `json:"collection_name" binding:"required"`
	ISBN     string `json:"isbn" binding:"required"`
	ToDelete bool   `json:"to_delete"`
}

var DB_CONN *pgx.Conn

func NewRouter() {
	var err error
	router := gin.Default()
	DB_CONN, err = newPostgres("postgres", os.Getenv("DATABASE_PASSWORD"), "localhost", "5432", "Book Management System")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	router.GET("/books", getBooks)
	router.POST("/books", postBooks)
	router.DELETE("/books", deleteBooks)
	router.GET("/collections", getCollections)
	router.POST("/collections", postCollections)
	router.DELETE("/collections", deleteCollections)
	router.PUT("/collections", putCollections)
	router.Run("localhost:8080")
	defer DB_CONN.Close(context.Background())
}

// Book Endpoints
func getBooks(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	if len(queryParams) == 0 {
		query := `
		SELECT * FROM books
	`
		rows, err := DB_CONN.Query(context.Background(), query)
		if errorHandler(err, c, "Error Fetching Book Details") {
			return
		}
		defer rows.Close()
		var books []Book
		for rows.Next() {
			var book Book
			err := rows.Scan(&book.ISBN, &book.Title, &book.Author, &book.Date, &book.Edition, &book.Genre, &book.Description)
			if errorHandler(err, c, "Error Fetching Book Details") {
				return
			}
			books = append(books, book)
		}
		c.JSON(http.StatusAccepted, books)
	} else {
		var query strings.Builder
		query.WriteString(fmt.Sprintf("SELECT * FROM books WHERE isbn = '%s'", queryParams["isbn"][0]))
		if len(c.Query("author")) > 0 {
			query.WriteString(fmt.Sprintf(" AND author = '%s'", c.Query("author")))
		}
		if len(c.Query("genre")) > 0 {
			query.WriteString(fmt.Sprintf(" AND genre = '%s'", c.Query("genre")))
		}
		if len(c.Query("published_date")) > 0 {
			query.WriteString(fmt.Sprintf(" AND date = '%s'", c.Query("date")))
		}
		s := query.String()
		row := DB_CONN.QueryRow(context.Background(), s)
		book, e := scanRow(row)
		if errorHandler(e, c, "Error Fetching Book Details") {
			return
		}
		c.JSON(http.StatusAccepted, book)

	}

}

func postBooks(c *gin.Context) {
	var newBook Book
	isValid := processor(c, &newBook)
	if !isValid {
		return
	}
	query := `
		INSERT INTO books (isbn, title, author, date, edition, genre, description) VALUES (@isbn, @title, @author, @date, @edition, @genre, @description)
	`
	args := pgx.NamedArgs{
		"isbn":        newBook.ISBN,
		"title":       newBook.Title,
		"author":      newBook.Author,
		"date":        newBook.Date,
		"edition":     newBook.Edition,
		"genre":       newBook.Genre,
		"description": newBook.Description,
	}
	_, err := DB_CONN.Exec(context.Background(), query, args)
	if errorHandler(err, c, "Error Inserting Book Details") {
		return
	}
	c.JSON(http.StatusCreated, newBook)
}

func deleteBooks(c *gin.Context) {
	query := `
        DELETE FROM books_collections WHERE isbn = @isbn
    `

	args := pgx.NamedArgs{
		"isbn": c.Query("isbn"),
	}

	_, err := DB_CONN.Exec(context.Background(), query, args)
	if errorHandler(err, c, "Error Fetching Book Details") {
		return
	}
	query = `
		DELETE FROM books WHERE isbn = @isbn
	`
	_, err = DB_CONN.Exec(context.Background(), query, args)
	if errorHandler(err, c, "Error Fetching Book Details") {
		return
	}
	c.Status(204)

}

// Collection Endpoints
func getCollections(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	if len(queryParams) == 0 {
		query := `
		SELECT * FROM collections
	`
		rows, err := DB_CONN.Query(context.Background(), query)
		if errorHandler(err, c, "Error Fetching Collection Details") {
			return
		}
		defer rows.Close()

		var collections []Collection
		for rows.Next() {
			var collection Collection
			err := rows.Scan(&collection.Name, &collection.Description)
			if errorHandler(err, c, "Error Fetching Collection Details") {
				return
			}
			collections = append(collections, collection)
		}
		c.JSON(http.StatusAccepted, collections)
	} else {
		query := `
		SELECT books.isbn, books.title, books.author, books.date, books.edition, books.genre, books.description
		 FROM books_collections
		 INNER JOIN books ON books_collections.isbn=books.isbn AND books_collections.collection_name = @name
	`
		args := pgx.NamedArgs{
			"name": c.Query("collection_name"),
		}
		row := DB_CONN.QueryRow(context.Background(), query, args)
		book, e := scanRow(row)
		if errorHandler(e, c, "Error Fetching Collection Details") {
			return
		}
		c.JSON(http.StatusAccepted, gin.H{c.Query("collection_name"): book})

	}

}

func postCollections(c *gin.Context) {
	var newCollection Collection
	isValid := processor(c, &newCollection)
	if !isValid {
		return
	}
	query := `
	INSERT INTO collections (collection_name, description) VALUES (@collection_name, @description)
	`
	args := pgx.NamedArgs{
		"collection_name": newCollection.Name,
		"description":     newCollection.Description,
	}
	_, err := DB_CONN.Exec(context.Background(), query, args)
	if errorHandler(err, c, "Error Inserting Book into Collection") {
		return
	}
	c.JSON(http.StatusCreated, newCollection)
}

func deleteCollections(c *gin.Context) {
	query := `
        DELETE FROM books_collections WHERE collection_name = @name
    `
	args := pgx.NamedArgs{
		"name": c.Query("collection_name"),
	}
	_, err := DB_CONN.Exec(context.Background(), query, args)
	if errorHandler(err, c, "Error Fetching Collection Details") {
		return
	}
	query = `
        DELETE FROM collections WHERE collection_name = @name
    `
	_, err = DB_CONN.Exec(context.Background(), query, args)
	if errorHandler(err, c, "Error Fetching Collection Details") {
		return
	}
	c.Status(204)

}

func putCollections(c *gin.Context) {
	var newPutBody PutBody
	isValid := processor(c, &newPutBody)
	if !isValid {
		return
	}
	query := `
		SELECT * FROM books WHERE isbn = @isbn
	`
	args := pgx.NamedArgs{
		"isbn": newPutBody.ISBN,
	}
	row := DB_CONN.QueryRow(context.Background(), query, args)
	_, e := scanRow(row)
	if errorHandler(e, c, "Error Fetching Book Details") {
		return
	}
	if newPutBody.ToDelete {
		query = `
			DELETE FROM books_collections WHERE collection_name = @name AND isbn = @isbn
			`
	} else {
		query = `
			INSERT INTO books_collections (isbn, collection_name) VALUES (@isbn, @collection_name)
			`
	}
	args = pgx.NamedArgs{
		"isbn":            newPutBody.ISBN,
		"collection_name": newPutBody.Name,
	}
	_, err := DB_CONN.Exec(context.Background(), query, args)
	if errorHandler(err, c, "Error Inserting Book into Collection") {
		return
	}
	c.Status(202)

}
