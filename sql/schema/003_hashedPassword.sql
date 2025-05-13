-- +goose Up
ALTER TABLE users ADD hashed_password varchar(255) DEFAULT 'unset' NOT NULL;