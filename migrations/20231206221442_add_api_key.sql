-- Create "api_keys" table
CREATE TABLE "api_keys" (
  "id" text NOT NULL,
  "created_at" bigint NULL,
  "updated_at" bigint NULL,
  "deleted_at" bigint NULL,
  "name" text NULL,
  "user_id" text NULL,
  "value" text NULL,
  "scopes" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_api_keys_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_api_keys_deleted_at" to table: "api_keys"
CREATE INDEX "idx_api_keys_deleted_at" ON "api_keys" ("deleted_at");
-- Create index "idx_api_keys_name" to table: "api_keys"
CREATE INDEX "idx_api_keys_name" ON "api_keys" ("name");
-- Create index "idx_api_keys_user_id" to table: "api_keys"
CREATE INDEX "idx_api_keys_user_id" ON "api_keys" ("user_id");
-- Modify "licenses" table
ALTER TABLE "licenses" ADD COLUMN "user_id" text NULL, ADD
 CONSTRAINT "fk_licenses_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
