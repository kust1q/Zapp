package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kust1q/Zapp/backend/internal/config"
	httpHandler "github.com/kust1q/Zapp/backend/internal/controllers/http/handler"
	"github.com/kust1q/Zapp/backend/internal/domain/entity"
	s3 "github.com/kust1q/Zapp/backend/internal/providers/db/minio"
	db "github.com/kust1q/Zapp/backend/internal/providers/db/postgres"
	"github.com/kust1q/Zapp/backend/internal/providers/db/redis/cache"
	"github.com/kust1q/Zapp/backend/internal/providers/db/redis/tokens"
	"github.com/kust1q/Zapp/backend/internal/providers/search/elastic"
	"github.com/kust1q/Zapp/backend/internal/service/auth"
	"github.com/kust1q/Zapp/backend/internal/service/feed"
	"github.com/kust1q/Zapp/backend/internal/service/media"
	"github.com/kust1q/Zapp/backend/internal/service/search"
	"github.com/kust1q/Zapp/backend/internal/service/tweets"
	"github.com/kust1q/Zapp/backend/internal/service/user"
	el "github.com/kust1q/Zapp/backend/pkg/elastic"
	mn "github.com/kust1q/Zapp/backend/pkg/minio"
	pg "github.com/kust1q/Zapp/backend/pkg/postgres"
	rs "github.com/kust1q/Zapp/backend/pkg/redis"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	// --- Configs ---
	if err := config.InitConfig(); err != nil {
		logrus.WithError(err).Fatal("error initializing config")
	}
	cfg := config.Get()
	if err := cfg.Validate(); err != nil {
		logrus.WithError(err).Fatal("invalid configuration")
	}

	redisClient, err := rs.NewRedisClient(&cfg.Redis)
	if err != nil {
		logrus.WithError(err).Fatal("failed to initialize redis client")
	}
	cache := cache.NewCache(redisClient, cfg.Cache.DefaultTtl, cfg.Cache.CountersTtl)

	elasticClient, err := el.NewElasticClient([]string{fmt.Sprintf("http://%s:%s", cfg.Elastic.Host, cfg.Elastic.Port)})
	if err != nil {
		logrus.Fatalf("failed to connect to elasticsearch: %v", err)
	}
	elasticRepo := elastic.NewElasticRepository(elasticClient)

	postgresConnect, err := pg.NewPostgresConnect(&cfg.Postgres)
	if err != nil {
		logrus.WithError(err).Fatal("failed to initialize pg connect")
	}
	pgDB := db.NewPostgresDB(postgresConnect, cache)

	minioClient, err := mn.NewMinioClient(&cfg.Minio)
	if err != nil {
		logrus.WithError(err).Fatal("failed to initialize minio client")
	}

	mediaTypeMap := map[entity.MediaType]s3.MediaPolicy{
		entity.MediaTypeImage: {
			MaxSize:     16 * 1024 * 1024,
			AllowedMime: []string{"image/jpeg", "image/png", "image/webp"},
			AllowedExt:  []string{".jpg", ".jpeg", ".png", ".webp", ".bmp", ".tiff"},
		},
		entity.MediaTypeVideo: {
			MaxSize:     512 * 1024 * 1024,
			AllowedMime: []string{"video/mp4", "video/quicktime", "video/x-m4v"},
			AllowedExt:  []string{".mp4", ".mov", ".m4v", ".avi", ".wmv", ".flv", ".webm"},
		},
		entity.MediaTypeGIF: {
			MaxSize:       16 * 1024 * 1024,
			AllowedMime:   []string{"image/gif"},
			AllowedExt:    []string{".gif"},
			ForceMimeType: "image/gif",
		},
		entity.MediaTypeAudio: {
			MaxSize:     64 * 1024 * 1024,
			AllowedMime: []string{"audio/mpeg", "audio/wav", "audio/x-wav", "audio/ogg", "audio/flac", "audio/aac", "audio/x-m4a", "audio/webm"},
			AllowedExt:  []string{".mp3", ".wav", ".ogg", ".flac", ".aac", ".m4a"},
		},
	}

	minioDB := s3.NewMinioDB(minioClient, &cfg.Minio, mediaTypeMap)

	mediaService := media.NewMediaService(pgDB, minioDB)
	authService := auth.NewAuthService(
		config.AuthServiceConfig{PrivateKey: cfg.JWT.PrivateKey, PublicKey: cfg.JWT.PublicKey, AccessTTL: cfg.Tokens.AccessTTL, RefreshTTL: cfg.Tokens.RefreshTTL},
		pgDB,
		mediaService,
		tokens.NewTokenStorage(redisClient),
		elasticRepo)
	tweetService := tweets.NewTweetService(pgDB, mediaService, elasticRepo)
	userService := user.NewUserService(pgDB, mediaService, elasticRepo)
	feedService := feed.NewFeedService(pgDB, tweetService)
	searchService := search.NewSearchService(pgDB, elasticRepo, mediaService)

	handler := httpHandler.NewHandler(
		authService,
		tweetService,
		userService,
		searchService,
		feedService,
		mediaService,
	)

	srv := &http.Server{
		Addr:    ":" + cfg.App.Port,
		Handler: handler.InitRouters(),
	}

	go func() {
		logrus.Infof("Starting server on port %s", cfg.App.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal("Server forced to shutdown:", err)
	}
	logrus.Info("Closing database connection...")
	if err := postgresConnect.Close(); err != nil {
		logrus.Errorf("Error closing DB: %v", err)
	}

	logrus.Info("Closing redis connection...")
	if err := redisClient.Close(); err != nil {
		logrus.Errorf("Error closing Redis: %v", err)
	}

	logrus.Info("Server exiting")
}
