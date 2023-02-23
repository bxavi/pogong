CREATE TABLE account (
  id bigserial PRIMARY KEY,
  email varchar UNIQUE NOT NULL,
  password varchar NOT NULL,
  created_at timestamptz NOT NULL DEFAULT (now())
);

