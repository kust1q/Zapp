package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kust1q/Zapp/backend/internal/dto"
	"golang.org/x/crypto/bcrypt"
)

type AccessClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (s *authService) SignIn(ctx context.Context, credential *dto.SignInRequest) (*dto.SignInResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	credential.Email = strings.ToLower(strings.TrimSpace(credential.Email))
	credential.Password = strings.TrimSpace(credential.Password)

	user, err := s.storage.GetUserByEmail(ctx, credential.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &dto.SignInResponse{}, ErrInvalidCredentials
		}
		return &dto.SignInResponse{}, fmt.Errorf("failed to find user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credential.Password)); err != nil {
		return &dto.SignInResponse{}, ErrInvalidCredentials
	}

	role := "user"
	if user.IsSuperuser {
		role = "admin"
	}

	accessToken, err := s.generateAccessToken(user.ID, user.Email, role)
	if err != nil {
		return &dto.SignInResponse{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.generateRefreshToken(ctx, strconv.Itoa(user.ID))
	if err != nil {
		return &dto.SignInResponse{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &dto.SignInResponse{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}

func (s *authService) generateAccessToken(userID int, email, role string) (string, error) {
	claims := AccessClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.cfg.AccessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "zapp",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.cfg.PrivateKey)
}

func (s *authService) generateRefreshToken(ctx context.Context, userID string) (string, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	refreshToken := base64.RawURLEncoding.EncodeToString(tokenBytes)
	if err := s.tokens.Store(ctx, refreshToken, userID, s.cfg.RefreshTTL); err != nil {
		return "", fmt.Errorf("failed to store refresh token: %w", err)
	}

	return refreshToken, nil
}

func (s *authService) VerifyAccessToken(tokenString string) (int, error) {
	claims := &AccessClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return s.cfg.PublicKey, nil
		},
		jwt.WithLeeway(10*time.Second),
	)

	if err != nil {
		return 0, fmt.Errorf("token validation failed: %w", err)
	}

	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	return claims.UserID, nil
}
