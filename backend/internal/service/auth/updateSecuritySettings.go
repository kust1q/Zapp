package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
	"github.com/kust1q/Zapp/backend/internal/dto"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidSecretAnswer = errors.New("Invalid secret answer")
)

func (s *authService) UpdateSecuritySettings(ctx context.Context, userID int, req *dto.UpdateSecuritySettingsRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	oldSecurityData, err := s.storage.GetSecurityDataByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(oldSecurityData.Answer), []byte(req.OldSecretAnswer)); err != nil {
		return ErrInvalidSecretAnswer
	}

	newHashedAnswer, err := bcrypt.GenerateFromPassword([]byte(req.NewSecretAnswer), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash new answer: %w", err)
	}

	return s.storage.SetSecretQuestion(ctx, &entity.SecretQuestion{
		UserID:         userID,
		SecretQuestion: req.NewSecretQuestion,
		Answer:         string(newHashedAnswer),
	})
}
