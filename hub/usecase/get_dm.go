package usecase

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/phamduytien1805/hub/domain"
)

type GetDMChannelUsecase struct {
	repo   domain.DirectChannelRepo
	logger *slog.Logger
}

func NewGetDMChannelUsecase(logger *slog.Logger, dmRepo domain.DirectChannelRepo) *GetDMChannelUsecase {
	return &GetDMChannelUsecase{
		repo:   dmRepo,
		logger: logger,
	}
}

func (s *GetDMChannelUsecase) ById(ctx context.Context, channelId uuid.UUID) (domain.DirectChannel, error) {
	channels, err := s.repo.GetDirectChannelsByIds(ctx, []uuid.UUID{channelId})
	if err != nil {
		return domain.DirectChannel{}, err
	}
	return channels[0], nil
}

func (s *GetDMChannelUsecase) ByIds(ctx context.Context, channelIds []uuid.UUID) ([]domain.DirectChannel, error) {
	return s.repo.GetDirectChannelsByIds(ctx, channelIds)
}

func (s *GetDMChannelUsecase) ByUserId(ctx context.Context, userId uuid.UUID) ([]domain.DirectChannel, error) {
	channelIds, err := s.repo.GetUserDMChannels(ctx, userId)
	if err != nil {
		return nil, err
	}
	return s.ByIds(ctx, channelIds)
}
