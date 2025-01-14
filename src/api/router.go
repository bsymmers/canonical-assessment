package api

import "github.com/gin-gonic/gin"

type Book struct {
	ISBN        string  `json:"isbn"`
	Title       string  `json:"title"`
	Author      string  `json:"author"`
	Date        float64 `json:"published_date"`
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
	router.Run("localhost:8080")

	router.GET("/books", getBooks)
	router.POST("/books", postBooks)
	router.DELETE("/books", deleteBooks)
	router.GET("/collections", getCollections)
	router.POST("/collections", postCollections)
	router.DELETE("/collections", deleteCollections)
	router.PUT("/collections", putCollections)
}

// Book Endpoints
func getBooks(c *gin.Context) {

}

func postBooks(c *gin.Context) {

}

func deleteBooks(c *gin.Context) {

}

// Collection Endpoints
func getCollections(c *gin.Context) {

}

func postCollections(c *gin.Context) {

}

func deleteCollections(c *gin.Context) {

}

func putCollections(c *gin.Context) {

}
