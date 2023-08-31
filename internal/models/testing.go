package models

import "database/sql"

func NewTestUser() *User {
	return &User{
		ID:   1,
		Name: "testName",
	}
}

func NewTestSegment() *Segment {
	return &Segment{
		Name: "TestSegment",
		Description: sql.NullString{
			String: "Test Description",
			Valid:  true,
		},
	}
}
