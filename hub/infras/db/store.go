package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/phamduytien1805/hub/domain"
)

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	connPool *pgxpool.Pool
	q        *Queries
}

// NewStore creates a new store
func NewStore(connPool *pgxpool.Pool) domain.DirectChannelRepo {
	return &SQLStore{
		connPool: connPool,
		q:        New(connPool),
	}
}

// CreateDirectChannel creates a new direct message channel between two distinct users
func (s *SQLStore) CreateDirectChannel(ctx context.Context, userID1, userID2 uuid.UUID) (domain.DirectChannel, error) {
	// Validate distinct users
	if userID1 == userID2 {
		return domain.DirectChannel{}, fmt.Errorf("cannot create direct channel between identical users")
	}
	userID1.Time()

	// Ensure consistent ordering of user IDs to prevent duplicate channels
	if userID1.Time() > userID2.Time() {
		userID1, userID2 = userID2, userID1
	}

	channelID, err := uuid.NewV7()
	if err != nil {
		return domain.DirectChannel{}, fmt.Errorf("failed to generate channel UUID: %w", err)
	}

	arg := CreateDMChannelParams{
		ChannelID: channelID,
		User1ID:   userID1,
		User2ID:   userID2,
	}

	dmChannel, err := s.q.CreateDMChannel(ctx, arg)
	if err != nil {
		return domain.DirectChannel{}, fmt.Errorf("failed to persist direct channel: %w", err)
	}

	return mapToDMChannel(dmChannel), nil
}

// GetDirectChannelsById returns a list of direct channels by their ids
func (s *SQLStore) GetDirectChannelsById(ctx context.Context, channelIds []uuid.UUID) ([]domain.DirectChannel, error) {
	dmChannels, err := s.q.GetDMChannelByChannelIds(ctx, channelIds)
	if err != nil {
		return nil, err
	}

	// Pre-allocate slice with exact length for direct assignment
	result := make([]domain.DirectChannel, len(dmChannels))

	for i, dmChannel := range dmChannels {
		// Direct assignment avoids append overhead and pointer dereferencing
		result[i] = mapToDMChannel(dmChannel)
	}

	return result, nil
}

// GetDirectChannelIdsByUserId returns a list of direct channel ids by user id
func (s *SQLStore) GetDirectChannelIdsByUserId(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	return s.q.GetUserDMChannels(ctx, userID)
}
