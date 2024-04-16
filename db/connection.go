package db

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"time"

	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// global postgresql connection
var conn *sqlx.DB

func PGConnect() (*sqlx.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "postgres", "admin", "postgres")

	dbc, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		log.Fatalln(err)
	}

	conn = dbc

	var count int
	err = conn.Get(&count, "SELECT COUNT(*) FROM notification")
	if err != nil {
		log.Println("Failed to execute query: ", err.Error())
	}

	err = dbc.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("DB Successfully connected!")
	// Use dbc for further database operations
	return conn, nil
}

func GetDBContext(c echo.Context, timeout ...int) (context.Context, context.CancelFunc) {
	r := c.Request()
	tm := 15
	if len(timeout) > 0 {
		tm = timeout[0]
	}
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(tm)*time.Second)
	return ctx, cancel
}
