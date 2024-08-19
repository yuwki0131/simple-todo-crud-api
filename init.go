package main

import (
    "database/sql"
    "fmt"
    "io/ioutil"
    "log"
    "time"

    "github.com/BurntSushi/toml"
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
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?tls=false",
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

    ddl, err := ioutil.ReadFile("db/schema.sql")
    if err != nil {
        log.Fatal("Could not read schema.sql: ", err)
    }

    _, err = db.Exec(string(ddl))
    if err != nil {
        log.Fatal(err)
    }
}

func main() {
    initConfig()
    initDB()
}

// その他のCRUD操作のハンドラは省略（前述の例と同様です）
