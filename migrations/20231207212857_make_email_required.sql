-- Modify "users" table
ALTER TABLE "users" ALTER COLUMN "email" SET NOT NULL;
-- Create index "users_email_key" to table: "users"
CREATE UNIQUE INDEX "users_email_key" ON "users" ("email");
