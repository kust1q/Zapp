package feed

import (
	"context"
	"fmt"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
)

type feedService struct {
	DB           db
	tweetService tweetService
}

func NewFeedService(db db, tweetService tweetService) *feedService {
	return &feedService{
		DB:           db,
		tweetService: tweetService,
	}
}

func (s *feedService) GetUserFeedByUserId(ctx context.Context, userID int) ([]entity.Tweet, error) {
	user, err := s.DB.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	followings, err := s.DB.GetFollowingsIds(ctx, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user followings by username: %w", err)
	}
	feed, err := s.DB.GetFeedByIds(ctx, followings)
	if err != nil {
		return nil, fmt.Errorf("failed to get feed by ids: %w", err)
	}
	res := make([]entity.Tweet, 0, len(feed))
	for _, t := range feed {
		tr, err := s.tweetService.BuildEntityTweetToResponse(ctx, &t)
		if err != nil {
			return nil, fmt.Errorf("failed to change tweet entity to response: %w", err)
		}
		res = append(res, *tr)
	}
	return res, nil
}
