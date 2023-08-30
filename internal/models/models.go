package models

import (
	"database/sql"
	"time"
)

type Segment struct {
	Name        string
	Description sql.NullString
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UsersInSegment struct {
	UserID      int64
	SegmentName string
	ExpireAt    sql.NullTime
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UsersInSegmentsHistory struct {
	UserID      int64
	SegmentName string
	ExpireAt    sql.NullTime
	ActionType  string
	ActionDate  time.Time
}

type AddUserIntoSegmentWithTTLInHoursParams struct {
	UserID        int64
	SegmentName   string
	NumberOfHours int32
}
