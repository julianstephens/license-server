-- Create "group_members" table
CREATE TABLE "group_members" (
  "user_group_id" text NOT NULL,
  "user_id" text NOT NULL,
  PRIMARY KEY ("user_group_id", "user_id"),
  CONSTRAINT "fk_group_members_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_group_members_user_group" FOREIGN KEY ("user_group_id") REFERENCES "user_groups" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "group_rules" table
CREATE TABLE "group_rules" (
  "user_group_id" text NOT NULL,
  "rule_id" text NOT NULL,
  PRIMARY KEY ("user_group_id", "rule_id"),
  CONSTRAINT "fk_group_rules_rule" FOREIGN KEY ("rule_id") REFERENCES "rules" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_group_rules_user_group" FOREIGN KEY ("user_group_id") REFERENCES "user_groups" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "groups" table
CREATE TABLE "groups" (
  "user_id" text NOT NULL,
  "user_group_id" text NOT NULL,
  PRIMARY KEY ("user_id", "user_group_id"),
  CONSTRAINT "fk_groups_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_groups_user_group" FOREIGN KEY ("user_group_id") REFERENCES "user_groups" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
