package media

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
)

var (
	ErrFileTooLarge     = errors.New("file too large")
	ErrInvalidMediaType = errors.New("invalid media type")
)

type MediaType string

const (
	TypeAvatar MediaType = "avatar"
	TypeImage  MediaType = "image"
	TypeVideo  MediaType = "video"
	TypeAudio  MediaType = "audio"
	TypeGIF    MediaType = "gif"
)

type MediaTypeConfig struct {
	MaxSize       int64
	AllowedMime   []string
	AllowedExt    []string
	ForceMimeType string
}

type MediaStorageConfig struct {
	Endpoint   string
	BucketName string
	UseSSL     bool
}

type mediaStorage struct {
	mc        *minio.Client
	cfg       MediaStorageConfig
	mediaCfgs map[MediaType]MediaTypeConfig
}

func NewMediaStorage(mc *minio.Client, cfg MediaStorageConfig, mediaCfgs map[MediaType]MediaTypeConfig) *mediaStorage {
	return &mediaStorage{
		mc:        mc,
		cfg:       cfg,
		mediaCfgs: mediaCfgs,
	}
}

func (s *mediaStorage) Upload(ctx context.Context, file io.Reader, mediaType MediaType, filename string) (path string, mimeType string, err error) {
	cfg, ok := s.mediaCfgs[mediaType]
	if !ok {
		return "", "", fmt.Errorf("unsupported media type: %s", mediaType)
	}
	ext := strings.ToLower(filepath.Ext(filename))
	data, err := s.readAndValidate(file, ext, cfg)
	if err != nil {
		return "", "", err
	}

	if cfg.ForceMimeType != "" {
		mimeType = cfg.ForceMimeType
	} else {
		mimeType = http.DetectContentType(data)
	}

	path = filepath.Join(string(mediaType), time.Now().Format("2006/01/02"), filename)

	_, err = s.mc.PutObject(
		ctx,
		s.cfg.BucketName,
		path,
		bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{
			ContentType: mimeType,
		})

	if err != nil {
		return "", "", fmt.Errorf("minio upload failed: %w", err)
	}

	return path, mimeType, nil
}

func (s *mediaStorage) GetURL(objectPath string) string {
	protocol := "http"
	if s.cfg.UseSSL {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s",
		protocol,
		s.cfg.Endpoint,
		s.cfg.BucketName,
		objectPath,
	)
}

func (s *mediaStorage) Remove(ctx context.Context, objectPath string) error {
	return s.mc.RemoveObject(
		ctx,
		s.cfg.BucketName,
		objectPath,
		minio.RemoveObjectOptions{},
	)
}

func (s *mediaStorage) readAndValidate(reader io.Reader, ext string, cfg MediaTypeConfig) ([]byte, error) {
	limitedReader := io.LimitReader(reader, cfg.MaxSize+1)
	data, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, err
	}

	if int64(len(data)) > cfg.MaxSize {
		return nil, ErrFileTooLarge
	}

	detectedType := http.DetectContentType(data)
	if !s.isAllowedMimeType(detectedType, cfg) {
		return nil, fmt.Errorf("invalid media type: %s", detectedType)
	}
	if !s.isValidExtension(ext, cfg) {
		return nil, fmt.Errorf("invalid extension: %s", ext)
	}

	return data, nil
}

func (s *mediaStorage) isValidExtension(ext string, cfg MediaTypeConfig) bool {
	if len(cfg.AllowedExt) == 0 {
		return true
	}

	for _, allowed := range cfg.AllowedExt {
		if strings.EqualFold(ext, allowed) {
			return true
		}
	}
	return false
}

func (s *mediaStorage) isAllowedMimeType(mimeType string, cfg MediaTypeConfig) bool {
	for _, allowed := range cfg.AllowedMime {
		if mimeType == allowed {
			return true
		}
	}
	return false
}
