CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ----------------------------

CREATE TABLE IF NOT EXISTS stores(
   id serial PRIMARY KEY,
   email varchar(512) NOT NULL UNIQUE,
   password varchar(64) NOT NULL,
   name varchar(64) NOT NULL,
   description text DEFAULT '',
   created_at timestamp NOT NULL DEFAULT clock_timestamp(),
   session_id uuid NOT NULL DEFAULT uuid_generate_v4()
);

-- ----------------------------

CREATE TYPE sign_key_types AS ENUM ('signin', 'email');

CREATE TABLE IF NOT EXISTS sign_keys(
   id serial PRIMARY KEY,
   store_id integer REFERENCES stores ON DELETE CASCADE,
   sign_key varchar(64) NOT NULL,
   type sign_key_types NOT NULL,
   created_at timestamp NOT NULL DEFAULT clock_timestamp()
);

-- ----------------------------

CREATE TABLE IF NOT EXISTS queues(
   id serial PRIMARY KEY,
   name varchar(64) NOT NULL,
   store_id integer REFERENCES stores ON DELETE CASCADE
);

-- ----------------------------

CREATE TYPE customer_status AS ENUM ('processing', 'done', 'delete');

CREATE TABLE IF NOT EXISTS customers(
   id serial PRIMARY KEY,
   name varchar(64) NOT NULL,
   phone varchar(30) NOT NULL,
   queue_id integer REFERENCES queues ON DELETE CASCADE,
   status customer_status NOT NULL,
   created_at timestamp NOT NULL DEFAULT clock_timestamp()
);

-- ----------------------------