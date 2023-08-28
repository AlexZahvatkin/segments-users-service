-- name: GetSegmentsByUserId :many
SELECT segment_id, name 
FROM users_in_segments
INNER JOIN segments ON segment_id = segments.id
WHERE user_id = @user_id AND 
CASE WHEN expire_at IS NOT NULL
THEN expire_at > now()
ELSE TRUE
END;

-- name: AddUserIntoSegment :one
INSERT INTO users_in_segments (user_id, segment_id, created_at, updated_at, expire_at) 
VALUES (@user_id, (SELECT id FROM segments WHERE @segment_name = name), now(), now(), null)
ON CONFLICT (user_id, segment_id) DO UPDATE
	SET updated_at = now(), expire_at = null
RETURNING *;

-- name: AddUserIntoSegmentWithTTLInHours :one
INSERT INTO users_in_segments (user_id, segment_id, created_at, updated_at, expire_at) 
VALUES (@user_id, (SELECT id FROM segments WHERE @segment_name = name), now(), now(), now() + make_interval(hours => @number_of_hours))
ON CONFLICT (user_id, segment_id) DO UPDATE
	SET updated_at = now(), expire_at = now() + make_interval(hours => @number_of_hours)
RETURNING *;

-- name: RemoveUserFromSegment :exec 
DELETE 
FROM users_in_segments
WHERE user_id = @user_id AND 
segment_id = (SELECT id FROM segments WHERE name = @segment_name);