package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"wsclean/notification"
	wsPkg "wsclean/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WsNotificationHandler(c echo.Context, hub *wsPkg.OnlineHub) error {
	fmt.Printf("WsNoficationHandler stared\n")

	userId := c.Param("userid")
	if userId == "" {
		return echo.NewHTTPError(400, "User ID is required")
	}
	fmt.Println("userId: ", userId)

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v", err)
		return err
	}
	fmt.Printf("Websocket connection established\n")

	client := &wsPkg.OnlineClient{
		OnlineHub:           hub,
		WebSocketConnection: ws,
		Send:                make(chan wsPkg.OnlineEventStructure),
		UserId:              userId,
	}
	go client.ReadPump()
	go client.WritePump()
	client.OnlineHub.Register <- client

	return nil
}

func WsHandler(c echo.Context) error {
	fmt.Printf("WsHandler stared\n")

	userId := c.Param("userid")
	if userId == "" {
		return echo.NewHTTPError(400, "User ID is required")
	}
	fmt.Println("userId: ", userId)

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v", err)
		return err
	}
	fmt.Printf("Websocket connection established\n")

	cl := wsPkg.NotificationClient{Client: ws}

	wsPkg.NotificationPools[userId] = cl

	return nil
}

func StoreNotification(c echo.Context, ntc notification.NotificationParams) {
	ntf, err := notification.NewNotification()
	if err != nil {
		fmt.Println("\n store_notification - fetch_user_from_list", err)
		return
	}
	// Note: this is a stream of data
	// should not defer cancel() it.
	// ctx, _ := getDBContext(c)
	ctx := context.Background()

	attrs := ntc.Attributes
	ntf.UserID = ntc.UserID
	ntf.ActionUserID = ntc.ActionUserID
	ntf.Description = ntc.Description
	ntf.Status = ntc.Status
	ntf.Type = ntc.Type
	ntf.SourceId = ntc.SourceId
	ntf.Attributes = attrs

	err = ntf.StoreNotification(ctx)
	if err != nil {
		fmt.Println("store_notification", err.Error())
		return
	}

	pool := wsPkg.NotificationPools[ntc.UserID]

	fmt.Println("Trying to get the client from pool for user_id: ", ntc.UserID)

	if pool.Client != nil {

		// Reading from the db existing notifications for the current user

		fmt.Printf("/n store_notification - pool.Client != nil for user_id= %s", ntc.UserID)
		resp, err := notification.GetNotificationByUserId(ctx, ntc.UserID)
		if err != nil {
			fmt.Println("store_notification", err.Error())
			return
		}

		ntfs := make([]notification.NotificaitionReponse, 0)

		for _, ntf := range resp {
			newntf := notification.NotificaitionReponse{}
			if ntf.Type == "SUG" {
				//do nothing

			}

			bt, err := json.Marshal(ntf)
			if err != nil {
				fmt.Println("store_notification error ", err.Error())
				// continue
			}

			if err := json.Unmarshal(bt, &newntf); err != nil {
				fmt.Println("store_notification", err.Error())
				return
			}

			newntf.Attributes = ntf.Attributes
			ntfs = append(ntfs, newntf)
		}

		ntfsJson, err := json.Marshal(ntfs)
		if err != nil {
			fmt.Println("Error marshaling ntfs to JSON: ", err.Error())
			return
		}

		err = pool.BroadcastMessage(ntfsJson)
		if err != nil {
			fmt.Println("store_notification", err.Error)
		}
		fmt.Println("Broadsast message sent")
	}
}
