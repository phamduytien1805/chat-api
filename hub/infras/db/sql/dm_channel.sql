-- name: CreateDMChannel :one
INSERT INTO dm_channels (
  channel_id,
  user1_id,
  user2_id
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: CreateUserDMChannel :one
INSERT INTO user_dm_channels (
  user_id,
  channel_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetDMChannelByChannelId :one
Select * from dm_channels dc where dc.channel_id = $1;

-- name: GetDMChannelByChannelIds :many
Select * from dm_channels dc where dc.channel_id = ANY($1::uuid[]);

-- name: GetUserDMChannels :many
Select udc.channel_id from user_dm_channels udc where udc.user_id = $1;
