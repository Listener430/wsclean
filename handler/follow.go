package handler

import (
	"github.com/labstack/echo/v4"
	"wsclean/notification"
)

type RequestBody struct {
	UserID       string                 `json:"userID"`
	ActionUserID string                 `json:"actionUserID"`
	Description  string                 `json:"description"`
	Status       bool                   `json:"status"`
	Type         string                 `json:"type"`
	SourceId     string                 `json:"sourceId"`
	Attributes   map[string]interface{} `json:"attributes"`
}

// FollowsStore ...
func FollowsStore(c echo.Context) error {
	//ctx, cancel := db.GetDBContext(c)
	//defer cancel()

	// Create an instance of RequestBody
	var body RequestBody

	// Bind the request body to the body variable
	if err := c.Bind(&body); err != nil {
		return err
	}

	// Assign values from the request body to the notificationParam fields
	notificationParam := notification.NotificationParams{
		UserID:       body.UserID,
		ActionUserID: body.ActionUserID,
		Description:  body.Description,
		Status:       body.Status,
		Type:         body.Type,
		SourceId:     body.SourceId,
		Attributes:   body.Attributes,
	}

	go StoreNotification(c, notificationParam)

	return c.JSON(201, body)
}
