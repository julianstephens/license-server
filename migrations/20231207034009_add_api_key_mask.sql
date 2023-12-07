-- Modify "api_keys" table
ALTER TABLE "api_keys" ADD COLUMN "mask" text NULL;
-- Create index "idx_api_keys_mask" to table: "api_keys"
CREATE INDEX "idx_api_keys_mask" ON "api_keys" ("mask");
