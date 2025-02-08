CREATE TABLE "dm_channels" (
  "channel_id" uuid PRIMARY KEY,
  "user1_id" uuid NOT NULL,
  "user2_id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  UNIQUE (user1_id, user2_id)
);

ALTER TABLE "dm_channels" ADD FOREIGN KEY ("user1_id") REFERENCES "users" ("id");
ALTER TABLE "dm_channels" ADD FOREIGN KEY ("user2_id") REFERENCES "users" ("id");

CREATE TABLE "user_dm_channels" (
  "user_id" uuid NOT NULL,
  "channel_id" uuid NOT NULL,
  PRIMARY KEY (user_id, thread_id)
);


ALTER TABLE "user_dm_channels" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "user_dm_channels" ADD FOREIGN KEY ("channel_id") REFERENCES "dm_channels" ("channel_id");

