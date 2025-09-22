package main

import (
	"crypto/rsa"
	"fmt"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kust1q/Zapp/backend/internal/config"
	"github.com/kust1q/Zapp/backend/internal/delivery/http"
	"github.com/kust1q/Zapp/backend/internal/security"
	"github.com/kust1q/Zapp/backend/internal/servers"
	"github.com/kust1q/Zapp/backend/internal/service/auth"
	"github.com/kust1q/Zapp/backend/internal/service/media"
	"github.com/kust1q/Zapp/backend/internal/service/tweets"
	"github.com/kust1q/Zapp/backend/internal/service/user"
	"github.com/kust1q/Zapp/backend/internal/storage/cache"
	"github.com/kust1q/Zapp/backend/internal/storage/data"
	"github.com/kust1q/Zapp/backend/internal/storage/minio"
	"github.com/kust1q/Zapp/backend/internal/storage/objects"
	"github.com/kust1q/Zapp/backend/internal/storage/postgres"
	"github.com/kust1q/Zapp/backend/internal/storage/redis"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	//Configs
	if err := config.InitConfig(); err != nil {
		logrus.WithError(err).Fatal("Error initializing config")
	}

	cfg := config.Get()

	if err := cfg.Validate(); err != nil {
		logrus.WithError(err).Fatal("Invalid configuration")
	}
	//DBs
	postgres, err := postgres.NewPostgresDB(cfg.Postgres)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize database")
	}
	defer func() {
		if err := postgres.Close(); err != nil {
			logrus.WithError(err).Error("Failed to close database connection")
		}
	}()

	minio, err := minio.NewMinioClient(cfg.Minio)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize minio client")
	}

	redis, err := redis.NewRedisClient(cfg.Redis)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize redis client")
	}
	defer func() {
		if err := redis.Close(); err != nil {
			logrus.WithError(err).Error("Failed to close Redis connection")
		}
	}()

	hasher := security.NewHasher(cfg.Cache.HashSecret)

	privateKey, publicKey, err := loadRSAKeys(cfg.JWT.PrivateKeyPath, cfg.JWT.PublicKeyPath)
	if err != nil {
		log.Fatalf("Failed to load RSA keys: %v", err)
	}

	//Init storage
	userCache := cache.NewAuthCache(redis, hasher, cfg.Cache.TTL)
	dataStorage := data.NewDataStorage(postgres, userCache)
	mediaTypeMap := map[objects.MediaType]objects.MediaTypeConfig{
		objects.TypeImage: {
			MaxSize:     16 * 1024 * 1024, // 16MB
			AllowedMime: []string{"image/jpeg", "image/png", "image/webp"},
			AllowedExt:  []string{".jpg", ".jpeg", ".png", ".webp", ".bmp", ".tiff"},
		},
		objects.TypeVideo: {
			MaxSize:     512 * 1024 * 1024, // 512 MB
			AllowedMime: []string{"video/mp4", "video/quicktime", "video/x-m4v"},
			AllowedExt:  []string{".mp4", ".mov", ".m4v", ".avi", ".wmv", ".flv", ".webm"},
		},
		objects.TypeGIF: {
			MaxSize:       16 * 1024 * 1024, // 16MB
			AllowedMime:   []string{"image/gif"},
			AllowedExt:    []string{".gif"},
			ForceMimeType: "image/gif",
		},
		objects.TypeAudio: {
			MaxSize: 64 * 1024 * 1024, // 64 MB
			AllowedMime: []string{
				"audio/mpeg",
				"audio/wav",
				"audio/x-wav",
				"audio/ogg",
				"audio/flac",
				"audio/aac",
				"audio/x-m4a",
				"audio/webm",
			},
			AllowedExt: []string{".mp3", ".wav", ".ogg", ".flac", ".aac", ".m4a"},
		},
	}
	objectStorage := objects.NewObjectStorage(minio, objects.ObjectStorageConfig{Endpoint: cfg.Minio.Endpoint, BucketName: cfg.Minio.BucketName, UseSSL: cfg.Minio.UseSSL}, mediaTypeMap)

	//Init services
	mediaService := media.NewMediaService(dataStorage, objectStorage)
	authService := auth.NewAuthService(
		auth.AuthServiceConfig{PrivateKey: privateKey, PublicKey: publicKey, AccessTTL: cfg.JWT.AccessTTL, RefreshTTL: cfg.JWT.RefreshTTL},
		dataStorage,
		userCache,
		mediaService,
		data.NewTokenStorage(redis))
	tweetService := tweets.NewTweetService(dataStorage, mediaService)
	userService := user.NewUserService(dataStorage, mediaService)

	//Init handler
	handler := http.NewHandler(
		authService,
		tweetService,
		userService,
		authService,
		authService,
		mediaService,
	)
	srv := new(servers.Server)
	if err := srv.Run(cfg.App.Port, handler.InitRouters()); err != nil {
		logrus.Fatalf("error occurred while running http server: %s", err.Error())
	}
}

func loadRSAKeys(privateKeyPath, publicKeyPath string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privatePEM, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read private key: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	publicPEM, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read public key: %w", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicPEM)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	return privateKey, publicKey, nil
}
