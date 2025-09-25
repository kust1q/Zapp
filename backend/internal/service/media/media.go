package media

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
	"github.com/kust1q/Zapp/backend/internal/dto"
	"github.com/kust1q/Zapp/backend/internal/storage/objects"
	"github.com/sirupsen/logrus"
)

type mediaStorage interface {
	// Tweet
	UpsertByTweetIdTx(ctx context.Context, tx *sql.Tx, media *entity.TweetMedia) (*entity.TweetMedia, error)
	GetMediaUrlByTweetID(ctx context.Context, tweetID int) (string, error)
	GetMediaDataByTweetID(ctx context.Context, tweetID int) (*entity.TweetMedia, error)
	DeleteMediaByTweetID(ctx context.Context, tweetID, userID int) error
	// User
	UploadAvatarTx(ctx context.Context, tx *sql.Tx, avatar *entity.Avatar) (*entity.Avatar, error)
	GetAvatarUrlByUserID(ctx context.Context, userID int) (string, error)
	GetAvatarDataByUserID(ctx context.Context, userID int) (*entity.Avatar, error)
	DeleteAvatarByUserID(ctx context.Context, userID int) error
}

type objectStorage interface {
	Upload(ctx context.Context, file io.Reader, mediaType objects.MediaType, filename string) (path string, mimeType string, err error)
	Remove(ctx context.Context, objectPath string) error
	GetPresignedURL(ctx context.Context, objectPath string) (string, error)
}

type mediaService struct {
	storage mediaStorage
	object  objectStorage
}

func NewMediaService(storage mediaStorage, object objectStorage) *mediaService {
	return &mediaService{storage: storage, object: object}
}

func (s *mediaService) UploadAndAttachTweetMediaTx(ctx context.Context, tweetID int, file io.Reader, filename string, tx *sql.Tx) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	mt := s.detectMediaType(filename)
	path, mime, err := s.object.Upload(ctx, file, mt, filename)
	if err != nil {
		return "", err
	}
	tweetMedia, err := s.storage.UpsertByTweetIdTx(ctx, tx, &entity.TweetMedia{
		TweetID:  tweetID,
		Path:     path,
		MimeType: mime,
	})
	if err != nil {
		s.CleanUpMedia(ctx, path)
		return "", fmt.Errorf("upsert media failed: %w", err)
	}
	return s.GetPresignedURL(ctx, tweetMedia.Path)
}

func (s *mediaService) GetMediaUrlByTweetID(ctx context.Context, tweetID int) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.storage.GetMediaUrlByTweetID(ctx, tweetID)
}

func (s *mediaService) GetMediaDataByTweetID(ctx context.Context, tweetID int) (*dto.TweetMediaDataResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	tweetMedia, err := s.storage.GetMediaDataByTweetID(ctx, tweetID)
	if err != nil {
		return &dto.TweetMediaDataResponse{}, fmt.Errorf("failed to get tweet media data: %w", err)
	}
	mediaURL, err := s.storage.GetMediaUrlByTweetID(ctx, tweetID)
	if err != nil {
		return &dto.TweetMediaDataResponse{}, fmt.Errorf("failed to get tweet media url: %w", err)
	}
	return &dto.TweetMediaDataResponse{
		ID:        tweetMedia.ID,
		TweetID:   tweetMedia.TweetID,
		MediaURL:  mediaURL,
		MimeType:  tweetMedia.MimeType,
		SizeBytes: tweetMedia.SizeBytes,
	}, nil
}

func (s *mediaService) DeleteTweetMedia(ctx context.Context, tweetID, userID int) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	media, err := s.storage.GetMediaDataByTweetID(ctx, tweetID)
	if err != nil {
		return err
	}
	if err := s.storage.DeleteMediaByTweetID(ctx, tweetID, userID); err != nil {
		return err
	}
	if media.Path != "" {
		s.CleanUpMedia(ctx, media.Path)
	}
	return nil
}

func (s *mediaService) UploadAvatarTx(ctx context.Context, userID int, file io.Reader, filename string, tx *sql.Tx) (*entity.Avatar, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	mt := s.detectMediaType(filename)
	path, mime, err := s.object.Upload(ctx, file, mt, filename)
	if err != nil {
		return &entity.Avatar{}, err
	}
	avatar, err := s.storage.UploadAvatarTx(ctx, tx, &entity.Avatar{
		UserID:   userID,
		Path:     path,
		MimeType: mime,
	})
	if err != nil {
		s.CleanUpMedia(ctx, path)
		return &entity.Avatar{}, fmt.Errorf("attach media failed: %w", err)
	}
	return avatar, nil
}

func (s *mediaService) GetAvatarUrlByUserID(ctx context.Context, userID int) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return s.storage.GetAvatarUrlByUserID(ctx, userID)
}

func (s *mediaService) GetAvatarDataByUserID(ctx context.Context, userID int) (*dto.AvatarDataResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	avatar, err := s.storage.GetAvatarDataByUserID(ctx, userID)
	if err != nil {
		return &dto.AvatarDataResponse{}, fmt.Errorf("failed to get tweet media data: %w", err)
	}
	avatarURL, err := s.storage.GetAvatarUrlByUserID(ctx, userID)
	if err != nil {
		return &dto.AvatarDataResponse{}, fmt.Errorf("failed to get tweet media url: %w", err)
	}
	return &dto.AvatarDataResponse{
		ID:        avatar.ID,
		UserID:    avatar.UserID,
		AvatarURL: avatarURL,
		MimeType:  avatar.MimeType,
		SizeBytes: avatar.SizeBytes,
	}, nil
}

func (s *mediaService) DeleteAvatar(ctx context.Context, userID int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	avatar, err := s.storage.GetAvatarDataByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get data of useravatar: %w", err)
	}
	if err := s.storage.DeleteAvatarByUserID(ctx, userID); err != nil {
		return err
	}
	if avatar.Path != "" {
		s.CleanUpMedia(ctx, avatar.Path)
	}
	return nil
}

func (s *mediaService) GetPresignedURL(ctx context.Context, objectPath string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.object.GetPresignedURL(ctx, objectPath)
}

func (s *mediaService) CleanUpMedia(ctx context.Context, Path string) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := s.object.Remove(ctx, Path); err != nil {
		logrus.WithField("media_url", Path).Warnf("failed to remove media: %s", err)
	}
}

func (s *mediaService) detectMediaType(filename string) objects.MediaType {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp":
		return objects.TypeImage
	case ".gif":
		return objects.TypeGIF
	case ".mp4", ".mov", ".m4v":
		return objects.TypeVideo
	case ".mp3", ".wav", ".ogg", ".flac", ".aac", ".m4a", ".webm":
		return objects.TypeAudio
	default:
		return objects.TypeImage
	}
}
