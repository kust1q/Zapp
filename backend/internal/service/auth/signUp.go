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
	"github.com/kust1q/Zapp/backend/internal/providers/db/redis/cache"
	"github.com/o1egl/govatar"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func (s *authService) SignUp(ctx context.Context, req *entity.User) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	req.Credential.Email = strings.ToLower(strings.TrimSpace(req.Credential.Email))
	req.Username = strings.TrimSpace(req.Username)

	if err := s.checkUserExists(ctx, req.Credential.Email, req.Username); err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Credential.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("password hashing failed: %w", err)
	}

	user := entity.User{
		Username:    req.Username,
		Bio:         req.Bio,
		Gen:         req.Gen,
		CreatedAt:   req.CreatedAt,
		IsSuperuser: req.IsSuperuser,
		Credential: &entity.Credential{
			Email:    req.Credential.Email,
			Password: string(hashedPassword),
		},
	}

	createdUser, err := s.db.CreateUserTx(ctx, tx, &user)
	if err != nil {
		return nil, fmt.Errorf("user creation failed: %w", err)
	}

	avatar, err := s.generateAndUploadAvatar(ctx, createdUser.ID, createdUser.Username, createdUser.Gen, tx)

	if err != nil {
		return nil, fmt.Errorf("failed to generate or upload avatar: %w", err)
	}

	avatarURL, err := s.media.GetPresignedURL(ctx, avatar.Path)
	if err != nil {
		return nil, fmt.Errorf("get avatar presigned url failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction failed: %w", err)
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.cache.Add(ctx, cache.EmailType, req.Credential.Email); err != nil {
			logrus.Warnf("failed to cache email for user %s: %v", req.Username, err)
		}

		if err := s.cache.Add(ctx, cache.UsernameType, req.Username); err != nil {
			logrus.Warnf("failed to cache username for user %s: %v", req.Username, err)
		}
	}()

	return &entity.User{
		ID:        createdUser.ID,
		Username:  createdUser.Username,
		Bio:       createdUser.Bio,
		Gen:       createdUser.Gen,
		AvatarURL: avatarURL,
		CreatedAt: createdUser.CreatedAt,
		Credential: &entity.Credential{
			Email: createdUser.Credential.Email,
		},
	}, nil
}

func (s *authService) generateAndUploadAvatar(ctx context.Context, userID int, username, gender string, tx *sql.Tx) (*entity.Avatar, error) {
	genMap := map[string]govatar.Gender{
		"male":   govatar.MALE,
		"female": govatar.FEMALE,
	}
	gen, ok := genMap[strings.ToLower(gender)]
	if !ok {
		return nil, fmt.Errorf("invalid gender")
	}

	img, err := govatar.GenerateForUsername(gen, username)
	if err != nil {
		return nil, fmt.Errorf("avatar generation failed: %w", err)
	}
	/*
		var avatarBuffer bytes.Buffer
		if err := png.Encode(&avatarBuffer, img); err != nil {
			return entity.User{}, err
		}
	*/
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80}); err != nil {
		return nil, fmt.Errorf("JPEG encoding failed: %w", err)
	}
	avatarSaveName := uuid.New().String() + ".jpg"

	avatar, err := s.media.UploadAvatarTx(ctx, userID, &buf, avatarSaveName, tx)
	if err != nil {
		return nil, fmt.Errorf("avatar upload failed: %w", err)
	}
	return &entity.Avatar{
		ID:        avatar.ID,
		UserID:    avatar.UserID,
		Path:      avatar.Path,
		MimeType:  avatar.MimeType,
		SizeBytes: avatar.SizeBytes,
	}, nil
}

func (s *authService) checkUserExists(ctx context.Context, email, username string) error {
	_, err := s.db.GetUserByEmail(ctx, email)
	if err == nil {
		return ErrEmailAlreadyUsed
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("email check failed: %w", err)
	}

	_, err = s.db.GetUserByUsername(ctx, username)
	if err == nil {
		return ErrUsernameAlreadyUsed
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("username check failed: %w", err)
	}
	return nil
}
