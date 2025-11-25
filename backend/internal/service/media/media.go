package media

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
	"github.com/kust1q/Zapp/backend/internal/storage/objects"
	"github.com/sirupsen/logrus"
)

type mediaStorage interface {
	// Tweet
	UpsertByTweetIdTx(ctx context.Context, tx *sql.Tx, media *entity.TweetMedia) (*entity.TweetMedia, error)
	GetMediaPathByTweetID(ctx context.Context, tweetID int) (string, error)
	GetMediaDataByTweetID(ctx context.Context, tweetID int) (*entity.TweetMedia, error)
	DeleteMediaByTweetID(ctx context.Context, tweetID, userID int) error
	// User
	UploadAvatarTx(ctx context.Context, tx *sql.Tx, avatar *entity.Avatar) (*entity.Avatar, error)
	GetAvatarPathByUserID(ctx context.Context, userID int) (string, error)
	GetAvatarDataByUserID(ctx context.Context, userID int) (*entity.Avatar, error)
	DeleteAvatarByUserID(ctx context.Context, userID int) error
}

type objectStorage interface {
	Upload(ctx context.Context, file io.Reader, mediaType objects.MediaType, filename string) (path string, mimeType string, err error)
	Remove(ctx context.Context, objectPath string) error
	GetPresignedURL(ctx context.Context, objectPath string) (string, error)
}

type mediaService struct {
	db     mediaStorage
	object objectStorage
}

func NewMediaService(db mediaStorage, object objectStorage) *mediaService {
	return &mediaService{db: db, object: object}
}

func (s *mediaService) UploadAndAttachTweetMediaTx(ctx context.Context, tweetID int, file io.Reader, filename string, tx *sql.Tx) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	mt, err := s.detectMediaType(filename)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	size, err := io.Copy(&buf, file)
	if err != nil {
		return "", fmt.Errorf("failed to read file for size calculation: %w", err)
	}

	bufReader := bytes.NewReader(buf.Bytes())

	path, mime, err := s.object.Upload(ctx, bufReader, mt, filename)
	if err != nil {
		return "", err
	}
	tweetMedia, err := s.db.UpsertByTweetIdTx(ctx, tx, &entity.TweetMedia{
		TweetID:   tweetID,
		Path:      path,
		MimeType:  mime,
		SizeBytes: size,
	})
	if err != nil {
		s.CleanUpMedia(ctx, path)
		return "", fmt.Errorf("upsert media failed: %w", err)
	}
	return s.object.GetPresignedURL(ctx, tweetMedia.Path)
}

func (s *mediaService) GetMediaUrlByTweetID(ctx context.Context, tweetID int) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	mediaPath, err := s.db.GetMediaPathByTweetID(ctx, tweetID)
	if err != nil {
		return "", fmt.Errorf("failed to get media url: %w", err)
	}

	return s.object.GetPresignedURL(ctx, mediaPath)
}

func (s *mediaService) GetMediaDataByTweetID(ctx context.Context, tweetID int) (*entity.TweetMedia, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	tweetMedia, err := s.db.GetMediaDataByTweetID(ctx, tweetID)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get tweet media data: %w", err)
	} else if err == sql.ErrNoRows {
		return &entity.TweetMedia{}, nil
	}

	mediaURL, err := s.object.GetPresignedURL(ctx, tweetMedia.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to get tweet media url: %w", err)
	}

	return &entity.TweetMedia{
		ID:        tweetMedia.ID,
		TweetID:   tweetMedia.TweetID,
		Path:      mediaURL,
		MimeType:  tweetMedia.MimeType,
		SizeBytes: tweetMedia.SizeBytes,
	}, nil
}

func (s *mediaService) DeleteTweetMedia(ctx context.Context, tweetID, userID int) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	media, err := s.db.GetMediaDataByTweetID(ctx, tweetID)
	if err != nil {
		return err
	}
	if err := s.db.DeleteMediaByTweetID(ctx, tweetID, userID); err != nil {
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

	mt, err := s.detectMediaType(filename)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	size, err := io.Copy(&buf, file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file for size calculation: %w", err)
	}

	bufReader := bytes.NewReader(buf.Bytes())

	path, mime, err := s.object.Upload(ctx, bufReader, mt, filename)
	if err != nil {
		return nil, err
	}

	avatar, err := s.db.UploadAvatarTx(ctx, tx, &entity.Avatar{
		UserID:    userID,
		Path:      path,
		MimeType:  mime,
		SizeBytes: size,
	})

	if err != nil {
		s.CleanUpMedia(ctx, path)
		return nil, fmt.Errorf("upload avatar failed: %w", err)
	}

	avatar.Path, err = s.object.GetPresignedURL(ctx, avatar.Path)
	if err != nil {
		return nil, fmt.Errorf("get avatar url failed: %w", err)
	}

	return avatar, nil
}

func (s *mediaService) GetAvatarUrlByUserID(ctx context.Context, userID int) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	avatarPath, err := s.db.GetMediaPathByTweetID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to get avatar url: %w", err)
	}

	return s.object.GetPresignedURL(ctx, avatarPath)
}

func (s *mediaService) GetAvatarDataByUserID(ctx context.Context, userID int) (*entity.Avatar, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	avatar, err := s.db.GetAvatarDataByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get avatar data: %w", err)
	}

	avatarPath, err := s.db.GetAvatarPathByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get avatar path: %w", err)
	}

	avatarUrl, err := s.object.GetPresignedURL(ctx, avatarPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get avatar url: %w", err)
	}

	avatar.Path = avatarUrl

	return avatar, nil
}

func (s *mediaService) DeleteAvatar(ctx context.Context, userID int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	avatar, err := s.db.GetAvatarDataByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get data of useravatar: %w", err)
	}

	if err := s.db.DeleteAvatarByUserID(ctx, userID); err != nil {
		return err
	}

	if avatar.Path != "" {
		s.CleanUpMedia(ctx, avatar.Path)
	}
	return nil
}

func (s *mediaService) CleanUpMedia(ctx context.Context, Path string) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := s.object.Remove(ctx, Path); err != nil {
		logrus.WithField("media_url", Path).Warnf("failed to remove media: %s", err)
	}
}

func (s *mediaService) detectMediaType(filename string) (objects.MediaType, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp":
		return objects.TypeImage, nil
	case ".gif":
		return objects.TypeGIF, nil
	case ".mp4", ".mov", ".m4v":
		return objects.TypeVideo, nil
	case ".mp3", ".wav", ".ogg", ".flac", ".aac", ".m4a", ".webm":
		return objects.TypeAudio, nil
	default:
		return "", fmt.Errorf("invalid media type")
	}
}
