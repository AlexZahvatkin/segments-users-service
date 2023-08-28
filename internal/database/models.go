// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1

package database

import (
	"database/sql"
	"time"
)

type Segment struct {
	Name        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Description sql.NullString
}

type User struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
}

type UsersInSegment struct {
	UserID      int64
	SegmentName string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ExpireAt    sql.NullTime
}

type UsersInSegmentsHistory struct {
	UserID      int64
	SegmentName string
	ExpireAt    sql.NullTime
	ActionType  string
	ActionDate  time.Time
}
