package auth

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"image/jpeg"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kust1q/Zapp/backend/internal/domain/entity"
	"github.com/kust1q/Zapp/backend/internal/dto"
	"github.com/kust1q/Zapp/backend/internal/storage/cache"
	media "github.com/kust1q/Zapp/backend/internal/storage/objects"
	"github.com/o1egl/govatar"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUsernameAlreadyUsed = errors.New("username already used")
	ErrEmailAlreadyUsed    = errors.New("email already used")
	ErrInvalidGender       = errors.New("invalid gender")
	ErrCacheUnavailable    = errors.New("cache service unavailable")
	ErrInvalidInput        = errors.New("invalid input data")
)

func (s *authService) SignUp(ctx context.Context, user *dto.SignUpRequest) (dto.SignUpResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	user.Username = strings.TrimSpace(user.Username)

	if err := s.checkUserExists(ctx, user.Email, user.Username); err != nil {
		return dto.SignUpResponse{}, err
	}

	avatarURL, err := s.generateAndUploadAvatar(ctx, user.Username, user.Gen)
	if err != nil {
		return dto.SignUpResponse{}, err
	}

	tx, err := s.storage.BeginTx(ctx)
	if err != nil {
		s.cleanupAvatar(ctx, avatarURL)
		return dto.SignUpResponse{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		s.cleanupAvatar(ctx, avatarURL)
		return dto.SignUpResponse{}, fmt.Errorf("password hashing failed: %w", err)
	}

	domainUser := entity.User{
		Username:    user.Username,
		Email:       user.Email,
		Password:    string(hashedPassword),
		Bio:         user.Bio,
		Gen:         user.Gen,
		AvatarURL:   avatarURL,
		CreatedAt:   time.Now(),
		IsSuperuser: false,
	}

	createdUser, err := s.storage.CreateUserTx(ctx, tx, &domainUser)
	if err != nil {
		if delErr := s.media.Remove(ctx, avatarURL); delErr != nil {
			logrus.Warnf("failed to remove avatar: %v", delErr)
		}
		return dto.SignUpResponse{}, fmt.Errorf("user creation failed: %w", err)
	}

	hashedAnswer, err := bcrypt.GenerateFromPassword([]byte(user.SecretAnswer), bcrypt.DefaultCost)
	if err != nil {
		return dto.SignUpResponse{}, fmt.Errorf("secret answer hashing failed: %w", err)
	}

	if err := s.storage.SetSecretQuestionTx(ctx, tx, &entity.SecretQuestion{UserID: createdUser.ID, SecretQuestion: user.SecretQuestion, Answer: string(hashedAnswer)}); err != nil {
		return dto.SignUpResponse{}, fmt.Errorf("set secret question failed: %w", err)
	}
	if err := tx.Commit(); err != nil {
		s.cleanupAvatar(ctx, avatarURL)
		return dto.SignUpResponse{}, fmt.Errorf("commit transaction failed: %w", err)
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.cache.Add(ctx, cache.EmailType, user.Email); err != nil {
			logrus.Warnf("failed to cache email for user %s: %v", user.Username, err)
		}

		if err := s.cache.Add(ctx, cache.UsernameType, user.Username); err != nil {
			logrus.Warnf("failed to cache username for user %s: %v", user.Username, err)
		}
	}()

	return dto.SignUpResponse{
		ID:        createdUser.ID,
		Username:  createdUser.Username,
		Email:     createdUser.Email,
		Bio:       createdUser.Bio,
		Gen:       createdUser.Gen,
		AvatarURL: createdUser.AvatarURL,
		CreatedAt: createdUser.CreatedAt,
	}, nil
}

func (s *authService) generateAndUploadAvatar(ctx context.Context, username, gender string) (string, error) {
	genMap := map[string]govatar.Gender{
		"male":   govatar.MALE,
		"female": govatar.FEMALE,
	}
	gen, ok := genMap[strings.ToLower(gender)]
	if !ok {
		return "", ErrInvalidGender
	}

	img, err := govatar.GenerateForUsername(gen, username)
	if err != nil {
		return "", fmt.Errorf("avatar generation failed: %w", err)
	}
	/*
		var avatarBuffer bytes.Buffer
		if err := png.Encode(&avatarBuffer, img); err != nil {
			return entity.User{}, err
		}
	*/
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80}); err != nil {
		return "", fmt.Errorf("JPEG encoding failed: %w", err)
	}

	avatarSaveName := uuid.New().String() + ".jpg"

	avatarURL, _, err := s.media.Upload(ctx, &buf, media.TypeAvatar, avatarSaveName)
	if err != nil {
		return "", fmt.Errorf("avatar upload failed: %w", err)
	}
	return avatarURL, nil
}

func (s *authService) checkUserExists(ctx context.Context, email, username string) error {
	_, err := s.storage.GetUserByEmail(ctx, email)
	if err == nil {
		return ErrEmailAlreadyUsed
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("email check failed: %w", err)
	}

	_, err = s.storage.GetUserByUsername(ctx, username)
	if err == nil {
		return ErrUsernameAlreadyUsed
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("username check failed: %w", err)
	}
	return nil
}

func (s *authService) cleanupAvatar(ctx context.Context, avatarURL string) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := s.media.Remove(ctx, avatarURL); err != nil {
		logrus.Warnf("failed to remove avatar %s: %v", avatarURL, err)
	}
}
