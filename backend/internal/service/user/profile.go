package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
	"github.com/kust1q/Zapp/backend/internal/dto"
)

func (s *userService) GetMe(ctx context.Context, userID int) (*dto.UserProfileResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	user, err := s.storage.GetUserByID(ctx, userID)
	if err != nil {
		return &dto.UserProfileResponse{}, fmt.Errorf("failed to get user by id: %w", err)
	}
	return s.GetUserProfile(ctx, user.Username)
}

func (s *userService) GetUserProfile(ctx context.Context, username string) (*dto.UserProfileResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	user, err := s.storage.GetUserByUsername(ctx, username)
	if err != nil {
		return &dto.UserProfileResponse{}, fmt.Errorf("failed to get user by user id: %w", err)
	}
	avatar, err := s.media.GetAvatarByUserID(ctx, user.ID)
	if err != nil {
		return &dto.UserProfileResponse{}, fmt.Errorf("failed to get avatar by user id: %w", err)
	}
	tweets, err := s.storage.GetTweetsAndRetweetsByUsername(ctx, user.Username)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return &dto.UserProfileResponse{}, fmt.Errorf("failed to get tweets by username: %w", err)
		}
	}

	tweetsRes := make([]dto.TweetResponse, 0, len(tweets))
	for _, t := range tweets {
		tr, err := s.tweetResponseByTweet(ctx, &t, user, avatar)
		if err != nil {
			return &dto.UserProfileResponse{}, fmt.Errorf("failed to get tweet responses by tweet: %w", err)
		}
		tweetsRes = append(tweetsRes, *tr)
	}
	return &dto.UserProfileResponse{
		UserResponse: dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Bio:       user.Bio,
			Gen:       user.Gen,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			Avatar: dto.Avatar{
				MediaURL:  avatar.MediaURL,
				MimeType:  avatar.MimeType,
				SizeBytes: avatar.SizeBytes,
			},
		},
		Tweets: tweetsRes,
	}, nil
}

func (s *userService) tweetResponseByTweet(ctx context.Context, tweet *entity.Tweet, author *entity.User, avatar *dto.Avatar) (*dto.TweetResponse, error) {
	var m dto.TweetMedia
	media, err := s.media.GetMediaByTweetID(ctx, tweet.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			m = dto.TweetMedia{}
		} else {
			return &dto.TweetResponse{}, fmt.Errorf("failed to get tweet by id: %w", err)
		}
	} else {
		m = dto.TweetMedia{
			ID:        media.ID,
			TweetID:   media.TweetID,
			MediaURL:  media.MediaURL,
			MimeType:  media.MimeType,
			SizeBytes: media.SizeBytes,
		}
	}

	return &dto.TweetResponse{
		ID:            tweet.ID,
		Content:       tweet.Content,
		CreatedAt:     tweet.CreatedAt,
		UpdatedAt:     tweet.UpdatedAt,
		ParentTweetID: &tweet.ParentTweetID,
		Media:         m,
		Author: dto.SmallUserResponse{
			ID:       author.ID,
			Username: author.Username,
			Avatar: dto.Avatar{
				MediaURL:  avatar.MediaURL,
				MimeType:  avatar.MimeType,
				SizeBytes: avatar.SizeBytes,
			},
		},
	}, nil
}
