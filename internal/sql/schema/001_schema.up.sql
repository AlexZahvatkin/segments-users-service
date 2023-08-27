CREATE TABLE users (
	id BIGSERIAL PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	name TEXT);
	
CREATE TABLE segments( 
	id BIGSERIAL PRIMARY KEY,
	name TEXT UNIQUE NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL
);
	
CREATE TABLE users_in_segments(
	user_id BIGSERIAL NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	segment_id BIGSERIAL NOT NULL REFERENCES segments(id) ON DELETE CASCADE,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	will_expire_at TIMESTAMP,
	PRIMARY KEY (user_id, segment_id)
);
	
CREATE TABLE users_in_segments_history(
	user_id BIGSERIAL NOT NULL,
	segment_id BIGSERIAL NOT NULL,
	segment_name TEXT NOT NULL,
	action_type TEXT NOT NULL,
	action_date TIMESTAMP NOT NULL
);

CREATE OR REPLACE FUNCTION users_in_segments_insert() 
RETURNS TRIGGER 
AS 
$$
BEGIN 
	INSERT INTO users_in_segments_history(user_id, segment_id, segment_name, action_type, action_date)
	VALUES (NEW.user_id, NEW.segment_id, (SELECT name FROM segments WHERE id = NEW.segment_id), 'inserted', now());
	
RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_in_segments_after_insert 
	AFTER INSERT ON users_in_segments
	FOR EACH ROW
	EXECUTE PROCEDURE users_in_segments_insert();

CREATE OR REPLACE FUNCTION users_in_segments_delete() 
RETURNS TRIGGER 
AS 
$$
BEGIN 
	INSERT INTO users_in_segments_history(user_id, segment_id, segment_name, action_type, action_date)
	VALUES (OLD.user_id, OLD.segment_id, (SELECT name FROM segments WHERE id = OLD.segment_id), 'deleted', now());
	
RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_in_segments_after_delete 
	AFTER DELETE ON users_in_segments
	FOR EACH ROW
	EXECUTE PROCEDURE users_in_segments_delete();
