package usecase

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/phamduytien1805/hub/domain"
)

type CreateDMChannelUsecase struct {
	repo   domain.DirectChannelRepo
	logger *slog.Logger
}

func NewCreateDMChannelUsecase(logger *slog.Logger, dmRepo domain.DirectChannelRepo) *CreateDMChannelUsecase {
	return &CreateDMChannelUsecase{
		repo:   dmRepo,
		logger: logger,
	}
}

func (s *CreateDMChannelUsecase) Exec(ctx context.Context, user1ID, user2ID uuid.UUID) (domain.DirectChannel, error) {
	return s.repo.CreateDirectChannel(ctx, user1ID, user2ID)
}
