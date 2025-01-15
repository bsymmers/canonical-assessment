package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Book struct {
	ISBN        string      `json:"isbn"`
	Title       string      `json:"title"`
	Author      string      `json:"author"`
	Date        pgtype.Date `json:"published_date"`
	Edition     float64     `json:"edition"`
	Genre       string      `json:"genre"`
	Description string      `json:"description"`
}

type Collection struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// TODO: add test data
var DB_CONN *pgx.Conn

func newPostgres(username, password, host, port, dbName string) (*pgx.Conn, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, dbName)
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

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

func processor(c *gin.Context, incomingObj any) bool {
	if err := c.ShouldBind(incomingObj); err != nil {
		c.String(http.StatusBadRequest, "bad request: %v", err)
		return false
	}
	// fmt.Println(incomingObj)
	return true
}

// Book Endpoints
// TODO: add filtering
func getBooks(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	if len(queryParams) == 0 {
		query := `
		SELECT * FROM books
	`
		rows, err := DB_CONN.Query(context.Background(), query)
		if err != nil {
			log.Printf("Error Fetching Book Details")
			c.String(http.StatusBadRequest, "bad request: %v", err)
			return
		}
		defer rows.Close()

		var books []Book
		// Iterate over the retrieved rows and scan each row into a Book struct.
		for rows.Next() {
			var book Book
			err := rows.Scan(&book.ISBN, &book.Title, &book.Author, &book.Date, &book.Edition, &book.Genre, &book.Description)
			if err != nil {
				log.Printf("Error Fetching Book Details")
				c.String(http.StatusBadRequest, "bad request: %v", err)
				return
			}
			books = append(books, book)
		}
		c.JSON(http.StatusAccepted, books)
	} else {
		query := `
		SELECT * FROM books WHERE isbn = @isbn
	`
		args := pgx.NamedArgs{
			"isbn": queryParams["isbn"][0],
		}
		row := DB_CONN.QueryRow(context.Background(), query, args)
		var book Book
		err := row.Scan(&book.ISBN, &book.Title, &book.Author, &book.Date, &book.Edition, &book.Genre, &book.Description)
		if err != nil {
			log.Printf("Error Fetching Book Details")
			c.String(http.StatusBadRequest, "bad request: %v", err)
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
	if err != nil {
		log.Println("Error Inserting Book Details")
		c.String(http.StatusBadRequest, "bad request: %v", err)
	}
	c.JSON(http.StatusAccepted, newBook)
}

func deleteBooks(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	query := `
		DELETE FROM books WHERE isbn = @isbn
	`
	args := pgx.NamedArgs{
		"isbn": queryParams["isbn"][0],
	}

	_, err := DB_CONN.Exec(context.Background(), query, args)
	if err != nil {
		log.Printf("Error Fetching Book Details")
		c.String(http.StatusBadRequest, "bad request: %v", err)
		return
	}
	c.Status(200)

}

// Collection Endpoints
func getCollections(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	if len(queryParams) == 0 {
		query := `
		SELECT * FROM collections
	`
		rows, err := DB_CONN.Query(context.Background(), query)
		if err != nil {
			log.Printf("Error Fetching Collection Details")
			c.String(http.StatusBadRequest, "bad request: %v", err)
			return
		}
		defer rows.Close()

		var collections []Collection
		// Iterate over the retrieved rows and scan each row into a Book struct.
		for rows.Next() {
			var collection Collection
			err := rows.Scan(&collection.Name, &collection.Description)
			if err != nil {
				log.Printf("Error Fetching Book Details")
				c.String(http.StatusBadRequest, "bad request: %v", err)
				return
			}
			collections = append(collections, collection)
		}
		c.JSON(http.StatusAccepted, collections)
	} else {
		query := `
		SELECT * FROM collections WHERE name = @name
	`
		args := pgx.NamedArgs{
			"name": queryParams["name"][0],
		}
		row := DB_CONN.QueryRow(context.Background(), query, args)
		var collection Collection
		err := row.Scan(&collection.Name, &collection.Description)
		if err != nil {
			log.Printf("Error Fetching Book Details")
			c.String(http.StatusBadRequest, "bad request: %v", err)
			return
		}
		c.JSON(http.StatusAccepted, collection)

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
	if err != nil {
		log.Println("Error Inserting Book into Collection")
		c.String(http.StatusBadRequest, "bad request: %v", err)
	}
	c.JSON(http.StatusAccepted, newCollection)
}

func deleteCollections(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	query := `
        DELETE FROM collections WHERE collection_name = @name
    `
	args := pgx.NamedArgs{
		"name": queryParams["name"][0],
	}
	_, err := DB_CONN.Exec(context.Background(), query, args)
	if err != nil {
		log.Printf("Error Fetching Book Details")
		c.String(http.StatusBadRequest, "bad request: %v", err)
		return
	}
	query = `
        DELETE FROM books_collections WHERE collection_name = @name
    `
	_, err = DB_CONN.Exec(context.Background(), query, args)
	if err != nil {
		log.Printf("Error Fetching Book Details")
		c.String(http.StatusBadRequest, "bad request: %v", err)
		return
	}
	c.Status(200)

}

func putCollections(c *gin.Context) {
	// TODO: call getbook and then depending on response add book to collections books table
	queryParams := c.Request.URL.Query()
	query := `
		SELECT * FROM books WHERE isbn = @isbn
	`
	args := pgx.NamedArgs{
		"isbn": queryParams["isbn"][0],
	}
	row := DB_CONN.QueryRow(context.Background(), query, args)
	var book Book
	err := row.Scan(&book.ISBN, &book.Title, &book.Author, &book.Date, &book.Edition, &book.Genre, &book.Description)
	if err != nil {
		log.Printf("Error Fetching Book Details")
		c.String(http.StatusBadRequest, "bad request: %v", err)
		return
	}
	toDelete, ok := queryParams["to_delete"]
	// If the key exists
	if ok && toDelete[0] == "false" {
		query = `
			DELETE FROM books_collections WHERE collection_name = @name AND isbn = @isbn
			`
	} else {
		query = `
			INSERT INTO books_collections (isbn, collection_name) VALUES (@isbn, @collection_name)
			`
	}
	args = pgx.NamedArgs{
		"isbn":            queryParams["isbn"][0],
		"collection_name": queryParams["collection_name"][0],
	}
	_, err = DB_CONN.Exec(context.Background(), query, args)
	if err != nil {
		log.Println("Error Inserting Book into Collection")
		c.String(http.StatusBadRequest, "bad request: %v", err)
	}
	c.Status(200)

}
