CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ----------------------------

CREATE TYPE store_status AS ENUM ('open', 'close');

CREATE TABLE IF NOT EXISTS stores(
   id serial PRIMARY KEY,
   email varchar(512) NOT NULL,
   password varchar(15) NOT NULL,
   name varchar(64) NOT NULL,
   description text,
   created_at timestamp NOT NULL DEFAULT clock_timestamp(),
   status store_status NOT NULL,
   session_id uuid NOT NULL DEFAULT uuid_generate_v4()
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