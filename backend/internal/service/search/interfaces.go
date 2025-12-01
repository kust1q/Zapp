package search

import (
	"context"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
)

type searchStorage interface {
	GetTweetsByIDs(ctx context.Context, ids []int) ([]entity.Tweet, error)
	GetUsersByIDs(ctx context.Context, ids []int) ([]entity.User, error)
}

type searchRepository interface {
	SearchTweets(ctx context.Context, query string) ([]int, error)
	SearchUsers(ctx context.Context, query string) ([]int, error)
}
