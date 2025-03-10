// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: dm_channel.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createDMChannel = `-- name: CreateDMChannel :one
INSERT INTO dm_channels (
  channel_id,
  user1_id,
  user2_id
) VALUES (
  $1, $2, $3
) RETURNING channel_id, user1_id, user2_id, created_at
`

type CreateDMChannelParams struct {
	ChannelID uuid.UUID `json:"channel_id"`
	User1ID   uuid.UUID `json:"user1_id"`
	User2ID   uuid.UUID `json:"user2_id"`
}

func (q *Queries) CreateDMChannel(ctx context.Context, arg CreateDMChannelParams) (DmChannel, error) {
	row := q.db.QueryRow(ctx, createDMChannel, arg.ChannelID, arg.User1ID, arg.User2ID)
	var i DmChannel
	err := row.Scan(
		&i.ChannelID,
		&i.User1ID,
		&i.User2ID,
		&i.CreatedAt,
	)
	return i, err
}

const createUserDMChannel = `-- name: CreateUserDMChannel :one
INSERT INTO user_dm_channels (
  user_id,
  channel_id
) VALUES (
  $1, $2
) RETURNING user_id, channel_id
`

type CreateUserDMChannelParams struct {
	UserID    uuid.UUID `json:"user_id"`
	ChannelID uuid.UUID `json:"channel_id"`
}

func (q *Queries) CreateUserDMChannel(ctx context.Context, arg CreateUserDMChannelParams) (UserDmChannel, error) {
	row := q.db.QueryRow(ctx, createUserDMChannel, arg.UserID, arg.ChannelID)
	var i UserDmChannel
	err := row.Scan(&i.UserID, &i.ChannelID)
	return i, err
}

const getDMChannelByChannelId = `-- name: GetDMChannelByChannelId :one
Select channel_id, user1_id, user2_id, created_at from dm_channels dc where dc.channel_id = $1
`

func (q *Queries) GetDMChannelByChannelId(ctx context.Context, channelID uuid.UUID) (DmChannel, error) {
	row := q.db.QueryRow(ctx, getDMChannelByChannelId, channelID)
	var i DmChannel
	err := row.Scan(
		&i.ChannelID,
		&i.User1ID,
		&i.User2ID,
		&i.CreatedAt,
	)
	return i, err
}

const getDMChannelByChannelIds = `-- name: GetDMChannelByChannelIds :many
Select channel_id, user1_id, user2_id, created_at from dm_channels dc where dc.channel_id = ANY($1::uuid[])
`

func (q *Queries) GetDMChannelByChannelIds(ctx context.Context, dollar_1 []uuid.UUID) ([]DmChannel, error) {
	rows, err := q.db.Query(ctx, getDMChannelByChannelIds, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []DmChannel{}
	for rows.Next() {
		var i DmChannel
		if err := rows.Scan(
			&i.ChannelID,
			&i.User1ID,
			&i.User2ID,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserDMChannels = `-- name: GetUserDMChannels :many
Select udc.channel_id from user_dm_channels udc where udc.user_id = $1
`

func (q *Queries) GetUserDMChannels(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := q.db.Query(ctx, getUserDMChannels, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []uuid.UUID{}
	for rows.Next() {
		var channel_id uuid.UUID
		if err := rows.Scan(&channel_id); err != nil {
			return nil, err
		}
		items = append(items, channel_id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
