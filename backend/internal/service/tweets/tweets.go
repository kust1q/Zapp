package tweets

import (
	"context"
	"errors"
	"fmt"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
)

var (
	ErrTweetNotFound = errors.New("tweet not found")
	ErrUserNotFound  = errors.New("user not found ")
)

type tweetService struct {
	db     tweetStorage
	media  mediaService
	search searchRepository
}

func NewTweetService(db tweetStorage, media mediaService, search searchRepository) *tweetService {
	return &tweetService{
		db:     db,
		media:  media,
		search: search,
	}
}

func (s *tweetService) BuildEntityTweetToResponse(ctx context.Context, tweet *entity.Tweet) (*entity.Tweet, error) {
	author, err := s.db.GetUserByID(ctx, tweet.Author.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tweet author: %w", err)
	}

	avatarUrl, err := s.media.GetAvatarUrlByUserID(ctx, tweet.Author.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user avatar: %w", err)
	}

	counts, err := s.db.GetCounts(ctx, tweet.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get counters: %w", err)
	}

	return &entity.Tweet{
		ID:            tweet.ID,
		ParentTweetID: tweet.ParentTweetID,
		Content:       tweet.Content,
		CreatedAt:     tweet.CreatedAt,
		UpdatedAt:     tweet.UpdatedAt,
		MediaUrl:      tweet.MediaUrl,
		Author: &entity.SmallUser{
			ID:        author.ID,
			Username:  author.Username,
			AvatarUrl: avatarUrl,
		},
		Counters: counts,
	}, nil
}
