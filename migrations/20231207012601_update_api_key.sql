-- Modify "api_keys" table
ALTER TABLE "api_keys" DROP COLUMN "name", DROP COLUMN "value", ADD COLUMN "key" bytea NULL;
