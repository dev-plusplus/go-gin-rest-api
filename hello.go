package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// CRUD: Create, Read, Update and Delete

type Task struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CompletedAt string `json:"completedAt"`
}

var tasks []Task = []Task{
	{Name: "Task Name", Description: "Task Description"},
}

func MySQLConnect() *sql.DB {
	db, err := sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db
}

func CreateTask(c *gin.Context) {
	db := MySQLConnect()
	defer db.Close()

	stmtIns, err := db.Prepare("INSERT INTO Task VALUES( ?, ?, ?, ? )") // ? = placeholder
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Malformed JSON when creating a Task",
		})
		return
	}
	defer stmtIns.Close()

	var newTask Task
	var error = c.BindJSON(&newTask)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Malformed JSON when creating a Task",
		})
		return
	}

	result, err := stmtIns.Exec(nil, newTask.Name, newTask.Description, nil) // Insert tuples (i, i^2)
	fmt.Println(result)
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"data":    newTask,
	})
}

func ReadTasks(c *gin.Context) {
	db := MySQLConnect()
	defer db.Close()

	rows, err := db.Query("SELECT id, name, description, completedAt FROM Task ORDER BY id DESC")

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error executing the query",
		})
		return
	}

	var tasks []Task = []Task{}
	var id, name, description, completedAt []byte
	// Fetch rows
	for rows.Next() {
		var newTask Task
		// get RawBytes from data
		err = rows.Scan(&id, &name, &description, &completedAt)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		fmt.Println("ID:", string(id))
		fmt.Println("NAME:", string(name))
		fmt.Println("DESCRIPTION:", string(description))
		fmt.Println("completedAt:", string(completedAt))

		newTask.ID = string(id)
		newTask.CompletedAt = string(completedAt)
		newTask.Name = string(name)
		newTask.Description = string(description)

		tasks = append(tasks, newTask)
	}
	if err = rows.Err(); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error executing the query",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": tasks,
	})
}

func ReadTask(c *gin.Context) {
	var taskId, error = strconv.Atoi(c.Param("id"))
	fmt.Println("taskId:", taskId)
	db := MySQLConnect()
	defer db.Close()

	query, err := db.Prepare("SELECT id, name, description, completedAt FROM Task WHERE id = ? ORDER BY id DESC")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error preparing the query",
		})
		return
	}

	var id, name, description, completedAt []byte
	row := query.QueryRow(taskId)
	err = row.Scan(&id, &name, &description, &completedAt)
	if err != nil {
		fmt.Println("ERROR:", error)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Not a valid ID",
		})
		return
	}

	var newTask Task
	newTask.ID = string(id)
	newTask.CompletedAt = string(completedAt)
	newTask.Name = string(name)
	newTask.Description = string(description)

	c.JSON(http.StatusNotFound, gin.H{
		"data": newTask,
	})
}

func UpdateTask(c *gin.Context) {
	var taskId, error = strconv.Atoi(c.Param("id"))
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Missing task ID",
		})
		return
	}

	var newTask Task
	var err = c.BindJSON(&newTask)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Malformed JSON when creating a Task",
		})
		return
	}

	db := MySQLConnect()
	defer db.Close()

	query, err := db.Prepare("UPDATE Task SET name = ?, description = ?, completedAt = ? WHERE id = ? ")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error preparing the query",
		})
		return
	}

	_, err = query.Exec(newTask.Name, newTask.Description, newTask.CompletedAt, taskId)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error running the query",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task Updated",
		"data":    newTask,
	})
}

func DeleteTask(c *gin.Context) {
	var taskId, error = strconv.Atoi(c.Param("id"))
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Missing task ID",
		})
		return
	}

	db := MySQLConnect()
	defer db.Close()

	query, err := db.Prepare("DELETE FROM Task WHERE id = ? ")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error preparing the query",
		})
		return
	}

	_, err = query.Exec(taskId)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error running the query",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task Deleted",
	})
}

func main() {
	// Loading environment variables
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Error loading .env file")
	}

	application := gin.Default()
	application.POST("/tasks", CreateTask)
	application.GET("/tasks", ReadTasks)
	application.GET("/tasks/:id", ReadTask)
	application.PUT("/tasks/:id", UpdateTask)
	application.DELETE("/tasks/:id", DeleteTask)
	application.Run()
}
