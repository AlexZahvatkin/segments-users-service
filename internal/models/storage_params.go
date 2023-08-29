package models

import (
	"database/sql"
	"time"
)

type GetSegmentsHistoryByUserIdParams struct {
	UserID   int64
	FromDate time.Time
	ToDate   time.Time
}

type RemoveUserFromSegmentParams struct {
	UserID      int64
	SegmentName string
}

type AddUserIntoSegmentParams struct {
	UserID      int64
	SegmentName string
}

type AddUserIntoSegmentWithExpireDatetimeParams struct {
	UserID      int64
	SegmentName string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ExpireAt    sql.NullTime
}

type AddSegmentParams struct {
	Name        string         `json:"name" validate:"required,min=4,max=255"`
	Description sql.NullString `json:"description" validate:"max=65535"`
}