package search

import (
	"context"
	"fmt"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
)

type searchService struct {
	db     searchStorage
	search searchRepository
	media  mediaService
}

func NewSearchService(db searchStorage, search searchRepository, media mediaService) *searchService {
	return &searchService{
		db:     db,
		search: search,
		media:  media,
	}
}

func (s *searchService) SearchTweets(ctx context.Context, query string) ([]entity.Tweet, error) {
	ids, err := s.search.SearchTweets(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search tweets: %w", err)
	}
	if len(ids) == 0 {
		return []entity.Tweet{}, nil
	}

	tweets, err := s.db.GetTweetsByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get tweets by ids: %w", err)
	}

	for i := range tweets {
		if tweets[i].MediaUrl != "" {
			tweets[i].MediaUrl, err = s.media.GetMediaUrlByTweetID(ctx, tweets[i].ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get tweet media by id: %w", err)
			}
		}
	}

	return tweets, nil
}

func (s *searchService) SearchUsers(ctx context.Context, query string) ([]entity.User, error) {
	ids, err := s.search.SearchUsers(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	if len(ids) == 0 {
		return []entity.User{}, nil
	}

	users, err := s.db.GetUsersByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by ids: %w", err)
	}

	for i := range users {
		users[i].AvatarUrl, err = s.media.GetAvatarUrlByUserID(ctx, users[i].ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get user avatar by id: %w", err)
		}
	}

	return users, nil
}
