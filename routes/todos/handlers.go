package todos

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Capitalize the first letter of each function name to make them exported
func GetTodos(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get all todos"})
}

func CreateTodo(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "Create a new todo"})
}

func GetTodo(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Get todo " + id})
}

func UpdateTodo(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Update todo " + id})
}

func DeleteTodo(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Delete todo " + id})
}