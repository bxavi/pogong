CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "email" varchar NOT NULL,
  "password" varchar NOT NULL,
  "created_at" timestamptz DEFAULT 'now()'
);

INSERT INTO accounts (email, password) VALUES ('admin@email.com', 'password');