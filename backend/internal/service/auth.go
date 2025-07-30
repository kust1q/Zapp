package service

import (
	"bytes"
	"context"
	"fmt"
	"image/png"

	"github.com/google/uuid"
	"github.com/kust1q/Zapp/backend/internal/domain/entity"
	"github.com/kust1q/Zapp/backend/internal/dto"
	"github.com/o1egl/govatar"
	"golang.org/x/crypto/bcrypt"
)

type AuthStorage interface {
	UploadAvatar(ctx context.Context, avatarBuffer bytes.Buffer, avatarSaveName string) (string, error)
	CreateUser(ctx context.Context, user entity.User) (entity.User, error)
}

type authService struct {
	storage AuthStorage
}

func NewAuthService(storage AuthStorage) *authService {
	return &authService{
		storage: storage,
	}
}

func (s *authService) CreateUser(ctx context.Context, user dto.CreateUserRequest) (entity.User, error) {
	hashed_password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return entity.User{}, err
	}

	var gen govatar.Gender
	if user.Gen == "Male" {
		gen = govatar.MALE
	} else {
		gen = govatar.FEMALE
	}

	img, err := govatar.GenerateForUsername(gen, user.Username)
	if err != nil {
		return entity.User{}, err
	}

	var avatarBuffer bytes.Buffer
	if err := png.Encode(&avatarBuffer, img); err != nil {
		return entity.User{}, err
	}

	avatarSaveName := fmt.Sprintf("%s/%s-%s%s", "avatars", user.Username, uuid.New().String(), ".png")

	avatarURL, err := s.storage.UploadAvatar(ctx, avatarBuffer, avatarSaveName)
	if err != nil {
		return entity.User{}, err
	}

	domainUser := entity.User{
		Username:    user.Username,
		Email:       user.Email,
		Password:    string(hashed_password),
		Bio:         user.Bio,
		Gen:         user.Gen,
		AvatarURL:   avatarURL,
		IsSuperuser: false,
	}

	return s.storage.CreateUser(ctx, domainUser)
}
