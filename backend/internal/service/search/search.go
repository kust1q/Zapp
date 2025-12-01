package search

import (
	"context"
	"fmt"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
)

type searchService struct {
	db     searchStorage
	search searchRepository
}

func NewSearchService(search searchRepository, db searchStorage) *searchService {
	return &searchService{
		db:     db,
		search: search,
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

	return s.db.GetTweetsByIDs(ctx, ids)
}

func (s *searchService) SearchUsers(ctx context.Context, query string) ([]entity.User, error) {
	ids, err := s.search.SearchUsers(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	if len(ids) == 0 {
		return []entity.User{}, nil
	}

	return s.db.GetUsersByIDs(ctx, ids)
}
