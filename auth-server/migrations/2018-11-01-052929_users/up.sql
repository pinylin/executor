-- Your SQL goes here
CREATE TABLE users (
  openid VARCHAR(100) NOT NULL PRIMARY KEY,
  password VARCHAR(64) NOT NULL, --bcrypt hash
  created_at TIMESTAMP NOT NULL
);
