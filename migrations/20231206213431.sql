-- Create "licenses" table
CREATE TABLE "licenses" (
  "id" text NOT NULL,
  "created_at" bigint NULL,
  "updated_at" bigint NULL,
  "deleted_at" bigint NULL,
  "product_id" text NULL,
  "value" text NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_licenses_deleted_at" to table: "licenses"
CREATE INDEX "idx_licenses_deleted_at" ON "licenses" ("deleted_at");
-- Create index "licenses_value_key" to table: "licenses"
CREATE UNIQUE INDEX "licenses_value_key" ON "licenses" ("value");
-- Create "products" table
CREATE TABLE "products" (
  "id" text NOT NULL,
  "created_at" bigint NULL,
  "updated_at" bigint NULL,
  "deleted_at" bigint NULL,
  "name" text NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_products_deleted_at" to table: "products"
CREATE INDEX "idx_products_deleted_at" ON "products" ("deleted_at");
-- Create index "idx_products_name" to table: "products"
CREATE UNIQUE INDEX "idx_products_name" ON "products" ("name");
-- Create "rules" table
CREATE TABLE "rules" (
  "id" text NOT NULL,
  "created_at" bigint NULL,
  "updated_at" bigint NULL,
  "deleted_at" bigint NULL,
  "name" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_rules_deleted_at" to table: "rules"
CREATE INDEX "idx_rules_deleted_at" ON "rules" ("deleted_at");
-- Create index "idx_rules_name" to table: "rules"
CREATE UNIQUE INDEX "idx_rules_name" ON "rules" ("name");
-- Create "users" table
CREATE TABLE "users" (
  "id" text NOT NULL,
  "created_at" bigint NULL,
  "updated_at" bigint NULL,
  "deleted_at" bigint NULL,
  "name" text NULL,
  "email" text NULL,
  "password" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "users" ("deleted_at");
-- Create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX "idx_users_email" ON "users" ("email");
-- Create index "idx_users_password" to table: "users"
CREATE UNIQUE INDEX "idx_users_password" ON "users" ("password");
