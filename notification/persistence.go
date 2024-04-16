package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"wsclean/db"
)

func (ntf *Notification) StoreNotification(c context.Context) error {
	dbc, err := db.PGConnect()
	if err != nil {
		return err
	}

	attrs, err := json.Marshal(ntf.Attributes)
	if err != nil {
		fmt.Println("store_notification", err.Error)
		return err
	}

	ntf.ID = uuid.New()

	if err != nil {
		return fmt.Errorf("failed to generate uuuid: %+v", err)

	}

	err = dbc.GetContext(c, ntf, insertNotification, ntf.ID, ntf.UserID, ntf.Description, ntf.ActionUserID, ntf.Status, ntf.Type, ntf.SourceId, attrs)
	if err != nil {
		return fmt.Errorf("could not create notificaion: %+v", err)
	}
	return nil
}

func GetNotificationByUserId(ctx context.Context, userId string) ([]Notification, error) {
	dbc, err := db.PGConnect()
	if err != nil {
		return nil, err
	}

	args := map[string]interface{}{
		"user_id": userId,
	}
	query := fmt.Sprintf(getNotificationByUserId)
	ntfs := []Notification{}
	q, err := dbc.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, err
	}
	err = q.SelectContext(ctx, &ntfs, args)
	if err != nil {
		return nil, err
	}
	return ntfs, nil
}

/*
func UpdateNotification(c context.Context, ntf *Notification, up map[string]interface{}, id string) error {
	return common.UpdateRow(c, up, ntf, updateNotification, "id", id)
}
*/

// Delete ...
func Delete(ctx context.Context, uid string) (int, error) {
	dbc, err := db.PGConnect()
	if err != nil {
		return 0, err
	}
	// defer dbc.close()
	args := map[string]interface{}{
		"user_id": uid,
	}
	res, err := dbc.NamedExecContext(ctx, deleteNotification, args)
	if err != nil {
		return 0, err
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(ra), nil
}

const insertNotification = `
INSERT INTO notification
	(id, user_id, description, action_user_id, status, type, source_id, attributes)
VALUES
	($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *
`

const getNotificationByUserId = `
SELECT id, action_user_id, _created, description, _modified, source_id, status, type, user_id, attributes, count(*) OVER() AS full_count  FROM notification
WHERE user_id = :user_id AND _deleted = FALSE
ORDER BY
		_created DESC
`

const updateNotification = `
UPDATE notification
 SET %s
 WHERE
		%s = $%d
 AND
  _deleted = FALSE
RETURNING *
`

// deleteNotification ...
const deleteNotification = `
UPDATE 
	notification
SET
	_deleted = TRUE
WHERE
	user_id = :user_id
AND
	_deleted = FALSE
`
