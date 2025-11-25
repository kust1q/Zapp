package tweets

import (
	"context"
	"fmt"
	"time"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
)

/*
	func (s *tweetService) CreateTweet(ctx context.Context, userID int, tweet *dto.CreateTweetRequest) (*dto.TweetResponse, error) {
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		domainTweet := entity.Tweet{
			UserID:    userID,
			Content:   tweet.Content,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		createdTweet, err := s.db.CreateTweet(ctx, &domainTweet)
		if err != nil {
			return &dto.TweetResponse{}, fmt.Errorf("tweet creation failed: %w", err)
		}

		author, err := s.db.GetUserByID(ctx, userID)
		if err != nil {
			return &dto.TweetResponse{}, fmt.Errorf("failed to get tweet author")
		}

		avatar, err := s.media.GetAvatarUrlByUserID(ctx, userID)
		if err != nil {
			return &dto.TweetResponse{}, fmt.Errorf("failed to get user avatar")
		}

		return &dto.TweetResponse{
			ID:        createdTweet.ID,
			Content:   createdTweet.Content,
			CreatedAt: createdTweet.CreatedAt,
			UpdatedAt: createdTweet.UpdatedAt,
			Author: dto.SmallUserResponse{
				ID:        author.ID,
				Username:  author.Username,
				AvatarURL: avatar,
			},
		}, nil
	}
*/
func (s *tweetService) CreateTweet(ctx context.Context, tweet *entity.Tweet) (*entity.Tweet, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	createdTweet, err := s.db.CreateTweetTx(ctx, tx, tweet)
	if err != nil {
		return nil, fmt.Errorf("user creation failed: %w", err)
	}

	var mediaUrl string
	if tweet.File != nil {
		mediaUrl, err = s.media.UploadAndAttachTweetMediaTx(ctx, createdTweet.ID, tweet.File.File, tweet.File.Header.Filename, tx)
		if err != nil {
			return nil, fmt.Errorf("failed to upload media: %w", err)
		}
	}

	author, err := s.db.GetUserByID(ctx, tweet.Author.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tweet author: %w", err)
	}

	avatarUrl, err := s.media.GetAvatarUrlByUserID(ctx, tweet.Author.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user avatar: %w", err)
	}

	likes, retweets, replies, err := s.db.GetCounts(ctx, createdTweet.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get counters: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction failed: %w", err)
	}

	return &entity.Tweet{
		ID:            createdTweet.ID,
		ParentTweetID: createdTweet.ParentTweetID,
		Content:       createdTweet.Content,
		CreatedAt:     createdTweet.CreatedAt,
		UpdatedAt:     createdTweet.UpdatedAt,
		MediaUrl:      mediaUrl,
		Author: &entity.SmallUser{
			ID:        author.ID,
			Username:  author.Username,
			AvatarURL: avatarUrl,
		},
		Counters: &entity.Counters{
			ReplyCount:   replies,
			RetweetCount: retweets,
			LikeCount:    likes,
		},
	}, nil
}
