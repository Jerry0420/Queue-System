CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE store_status AS ENUM ('normal', 'close');

CREATE TABLE IF NOT EXISTS store(
   id serial PRIMARY KEY,
   email varchar(512) NOT NULL,
   password varchar(15) NOT NULL,
   name varchar(64) NOT NULL,
   description text,
   created_at timestamp NOT NULL DEFAULT clock_timestamp(),
   status store_status NOT NULL,
   updated_at timestamp NOT NULL DEFAULT clock_timestamp(),
);