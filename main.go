package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	supa "github.com/supabase-community/supabase-go"
)

var supabase *supa.Client

type Todo struct {
    ID          int    `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
}



func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	log.SetReportCaller(true)
	// gin.SetMode(gin.ReleaseMode)
	r := gin.New()


	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage)
	}))

	r.Use(gin.Recovery())

	supabaseUrl := os.Getenv("SUPABASE_URL")
	supabaseAnonKey := os.Getenv("SUPABASE_KEY")
	
	supabaseClient, err := supa.NewClient(supabaseUrl, supabaseAnonKey, nil)
	if err != nil {
		log.Fatal("Error initializing Supabase client:", err)
		panic(err)
	}
	supabase = supabaseClient

	r.GET("/todos", getTodos)
	r.POST("/todos", createTodo)
	r.GET("/todos/:id", getTodo)
	r.PUT("/todos/:id", updateTodo)
	r.DELETE("/todos/:id", deleteTodo)

	log.Info("Starting server on :8080")
	log.WithFields(logrus.Fields{
		"port": 8080,
	}).Info("Listening for requests")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}

}


func getTodos(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("size", "10"))


	if page < 1{
		page = 1
	}
	if perPage < 1{
		perPage = 10
	}

	start := (page - 1) * perPage
    end := start + perPage

	var todos []Todo
	_, err := supabase.From("todos").Select("*", "exact", false).Range(start, end, "").ExecuteTo(&todos)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}


	c.JSON(http.StatusOK, gin.H{"todos": todos})
}

func createTodo(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "Create a new todo"})
}

func getTodo(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Get todo " + id})
}

func updateTodo(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Update todo " + id})
}

func deleteTodo(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Delete todo " + id})
}