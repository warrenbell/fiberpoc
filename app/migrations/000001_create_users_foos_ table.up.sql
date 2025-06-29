CREATE OR REPLACE FUNCTION current_epoch_milliseconds()
RETURNS bigint AS $$
BEGIN
  RETURN EXTRACT(EPOCH FROM NOW()) * 1000;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS foos(
   id serial PRIMARY KEY,
   name VARCHAR (50) NOT NULL,
   created_at bigint DEFAULT current_epoch_milliseconds(),
   updated_at bigint DEFAULT 0,
   deleted_at bigint DEFAULT 0
);

CREATE TABLE IF NOT EXISTS users(
   id serial PRIMARY KEY,
   name VARCHAR (50) NOT NULL,
   email VARCHAR (300) NOT NULL,
   created_at bigint DEFAULT current_epoch_milliseconds(),
   updated_at bigint DEFAULT 0,
   deleted_at bigint DEFAULT 0,
   CHECK (email <> '')
);