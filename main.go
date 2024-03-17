package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// album slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func main() {
	// Initialize a Gin router.
	router := gin.Default()

	// NOTE: Disallow trusted of all proxies.
	// $ ifconfig | grep "inet "
	// See https://github.com/gin-gonic/gin/issues/2809
	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"127.0.0.1"})

	// NOTE: this is handler function to associate HTTP method and path.

	// We pass the name of the `getAlbums` function, not the result (i.e. `getAlbums()`).
	router.GET("/albums", getAlbums)
	// The colon preceding an item in the path signifies that the item is a path parameter.
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	// Attach the router to an http.Server and start the server.
	router.Run("localhost:8080")
}

// NOTE: gin.Context carries request details, validates, serializes JSON, etc.

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	// Context.IndentedJSON serializes the struct into JSON and adds it to response.
	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum album

	// Call BindJSON to bind the received JSON to newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)

	// Add a 201 status code to the response, along with the JSON.
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter setn by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	// Context.Param retrieves the id path parameter from the URL.
	id := c.Param("id")

	// Loop over the list of albums.
	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
