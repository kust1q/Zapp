package http

import (
	"context"
	"time"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
)

type (
	AuthService interface {
		SignUp(ctx context.Context, req *entity.User) (*entity.User, error)
		SignIn(ctx context.Context, req *entity.Credential) (*entity.Tokens, error)
		Refresh(ctx context.Context, req *entity.Refresh) (*entity.Tokens, error)
		SignOut(ctx context.Context, req *entity.Refresh) error
		VerifyAccessToken(tokenString string) (int, error)
		UpdatePassword(ctx context.Context, req *entity.UpdatePassword) error
		ForgotPassword(ctx context.Context, req *entity.ForgotPassword) (*entity.Recovery, error)
		RecoveryPassword(ctx context.Context, req *entity.RecoveryPassword) error

		GetRefreshTTL() time.Duration
	}

	TweetService interface {
		CreateTweet(ctx context.Context, req *entity.Tweet) (*entity.Tweet, error)
		GetTweetById(ctx context.Context, tweetID int) (*entity.Tweet, error)
		UpdateTweet(ctx context.Context, req *entity.Tweet) (*entity.Tweet, error)
		LikeTweet(ctx context.Context, userID, tweetID int) error
		UnlikeTweet(ctx context.Context, userID, tweetID int) error
		//ReplyToTweet(ctx context.Context, req *entity.Tweet) (*entity.Tweet, error)
		GetRepliesToTweet(ctx context.Context, tweetID int) ([]entity.Tweet, error)
		CreateRetweet(ctx context.Context, userID, tweetID int) error
		DeleteRetweet(ctx context.Context, userID, retweetID int) error
		GetTweetsAndRetweetsByUsername(ctx context.Context, username string) ([]entity.Tweet, error)
		GetLikes(ctx context.Context, tweetID int) ([]entity.SmallUser, error)
		DeleteTweet(ctx context.Context, userID, tweetID int) error
	}

	UserService interface {
		Update(ctx context.Context, req *entity.UpdateBio) error
		FollowToUser(ctx context.Context, followerID, followingID int) (*entity.Follow, error)
		UnfollowUser(ctx context.Context, followerID, followingID int) error
		GetFollowers(ctx context.Context, username string) ([]entity.SmallUser, error)
		GetFollowings(ctx context.Context, username string) ([]entity.SmallUser, error)
		GetUserProfile(ctx context.Context, username string) (*entity.UserProfile, error)
		GetMe(ctx context.Context, userID int) (*entity.UserProfile, error)
		DeleteUser(ctx context.Context, userID int) error
	}

	SearchService interface {
	}

	FeedService interface {
		GetUserFeedByUserId(ctx context.Context, userID int) ([]entity.Tweet, error)
	}

	MediaService interface {
		GetMediaDataByTweetID(ctx context.Context, tweetID int) (*entity.TweetMedia, error)
		GetAvatarDataByUserID(ctx context.Context, userID int) (*entity.Avatar, error)
		DeleteTweetMedia(ctx context.Context, tweetID, userID int) error
	}
)
