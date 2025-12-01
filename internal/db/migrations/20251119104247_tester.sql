-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE "testers" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "email" VARCHAR(100) UNIQUE NOT NULL,
  "first_name" VARCHAR(100) NOT NULL,
  "last_name" VARCHAR(100) NOT NULL, 
  "phone_number" VARCHAR(20) UNIQUE,
  "is_verified" BOOLEAN NOT NULL DEFAULT FALSE,
  "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
  "updated_at" TIMESTAMP NOT NULL DEFAULT NOW(),

  -- Phone number format validation (E.164 international format)
  CONSTRAINT phone_format CHECK (phone_number ~ '^\+?[1-9]\d{1,14}$')
);


CREATE INDEX  "idx_testers_created_at" ON "testers"("created_at");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP INDEX IF EXISTS "idx_testers_created_at";
DROP TABLE IF EXISTS "testers";


-- +goose StatementEnd
