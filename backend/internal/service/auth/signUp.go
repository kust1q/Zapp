package auth

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"image/jpeg"
	"strings"

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
)

func (s *authService) SignUp(ctx context.Context, user dto.CreateUserRequest) (dto.UserResponse, error) {
	user.Email = strings.ToLower(user.Email)

	if exists, err := s.cache.Exists(ctx, cache.EmailType, user.Email); err != nil {
		return dto.UserResponse{}, err
	} else if exists {
		return dto.UserResponse{}, ErrEmailAlreadyUsed
	}

	if exists, err := s.cache.Exists(ctx, cache.UsernameType, user.Username); err != nil {
		return dto.UserResponse{}, err
	} else if exists {
		return dto.UserResponse{}, ErrUsernameAlreadyUsed
	}

	_, err := s.storage.GetUserByEmail(ctx, user.Email)
	if err == nil {
		if err := s.cache.Add(ctx, cache.EmailType, user.Email); err != nil {
			logrus.Printf("failed to cache email: %v", err)
		}
		return dto.UserResponse{}, ErrEmailAlreadyUsed
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return dto.UserResponse{}, fmt.Errorf("email check failed: %w", err)
	}

	_, err = s.storage.GetUserByUsername(ctx, user.Username)
	if err == nil {
		if err := s.cache.Add(ctx, cache.UsernameType, user.Username); err != nil {
			logrus.Printf("failed to cache username: %v", err)
		}
		return dto.UserResponse{}, ErrUsernameAlreadyUsed
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return dto.UserResponse{}, fmt.Errorf("username check failed: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.UserResponse{}, fmt.Errorf("password hashing failed: %w", err)
	}

	avatarURL, err := s.generateAndUploadAvatar(ctx, user.Username, user.Gen)
	if err != nil {
		return dto.UserResponse{}, err
	}

	domainUser := entity.User{
		Username:    user.Username,
		Email:       user.Email,
		Password:    string(hashedPassword),
		Bio:         user.Bio,
		Gen:         user.Gen,
		AvatarURL:   avatarURL,
		IsSuperuser: false,
	}

	createdUser, err := s.storage.CreateUser(ctx, domainUser)
	if err != nil {
		if delErr := s.media.Remove(ctx, avatarURL); delErr != nil {
			logrus.Printf("failed to remove avatar: %v", delErr)
		}
		return dto.UserResponse{}, fmt.Errorf("user creation failed: %w", err)
	}

	return dto.UserResponse{
		ID:          createdUser.ID,
		Username:    createdUser.Username,
		Email:       createdUser.Email,
		Bio:         createdUser.Bio,
		Gen:         createdUser.Gen,
		AvatarURL:   createdUser.AvatarURL,
		CreatedAt:   createdUser.CreatedAt,
		IsSuperuser: createdUser.IsSuperuser,
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
