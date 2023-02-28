CREATE TABLE "account" (
  "id" bigserial PRIMARY KEY,
  "email" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "account" ("email");


INSERT INTO account (email, password) VALUES ('admin@email.com', '$2a$10$1F9WK1GD8WYUEGLpTk3eE.QmiKEPA/wCJPjkc7RO5Cqf7G4kyJgme');