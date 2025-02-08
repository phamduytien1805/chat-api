package db

import "github.com/phamduytien1805/hub/domain"

func mapToDMChannel(u DmChannel) *domain.DirectChannel {
	return &domain.DirectChannel{
		ChannelId:    domain.ChannelId(u.ChannelID),
		FirstUserId:  u.User1ID,
		SecondUserId: u.User2ID,
		CreatedAt:    u.CreatedAt,
	}
}
