package user

import (
	"context"
	"time"

	"github.com/kust1q/Zapp/backend/internal/dto"
)

func (s *userService) Update(ctx context.Context, userID int, req *dto.UpdateBioRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return s.storage.UpdateUserBio(ctx, userID, req.Bio)
}
