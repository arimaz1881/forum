package adapters

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"forum/internal/domain"
)

type NotificationsRepositorySqlite3 struct {
	db *sql.DB
}

var _ domain.NotificationsRepository = (*NotificationsRepositorySqlite3)(nil)

func NewNotificationsRepositorySqlite3(db *sql.DB) *NotificationsRepositorySqlite3 {
	return &NotificationsRepositorySqlite3{
		db: db,
	}
}

const createNotificationFromReaction = `
INSERT INTO
    notifications (
        post_id,
		author_id,
		action_id,
		comment_id,
        action_type,
        seen
    )
VALUES
    (?, ?, ?, ?, ?, ?)
ON CONFLICT (action_id)
DO UPDATE SET
    action_type = EXCLUDED.action_type;
`

func (q *NotificationsRepositorySqlite3) Create(ctx context.Context, input domain.CreateNotificationInput) error {
	if input.CommentID == 0 {
		input.CommentID = time.Now().UnixNano()
	}

	if input.ActionID == 0 {
		input.ActionID = time.Now().UnixNano()
	}
	_, err := q.db.ExecContext(
		ctx,
		createNotificationFromReaction,
		input.PostID,
		input.AuthorID,
		input.ActionID,
		input.CommentID,
		input.Action,
		input.Seen,
	)

	if err != nil {
		return err
	}

	return nil
}

const deleteNotification = `
DELETE FROM
  notifications
WHERE
  action_id =?;
`

func (q *NotificationsRepositorySqlite3) Delete(ctx context.Context, reactionID int64) error {
	_, err := q.db.ExecContext(ctx, deleteNotification, reactionID)
	return err
}

const lookNotification = `
UPDATE
  notifications
SET
  seen = true
WHERE 
  id = ?;
`

func (q *NotificationsRepositorySqlite3) Look(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, lookNotification, id)
	return err
}

const getNotificationsList = `
SELECT
    n.id,
    n.created_at,
    n.post_id,
    n.author_id,
    n.action_type,
	n.seen,
	p.title
FROM
	notifications n
JOIN 
	posts p ON p.id = n.post_id
WHERE 
	n.author_id = ? AND n.seen = false 
ORDER BY
	n.created_at DESC;
`

func (q *NotificationsRepositorySqlite3) GetList(ctx context.Context, userID int64) ([]domain.Notification, error) {
	rows, err := q.db.QueryContext(ctx, getNotificationsList, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var notifications []domain.Notification

	for rows.Next() {
		var notification domain.Notification
		err := rows.Scan(
			&notification.ID,
			&notification.CreatedAt,
			&notification.PostID,
			&notification.AuthorID,
			&notification.Action,
			&notification.Seen,
			&notification.PostTitle,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

const getNotification = `
SELECT
    id,
    created_at,
    post_id,
    author_id,
    action_type
FROM
    notifications
WHERE id = ?;
`

func (q *NotificationsRepositorySqlite3) GetOne(ctx context.Context, notificationID string) (*domain.Notification, error) {
	row := q.db.QueryRowContext(ctx, getNotification, notificationID)
	var notification domain.Notification

	err := row.Scan(
		&notification.ID,
		&notification.CreatedAt,
		&notification.PostID,
		&notification.AuthorID,
		&notification.Action,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotificationNotFound
		}
		return nil, err
	}

	return &notification, nil
}
