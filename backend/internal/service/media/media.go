package media

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/kust1q/Zapp/backend/internal/dto"
	"github.com/kust1q/Zapp/backend/internal/storage/objects"
	"github.com/sirupsen/logrus"
)

type mediaStorage interface {
	// Tweet
	UpsertByTweetIdTx(ctx context.Context, tx *sql.Tx, media *dto.TweetMedia) (*dto.TweetMedia, error)
	GetMediaByTweetID(ctx context.Context, tweetID int) (*dto.TweetMedia, error)
	DeleteMediaByTweetID(ctx context.Context, tweetID int) error
	// User
	UploadAvatarTx(ctx context.Context, tx *sql.Tx, avatar *dto.Avatar) (*dto.Avatar, error)
	GetAvatarByUserID(ctx context.Context, userID int) (*dto.Avatar, error)
	DeleteAvatarByUserID(ctx context.Context, userID int) error
}

type objectStorage interface {
	Upload(ctx context.Context, file io.Reader, mediaType objects.MediaType, filename string) (path string, mimeType string, err error)
	Remove(ctx context.Context, objectPath string) error
}

type mediaService struct {
	storage mediaStorage
	object  objectStorage
}

func NewMediaService(storage mediaStorage, object objectStorage) *mediaService {
	return &mediaService{storage: storage, object: object}
}

func (s *mediaService) UploadAndAttachTweetMediaTx(ctx context.Context, tweetID int, file io.Reader, filename string, tx *sql.Tx) (*dto.TweetMedia, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	mt := s.detectMediaType(filename)
	path, mime, err := s.object.Upload(ctx, file, mt, filename)
	if err != nil {
		return &dto.TweetMedia{}, err
	}
	tweetMedia, err := s.storage.UpsertByTweetIdTx(ctx, tx, &dto.TweetMedia{
		TweetID:  tweetID,
		MediaURL: path,
		MimeType: mime,
	})
	if err != nil {
		s.CleanUpMedia(ctx, path)
		return &dto.TweetMedia{}, fmt.Errorf("upsert media failed: %w", err)
	}
	return tweetMedia, nil
}

func (s *mediaService) GetMediaByTweetID(ctx context.Context, tweetID int) (*dto.TweetMedia, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.storage.GetMediaByTweetID(ctx, tweetID)
}

func (s *mediaService) DeleteTweetMedia(ctx context.Context, tweetID int) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	media, err := s.storage.GetMediaByTweetID(ctx, tweetID)
	if err != nil {
		return err
	}
	if err := s.storage.DeleteMediaByTweetID(ctx, tweetID); err != nil {
		return err
	}
	if media.MediaURL != "" {
		s.CleanUpMedia(ctx, media.MediaURL)
	}
	return nil
}

func (s *mediaService) UploadAvatarTx(ctx context.Context, userID int, file io.Reader, filename string, tx *sql.Tx) (*dto.Avatar, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	mt := s.detectMediaType(filename)
	path, mime, err := s.object.Upload(ctx, file, mt, filename)
	if err != nil {
		return &dto.Avatar{}, err
	}
	avatar, err := s.storage.UploadAvatarTx(ctx, tx, &dto.Avatar{
		UserID:   userID,
		MediaURL: path,
		MimeType: mime,
	})
	if err != nil {
		s.CleanUpMedia(ctx, path)
		return &dto.Avatar{}, fmt.Errorf("attach media failed: %w", err)
	}
	return avatar, nil
}

func (s *mediaService) GetAvatarByUserID(ctx context.Context, userID int) (*dto.Avatar, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return s.storage.GetAvatarByUserID(ctx, userID)
}

func (s *mediaService) DeleteAvatar(ctx context.Context, userID int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	avatar, err := s.storage.GetAvatarByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get data of useravatar: %w", err)
	}
	if err := s.storage.DeleteAvatarByUserID(ctx, userID); err != nil {
		return err
	}
	if avatar.MediaURL != "" {
		s.CleanUpMedia(ctx, avatar.MediaURL)
	}
	return nil
}

func (s *mediaService) CleanUpMedia(ctx context.Context, mediaURL string) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := s.object.Remove(ctx, mediaURL); err != nil {
		logrus.WithField("media_url", mediaURL).Warnf("failed to remove media: %s", err)
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
