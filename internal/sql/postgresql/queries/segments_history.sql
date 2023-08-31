-- name: GetSegmentsHistoryByUserId :many
SELECT *
FROM users_in_segments_history
WHERE user_id = $1
    AND action_date > @from_date
    AND action_date < @to_date;