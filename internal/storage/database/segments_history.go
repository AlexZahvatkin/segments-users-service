// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: segments_history.sql

package database

import (
	"context"

	"github.com/AlexZahvatkin/segments-users-service/internal/models"
)

const getSegmentsHistoryByUserId = `-- name: GetSegmentsHistoryByUserId :many
SELECT user_id, segment_name, expire_at, action_type, action_date 
FROM users_in_segments_history
WHERE user_id = $1 AND action_date > $2 AND action_date < $3
`

func (q *Queries) GetSegmentsHistoryByUserId(ctx context.Context, arg models.GetSegmentsHistoryByUserIdParams) ([]models.UsersInSegmentsHistory, error) {
	rows, err := q.db.QueryContext(ctx, getSegmentsHistoryByUserId, arg.UserID, arg.FromDate, arg.ToDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.UsersInSegmentsHistory
	for rows.Next() {
		var i models.UsersInSegmentsHistory
		if err := rows.Scan(
			&i.UserID,
			&i.SegmentName,
			&i.ExpireAt,
			&i.ActionType,
			&i.ActionDate,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
