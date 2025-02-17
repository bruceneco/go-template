
-- +migrate Up
CREATE TABLE users (
  id UUID primary key,
  name VARCHAR NOT NULL,
  email VARCHAR NOT NULL UNIQUE,
  password CHAR(64) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
  deleted_at TIMESTAMP WITH TIME ZONE
);
-- +migrate Down
DROP TABLE users;