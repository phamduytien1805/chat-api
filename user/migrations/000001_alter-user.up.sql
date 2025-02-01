CREATE TABLE "users" (
  "id" uuid PRIMARY KEY,
  "username" VARCHAR(32) NOT NULL UNIQUE,
  "email" VARCHAR(254) NOT NULL UNIQUE,
  "email_verified" bool NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "user_credentials" (
  "user_id" uuid NOT NULL,
  "hashed_password" varchar NOT NULL
);


ALTER TABLE "user_credentials" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
