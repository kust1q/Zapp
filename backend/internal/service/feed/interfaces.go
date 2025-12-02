package feed

import (
	"context"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
)

type db interface {
	GetUserByID(ctx context.Context, userID int) (*entity.User, error)
	GetFollowingsIds(ctx context.Context, username string) ([]int, error)

	GetFeedByIds(ctx context.Context, userIDs []int) ([]entity.Tweet, error)
}

type tweetService interface {
	BuildEntityTweetToResponse(ctx context.Context, tweet *entity.Tweet) (*entity.Tweet, error)
}
