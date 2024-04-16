package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"wsclean/db"
	"wsclean/handler"
	wsPkg "wsclean/websocket"
)

func main() {

	conn, err := db.PGConnect()
	if err != nil {
		panic(err)
	}

	// test db connection
	err = conn.Ping()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("DB Successfully connected!")
	}

	onlineHub := wsPkg.NewOnlineHub()
	go onlineHub.Run()

	// HTTP server
	e := echo.New()

	// Serve home.html
	e.Static("/", "home.html")

	e.POST("/users/:userid/follows", handler.FollowsStore)

	e.GET("/ws/:userid", func(c echo.Context) error {
		err := handler.WsHandler(c)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Something went wrong")
		}
		return nil
	})

	e.GET("/ws/online/:userid", func(c echo.Context) error {
		err := handler.WsNotificationHandler(c, onlineHub)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Something went wrong")
		}
		return nil
	})

	e.Start(":8081")

	return
}
