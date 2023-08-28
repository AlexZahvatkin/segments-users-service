// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: user_in_segments.sql

package database

import (
	"context"
)

const addUserIntoSegment = `-- name: AddUserIntoSegment :one
INSERT INTO users_in_segments (user_id, segment_name, created_at, updated_at, expire_at) 
VALUES ($1, $2, now(), now(), null)
ON CONFLICT (user_id, segment_name) DO UPDATE
	SET updated_at = now(), expire_at = null
RETURNING user_id, segment_name, created_at, updated_at, expire_at
`

type AddUserIntoSegmentParams struct {
	UserID      int64
	SegmentName string
}

func (q *Queries) AddUserIntoSegment(ctx context.Context, arg AddUserIntoSegmentParams) (UsersInSegment, error) {
	row := q.db.QueryRowContext(ctx, addUserIntoSegment, arg.UserID, arg.SegmentName)
	var i UsersInSegment
	err := row.Scan(
		&i.UserID,
		&i.SegmentName,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ExpireAt,
	)
	return i, err
}

const addUserIntoSegmentWithTTLInHours = `-- name: AddUserIntoSegmentWithTTLInHours :one
INSERT INTO users_in_segments (user_id, segment_name, created_at, updated_at, expire_at) 
VALUES ($1, $2, now(), now(), now() + make_interval(hours => $3))
ON CONFLICT (user_id, segment_name) DO UPDATE
	SET updated_at = now(), expire_at = now() + make_interval(hours => $3)
RETURNING user_id, segment_name, created_at, updated_at, expire_at
`

type AddUserIntoSegmentWithTTLInHoursParams struct {
	UserID        int64
	SegmentName   string
	NumberOfHours int32
}

func (q *Queries) AddUserIntoSegmentWithTTLInHours(ctx context.Context, arg AddUserIntoSegmentWithTTLInHoursParams) (UsersInSegment, error) {
	row := q.db.QueryRowContext(ctx, addUserIntoSegmentWithTTLInHours, arg.UserID, arg.SegmentName, arg.NumberOfHours)
	var i UsersInSegment
	err := row.Scan(
		&i.UserID,
		&i.SegmentName,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ExpireAt,
	)
	return i, err
}

const getSegmentsByUserId = `-- name: GetSegmentsByUserId :many
SELECT segment_name 
FROM users_in_segments
WHERE user_id = $1 AND 
CASE WHEN expire_at IS NOT NULL
THEN expire_at > now()
ELSE TRUE
END
`

func (q *Queries) GetSegmentsByUserId(ctx context.Context, userID int64) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getSegmentsByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var segment_name string
		if err := rows.Scan(&segment_name); err != nil {
			return nil, err
		}
		items = append(items, segment_name)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const removeUserFromSegment = `-- name: RemoveUserFromSegment :exec
DELETE 
FROM users_in_segments
WHERE user_id = $1 AND 
segment_name = $2
`

type RemoveUserFromSegmentParams struct {
	UserID      int64
	SegmentName string
}

func (q *Queries) RemoveUserFromSegment(ctx context.Context, arg RemoveUserFromSegmentParams) error {
	_, err := q.db.ExecContext(ctx, removeUserFromSegment, arg.UserID, arg.SegmentName)
	return err
}
