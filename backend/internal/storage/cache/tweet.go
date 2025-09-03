package cache

import (
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	LikeType    = "like"
	RetweetType = "retweet"
	ReplyType   = "reply"

	prefixLike    = "l:"
	prefixReply   = "r:"
	prefixRetweet = "t:"
)

var (
	tweetMap = map[string]string{
		LikeType:    prefixLike,
		ReplyType:   prefixReply,
		RetweetType: prefixRetweet,
	}
)

type tweetCache struct {
	redis *redis.Client
	ttl   time.Duration
}

func NewTweetCache(redis *redis.Client, ttl time.Duration) *tweetCache {
	return &tweetCache{
		redis: redis,
		ttl:   ttl,
	}
}
