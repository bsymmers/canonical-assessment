package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Book struct {
	ISBN        string  `json:"isbn"`
	Title       string  `json:"title"`
	Author      string  `json:"author"`
	Date        string  `json:"published_date"`
	Edition     float64 `json:"edition"`
	Genre       string  `json:"genre"`
	Description string  `json:"description"`
}

type Collection struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// TODO: add test data

func NewRouter() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.POST("/books", postBooks)
	router.DELETE("/books", deleteBooks)
	router.GET("/collections", getCollections)
	router.POST("/collections", postCollections)
	router.DELETE("/collections", deleteCollections)
	router.PUT("/collections", putCollections)
	router.Run("localhost:8080")
}

func processor(c *gin.Context, incomingObj any) bool {
	if err := c.ShouldBind(&incomingObj); err != nil {
		c.String(http.StatusBadRequest, "bad request: %v", err)
		return false
	}
	return true
}

// Book Endpoints
func getBooks(c *gin.Context) {
	//TODO: check if isbn is in db

}

func postBooks(c *gin.Context) {
	var newBook Book

	isValid := processor(c, newBook)
	if !isValid {
		return
	}

}

func deleteBooks(c *gin.Context) {
	// TODO: check if isbn is in db

}

// Collection Endpoints
func getCollections(c *gin.Context) {
	// TODO: check if name is in db

}

func postCollections(c *gin.Context) {
	var newCollection Collection

	isValid := processor(c, newCollection)
	if !isValid {
		return
	}

}

func deleteCollections(c *gin.Context) {
	// TODO: check if name is in db

}

func putCollections(c *gin.Context) {
	// TODO: check if name is in db and if book isbn is in db

}
