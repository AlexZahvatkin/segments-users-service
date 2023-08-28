DROP TRIGGER IF EXISTS users_in_segments_after_delete ON users_in_segments;
DROP TRIGGER IF EXISTS users_in_segments_after_insert ON users_in_segments;
DROP TABLE IF EXISTS users_in_segments;
DROP TABLE IF EXISTS segments;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS users_in_segments_history;