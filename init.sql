CREATE TABLE IF NOT EXISTS users (
	id BIGSERIAL PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	name TEXT NOT NULL);
	
CREATE TABLE IF NOT EXISTS segments( 
	name TEXT PRIMARY KEY NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	description TEXT
);
	
CREATE TABLE IF NOT EXISTS users_in_segments(
	user_id BIGSERIAL NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	segment_name TEXT NOT NULL REFERENCES segments(name) ON DELETE CASCADE,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	expire_at TIMESTAMP,
	PRIMARY KEY (user_id, segment_name)
);
	
CREATE TABLE IF NOT EXISTS users_in_segments_history(
	user_id BIGSERIAL NOT NULL,
	segment_name TEXT NOT NULL,
	expire_at TIMESTAMP,
	action_type TEXT NOT NULL,
	action_date TIMESTAMP NOT NULL
);

CREATE OR REPLACE FUNCTION users_in_segments_insert() 
RETURNS TRIGGER 
AS 
$$
BEGIN 
	INSERT INTO users_in_segments_history(user_id, segment_name, expire_at, action_type, action_date)
	VALUES (NEW.user_id, NEW.segment_name, NEW.expire_at, 'inserted', now());
	
RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER users_in_segments_after_insert 
	AFTER INSERT ON users_in_segments
	FOR EACH ROW
	EXECUTE PROCEDURE users_in_segments_insert();

CREATE OR REPLACE FUNCTION users_in_segments_delete() 
RETURNS TRIGGER 
AS 
$$
BEGIN 
	INSERT INTO users_in_segments_history(user_id, segment_name, expire_at, action_type, action_date)
	VALUES (OLD.user_id, OLD.segment_name, OLD.expire_at, 'deleted', now());
	
RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER users_in_segments_after_delete 
	AFTER DELETE ON users_in_segments
	FOR EACH ROW
	EXECUTE PROCEDURE users_in_segments_delete();