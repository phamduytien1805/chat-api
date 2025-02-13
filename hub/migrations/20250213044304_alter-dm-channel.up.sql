CREATE TABLE "dm_channels" (
  "channel_id" uuid PRIMARY KEY,
  "user1_id" uuid NOT NULL,
  "user2_id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  UNIQUE (user1_id, user2_id)
);

CREATE TABLE "user_dm_channels" (
  "user_id" uuid NOT NULL,
  "channel_id" uuid NOT NULL,
  PRIMARY KEY (user_id, channel_id)
);

ALTER TABLE "user_dm_channels" ADD FOREIGN KEY ("channel_id") REFERENCES "dm_channels" ("channel_id");

