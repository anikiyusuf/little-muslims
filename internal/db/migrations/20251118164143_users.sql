-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TYPE gender_types AS ENUM ('male', 'female', 'other');

CREATE TABLE "users" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "email" VARCHAR(100) UNIQUE NOT NULL,
  "first_name" VARCHAR(100) NOT NULL,
  "last_name" VARCHAR(100) NOT NULL, 
  "password_hash" VARCHAR(255) NOT NULL,
  "gender" gender_types,
  "phone_number" VARCHAR(20) UNIQUE,
  "is_verified" BOOLEAN NOT NULL DEFAULT FALSE,
  "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
  "updated_at" TIMESTAMP NOT NULL DEFAULT NOW(),

  -- Phone number format validation (E.164 international format)
  CONSTRAINT phone_format CHECK (phone_number ~ '^\+?[1-9]\d{1,14}$')
);


CREATE INDEX  "idx_users_created_at" ON "users"("created_at");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP INDEX IF EXISTS "idx_users_created_at";
DROP TABLE IF EXISTS "users";
DROP TYPE IF EXISTS "gender_type";

-- +goose StatementEnd
