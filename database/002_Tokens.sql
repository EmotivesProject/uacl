\c uacl_db;

CREATE TABLE tokens (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,

	CONSTRAINT fk_username
	FOREIGN KEY(username)
	REFERENCES users(username)
);