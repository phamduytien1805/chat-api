package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ChannelId uuid.UUID

type DirectChannel struct {
	ChannelId    ChannelId `json:"channel_id"`
	FirstUserId  uuid.UUID `json:"user1_id"`
	SecondUserId uuid.UUID `json:"user2_id"`
	CreatedAt    time.Time `json:"created_at"`
}

type UserDirectChannel struct {
	ChannelId ChannelId `json:"channel_id"`
	UserId    uuid.UUID `json:"user_id"`
}

type DirectChannelRepo interface {
	CreateDirectChannel(ctx context.Context, firstUserId, secondUserId uuid.UUID) (DirectChannel, error)
	GetDirectChannelsByIds(ctx context.Context, channelIds []uuid.UUID) ([]DirectChannel, error)
	GetUserDMChannels(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error)
}
