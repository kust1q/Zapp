package conv

import (
	"time"

	"github.com/kust1q/Zapp/backend/internal/controllers/http/dto/request"
	"github.com/kust1q/Zapp/backend/internal/controllers/http/dto/response"
	"github.com/kust1q/Zapp/backend/internal/domain/entity"
)

// Requests
func FromTweetRequestToDomain(userID int, parent_tweet_id *int, file *entity.File, req *request.Tweet) *entity.Tweet {
	if req == nil {
		return nil
	}

	return &entity.Tweet{
		ParentTweetID: parent_tweet_id,
		Content:       req.Content,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		File:          file,
		Author: &entity.SmallUser{
			ID: userID,
		},
		Counters: nil,
	}
}

func FromTweetUpdateRequestToDomain(userID, tweetID int, file *entity.File, req *request.Tweet) *entity.Tweet {
	if req == nil {
		return nil
	}

	return &entity.Tweet{
		ID:        tweetID,
		Content:   req.Content,
		UpdatedAt: time.Now(),
		File:      file,
		Author: &entity.SmallUser{
			ID: userID,
		},
	}
}

// Responses
func FromDomainToTweetResponse(tweet *entity.Tweet) *response.Tweet {
	if tweet == nil {
		return nil
	}

	responseTweet := &response.Tweet{
		ID:            tweet.ID,
		Content:       tweet.Content,
		CreatedAt:     tweet.CreatedAt,
		UpdatedAt:     tweet.UpdatedAt,
		ParentTweetID: tweet.ParentTweetID,
		MediaURL:      tweet.MediaUrl,
		Author:        FromDomainToSmallUserResponse(tweet.Author),
	}

	// Добавляем счетчики только если они не nil
	if tweet.Counters != nil {
		responseTweet.Counters = &response.Counters{
			ReplyCount:   tweet.Counters.ReplyCount,
			RetweetCount: tweet.Counters.RetweetCount,
			LikeCount:    tweet.Counters.LikeCount,
		}
	}

	return responseTweet
}

/*
func FromDomainToLikesListResponse(users []entity.SmallUser) []response.UserLike {
    if users == nil {
        return nil
    }

    responses := make([]response.UserLike, 0, len(users))
    var userLike response.UserLike
    for _, u := range users {
        userLike = response.UserLike{
            Username:  u.Username,
            AvatarURL: u.AvatarURL,
        }
        responses = append(responses, userLike)
    }
    return responses
}
*/

func FromDomainToTweetListResponse(tweets []entity.Tweet) []response.Tweet {
	if tweets == nil {
		return nil
	}

	res := make([]response.Tweet, 0, len(tweets))
	for _, t := range tweets {
		tweetResponse := FromDomainToTweetResponse(&t)
		if tweetResponse != nil {
			res = append(res, *tweetResponse)
		}
	}
	return res
}
