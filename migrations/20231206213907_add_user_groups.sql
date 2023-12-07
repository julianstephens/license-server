-- Create "user_groups" table
CREATE TABLE "user_groups" (
  "id" text NOT NULL,
  "created_at" bigint NULL,
  "updated_at" bigint NULL,
  "deleted_at" bigint NULL,
  "name" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_user_groups_deleted_at" to table: "user_groups"
CREATE INDEX "idx_user_groups_deleted_at" ON "user_groups" ("deleted_at");
-- Create index "idx_user_groups_name" to table: "user_groups"
CREATE UNIQUE INDEX "idx_user_groups_name" ON "user_groups" ("name");
-- Modify "licenses" table
ALTER TABLE "licenses" ADD
 CONSTRAINT "fk_products_licenses" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Create "product_rules" table
CREATE TABLE "product_rules" (
  "rule_id" text NOT NULL,
  "product_id" text NOT NULL,
  PRIMARY KEY ("rule_id", "product_id"),
  CONSTRAINT "fk_product_rules_product" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_product_rules_rule" FOREIGN KEY ("rule_id") REFERENCES "rules" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "user_rules" table
CREATE TABLE "user_rules" (
  "rule_id" text NOT NULL,
  "user_id" text NOT NULL,
  PRIMARY KEY ("rule_id", "user_id"),
  CONSTRAINT "fk_user_rules_rule" FOREIGN KEY ("rule_id") REFERENCES "rules" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_user_rules_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
