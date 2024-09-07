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
	var newTodo Todo
	if err := c.BindJSON(&newTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	todo := map[string]interface{}{
        "title":       newTodo.Title,
        "description": newTodo.Description,
    }

	var inserted []Todo
	_, err := supabase.From("todos").Insert(todo, false, "", "", "exact").ExecuteTo(&inserted)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(inserted) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No data inserted"})
		return
	}


	c.JSON(http.StatusCreated, gin.H{"message": "Create a new todo", "data": inserted[0]})
}

func getTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var todo Todo
	_, err = supabase.From("todos").Select("*", "exact", false).Single().Eq("id", strconv.Itoa(id)).ExecuteTo(&todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}


	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Get todo %d", id), "data": todo})
}

func updateTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	fmt.Println("id", id)

	var newTodo Todo
	if err := c.BindJSON(&newTodo); err != nil {
		fmt.Println("err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	updateData := map[string]interface{}{
		"title":       newTodo.Title,
		"description": newTodo.Description,
		// Add other fields you want to update, but exclude 'id'
	}

	var updated []Todo
	_, err = supabase.From("todos").Update(updateData, "", "").Eq("id", strconv.Itoa(id)).ExecuteTo(&updated)
	if err != nil {
		fmt.Println("err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(updated) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Update todo %d", id), "data": updated[0]})
}

func deleteTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var deleted []Todo
	_, err = supabase.From("todos").Delete("","").Eq("id", strconv.Itoa(id)).ExecuteTo(&deleted)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(deleted) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Delete todo %d", id), "data": deleted[0]})
}