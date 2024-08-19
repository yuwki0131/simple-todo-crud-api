package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/BurntSushi/toml"
    "github.com/gin-gonic/gin"
    _ "github.com/go-sql-driver/mysql"
)

type Config struct {
    Database struct {
        User     string
        Password string
        Host     string
        Port     int
        Dbname   string
    }
}

type Todo struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    LimitedAt   time.Time `json:"limited_at"`
}

var db *sql.DB
var config Config

func initConfig() {
    if _, err := toml.DecodeFile("config.toml", &config); err != nil {
        log.Fatal("Could not load config file: ", err)
    }
}

func initDB() {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
        config.Database.User,
        config.Database.Password,
        config.Database.Host,
        config.Database.Port,
        config.Database.Dbname)

    var err error
    db, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal(err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }
}

func main() {
    initConfig()
    initDB()

    r := gin.Default()

    r.GET("/todos", getTodos)
    r.POST("/todos", createTodo)
    r.PUT("/todos/:id", updateTodo)
    r.DELETE("/todos/:id", deleteTodo)

    r.Run(":8080")
}

func getTodos(c *gin.Context) {
    rows, err := db.Query(
      "SELECT id, title, description, created_at, updated_at, limited_at FROM todos")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    var todos []Todo
    for rows.Next() {
        var todo Todo
        err := rows.Scan(&todo.ID,
                &todo.Title,
                &todo.Description,
                &todo.CreatedAt,
                &todo.UpdatedAt,
                &todo.LimitedAt)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        todos = append(todos, todo)
    }

    c.JSON(http.StatusOK, todos)
}

func createTodo(c *gin.Context) {
    var todo Todo
    if err := c.BindJSON(&todo); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    result, err := db.Exec(
      "INSERT INTO todos (title, description, limited_at) VALUES (?, ?, ?)",
      todo.Title,
      todo.Description,
      todo.LimitedAt)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    id, err := result.LastInsertId()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    todo.ID = int(id)
    todo.CreatedAt = time.Now()
    todo.UpdatedAt = time.Now()
    c.JSON(http.StatusOK, todo)
}

func updateTodo(c *gin.Context) {
    id := c.Param("id")
    var todo Todo
    if err := c.BindJSON(&todo); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    _, err := db.Exec(
      "UPDATE todos SET title = ?, description = ?, limited_at = ? WHERE id = ?",
      todo.Title,
      todo.Description,
      todo.LimitedAt,
      id)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, todo)
}

func deleteTodo(c *gin.Context) {
    id := c.Param("id")

    _, err := db.Exec("DELETE FROM todos WHERE id = ?", id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}
