package notification

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"time"
)

var NotificationTypes = map[string]string{
	"NotificationTypeBucketListTagged":    "You have been tagged in bucket list ",
	"NotificationTypeFollowed":            " started following you.",
	"NotificationTypeCommented":           " have commented on your post.",
	"NotificationTypeTripCopied":          " have copied your trip.",
	"NotificationTypeAddedToTrip":         "You have been added to trip ",
	"NotificationTypeArriveFollowRequest": "You have a new follow request from ",
	"NotificationTypeSuggestion":          " have suggested ",
}

type Notification struct {
	ID           uuid.UUID      `json:"id" db:"id"`
	Created      *time.Time     `json:"created,omitempty" db:"_created"`
	Modified     *time.Time     `json:"modified,omitempty" db:"_modified"`
	Deleted      bool           `json:"deleted,omitempty" db:"_deleted"`
	UserID       string         `json:"user_id" db:"user_id" validate:"required"`
	Description  string         `json:"description" db:"description" validated:"required"`
	ActionUserID string         `json:"action_user_id" db:"action_user_id" validate:"required"`
	Status       bool           `json:"status" db:"status" validate:"requried"`
	Type         string         `json:"type" db:"type" validate:"required"`
	SourceId     string         `json:"source_id" db:"source_id" validate:"required"`
	Attributes   AttributesType `json:"attributes" db:"attributes"`
}

type NotificaitionReponse struct {
	ID           uuid.UUID      `json:"id" db:"id"`
	Created      *time.Time     `json:"created,omitempty"`
	Modified     *time.Time     `json:"modified,omitempty"`
	Deleted      bool           `json:"deleted,omitempty"`
	UserID       string         `json:"user_id"`
	Description  string         `json:"description"`
	ActionUserID string         `json:"action_user_id"`
	Status       bool           `json:"status"`
	Type         string         `json:"type"`
	SourceId     string         `json:"source_id"`
	Attributes   AttributesType `json:"attributes"`
}

type AttributesType map[string]interface{}

func (a AttributesType) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *AttributesType) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

type NotificationParams struct {
	UserID       string                 `json:"user_id" db:"user_id" validate:"required"`
	Description  string                 `json:"description" db:"description" validated:"required"`
	ActionUserID string                 `json:"action_user_id" db:"action_user_id" validate:"required"`
	Status       bool                   `json:"status" db:"status" validate:"requried"`
	Type         string                 `json:"type" db:"type" validate:"required"`
	SourceId     string                 `json:"source_id" db:"source_id" validate:"required"`
	Attributes   map[string]interface{} `json:"attributes" db:"attributes"`
}

type NotificationUpdateRequestPayload struct {
	Status bool `json:"status"`
}

func NewNotification(options ...Option) (*Notification, error) {
	ntf := &Notification{}
	for _, o := range options {
		if err := o(ntf); err != nil {
			return nil, err
		}
	}
	return ntf, nil
}
