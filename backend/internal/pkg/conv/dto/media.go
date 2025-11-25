package conv

import (
	"github.com/kust1q/Zapp/backend/internal/controllers/http/dto/response"
	"github.com/kust1q/Zapp/backend/internal/domain/entity"
)

func FromDomainToAvatarResponse(avatar *entity.Avatar) *response.Avatar {
	if avatar == nil {
		return nil
	}

	return &response.Avatar{
		AvatarURL: avatar.Path,
		MimeType:  avatar.MimeType,
		SizeBytes: avatar.SizeBytes,
	}
}

func FromDomainToMediaResponse(media *entity.TweetMedia) *response.TweetMedia {
	if media == nil {
		return nil
	}

	return &response.TweetMedia{
		MediaURL:  media.Path,
		MimeType:  media.MimeType,
		SizeBytes: media.SizeBytes,
	}
}
