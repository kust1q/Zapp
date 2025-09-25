package http

import (
	"context"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kust1q/Zapp/backend/internal/dto"
)

type AuthService interface {
	SignUp(ctx context.Context, user *dto.SignUpRequest) (*dto.SignUpResponse, error)
	SignIn(ctx context.Context, credential *dto.SignInRequest) (*dto.SignInResponse, error)
	Refresh(ctx context.Context, token *dto.RefreshRequest) (*dto.SignInResponse, error)
	SignOut(ctx context.Context, token *dto.RefreshRequest) error
	VerifyAccessToken(tokenString string) (int, error)
	UpdateSecuritySettings(ctx context.Context, userID int, req *dto.UpdateSecuritySettingsRequest) error
	ResetPassword(ctx context.Context, userID int, req *dto.ResetPasswordRequest) error
	ForgotPassword(ctx context.Context, req *dto.ForgotPasswordRequest) error
}

type TweetService interface {
	CreateTweet(ctx context.Context, userID int, tweet *dto.CreateTweetRequest) (*dto.TweetResponse, error)
	CreateTweetWithMedia(ctx context.Context, userID int, tweet *dto.CreateTweetRequest, file *dto.FileData) (*dto.TweetResponse, error)
	GetTweetById(ctx context.Context, tweetID int) (*dto.TweetResponseWithCounters, error)
	UpdateTweet(ctx context.Context, userID, tweetID int, req *dto.UpdateTweetRequest) (*dto.UpdateTweetResponse, error)
	LikeTweet(ctx context.Context, userID, tweetID int) error
	UnlikeTweet(ctx context.Context, userID, tweetID int) error
	ReplyToTweet(ctx context.Context, userID, tweetID int, tweet *dto.CreateTweetRequest) (*dto.TweetResponse, error)
	ReplyToTweetWithMedia(ctx context.Context, userID, tweetID int, tweet *dto.CreateTweetRequest, file *dto.FileData) (*dto.TweetResponse, error)
	GetRepliesToTweet(ctx context.Context, tweetID int) ([]dto.TweetResponse, error)
	CreateRetweet(ctx context.Context, userID, tweetID int) error
	DeleteRetweet(ctx context.Context, userID, retweetID int) error
	GetTweetsAndRetweetsByUsername(ctx context.Context, username string) ([]dto.TweetResponse, error)
	GetLikes(ctx context.Context, tweetID int) ([]dto.UserLikeResponse, error)
	DeleteTweet(ctx context.Context, userID, tweetID int) error
}

type UserService interface {
	Update(ctx context.Context, userID int, req *dto.UpdateBioRequest) error
	FollowToUser(ctx context.Context, followerID, followingID int) (*dto.FollowResponse, error)
	UnfollowUser(ctx context.Context, followerID, followingID int) error
	GetFollowers(ctx context.Context, username string) ([]dto.SmallUserResponse, error)
	GetFollowings(ctx context.Context, username string) ([]dto.SmallUserResponse, error)
	GetUserProfile(ctx context.Context, username string) (*dto.UserProfileResponse, error)
	GetMe(ctx context.Context, userID int) (*dto.UserProfileResponse, error)
	DeleteUser(ctx context.Context, userID int) error
}

type SearchService interface {
}

type FeedService interface {
}

type MediaService interface {
	GetMediaDataByTweetID(ctx context.Context, tweetID int) (*dto.TweetMediaDataResponse, error)
	GetAvatarDataByUserID(ctx context.Context, userID int) (*dto.AvatarDataResponse, error)
	DeleteTweetMedia(ctx context.Context, tweetID, userID int) error
}

type Handler struct {
	authService   AuthService
	tweetService  TweetService
	userService   UserService
	searchService SearchService
	feedService   FeedService
	mediaService  MediaService
}

func NewHandler(
	authService AuthService,
	tweetService TweetService,
	userService UserService,
	searchService SearchService,
	feedService FeedService,
	mediaService MediaService,
) *Handler {
	return &Handler{
		authService:   authService,
		tweetService:  tweetService,
		userService:   userService,
		searchService: searchService,
		feedService:   feedService,
		mediaService:  mediaService,
	}
}

func (h *Handler) InitRouters() *gin.Engine {
	router := gin.New()

	router.Use(cors.Default())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api/v1")

	auth := api.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/refresh", h.refresh)
		auth.POST("/sign-out", h.signOut)
		auth.POST("/forgot-password", h.forgotPassword)
	}

	public := api.Group("/public")
	{
		public.GET("/tweets/:tweet_id", h.getTweetById)
		public.GET("/tweets/:tweet_id/replies", h.getReplies)
		public.GET("/tweets/:tweet_id/likes", h.getLikes)
		public.GET("/tweets/media/:tweet_id", h.getTweetMedia)
		public.GET("/users/:username/profile", h.getUserProfile)
		public.GET("/users/:username/tweets", h.getTweetsAndRetweetsByUsername)
		public.GET("/users/:username/followers", h.followers)
		public.GET("/users/:username/following", h.following)
		public.GET("/users/avatar/:user_id", h.getAvatar)
		public.GET("/search", h.search)
	}

	protected := api.Group("/protected", h.authMiddleware)
	{
		protected.PUT("/reset-password", h.resetPassword)
		protected.PUT("/security-settings", h.updateSecuritySettings)

		tweets := protected.Group("/tweets")
		{
			tweets.POST("", h.createTweet)
			tweets.PUT("/:tweet_id", h.updateTweet)
			tweets.DELETE("/:tweet_id", h.deleteTweet)

			tweets.POST("/:tweet_id/like", h.likeTweet)
			tweets.DELETE("/:tweet_id/like", h.unlikeTweet)

			tweets.POST("/:tweet_id/reply", h.replyToTweet)

			tweets.POST("/:tweet_id/retweet", h.retweet)
			tweets.DELETE("/retweets/:tweet_id", h.deleteRetweet)

			tweets.DELETE("tweets/:tweet_id", h.deleteTweetMedia)
		}

		users := protected.Group("/users")
		{
			users.GET("/me", h.getMe)
			users.PUT("/me", h.updateMe)
			users.DELETE("/me", h.deleteMe)
			users.POST("/:user_id/follow", h.followUser)
			users.DELETE("/:user_id/follow", h.unfollowUser)
		}

		feed := protected.Group("/feed")
		{
			feed.GET("/", h.getFeed)
			feed.GET("/explore", h.exploreFeed)
			feed.GET("/trends", h.trendsFeed)
		}
	}

	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Endpoint not found"})
	})

	return router
}
