package feed

import (
	"context"
	"fmt"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
)

type userStorage interface {
	GetUserByID(ctx context.Context, userID int) (*entity.User, error)
	GetFollowingsIds(ctx context.Context, username string) ([]int, error)
}

type feedStorage interface {
	GetFeedByIds(ctx context.Context, userIDs []int) ([]entity.Tweet, error)
}

type tweetService interface {
	TweetResponseByTweet(ctx context.Context, tweet *entity.Tweet) (*entity.Tweet, error)
}

type feedService struct {
	tweets tweetService
	users  userStorage
	feed   feedStorage
}

func NewFeedService(tweets tweetService, users userStorage, feed feedStorage) *feedService {
	return &feedService{
		tweets: tweets,
		users:  users,
		feed:   feed,
	}
}

func (s *feedService) GetUserFeedByUserId(ctx context.Context, userID int) ([]entity.Tweet, error) {
	user, err := s.users.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	followings, err := s.users.GetFollowingsIds(ctx, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user followings by username: %w", err)
	}
	feed, err := s.feed.GetFeedByIds(ctx, followings)
	if err != nil {
		return nil, fmt.Errorf("failed to get feed by ids: %w", err)
	}
	res := make([]entity.Tweet, 0, len(feed))
	for _, t := range feed {
		tr, err := s.tweets.TweetResponseByTweet(ctx, &t)
		if err != nil {
			return nil, fmt.Errorf("failed to change tweet entity to response: %w", err)
		}
		res = append(res, *tr)
	}
	return res, nil
}
