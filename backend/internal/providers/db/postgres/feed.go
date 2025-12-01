package postgres

import (
	"context"
	"fmt"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
	conv "github.com/kust1q/Zapp/backend/internal/providers/db/conv"
	"github.com/kust1q/Zapp/backend/internal/providers/db/models"
)

func (s *PostgresDB) GetFeedByIds(ctx context.Context, userIDs []int) ([]entity.Tweet, error) {
	query := fmt.Sprintf(`
        SELECT id, user_id, parent_tweet_id, content, created_at, updated_at
        FROM %s
        WHERE user_id = ANY($1)
		UNION ALL
		SELECT t.id, t.user_id, t.parent_tweet_id, t.content, r.created_at, t.updated_at
		FROM %s r
		JOIN %s t ON r.tweet_id = t.id
		WHERE r.user_id = ANY($1)
		ORDER BY created_at DESC
		LIMIT 30`,
		TweetsTable, RetweetsTable, TweetsTable)

	var tweetModels []models.Tweet
	if err := s.db.SelectContext(ctx, &tweetModels, query, userIDs); err != nil {
		return nil, fmt.Errorf("failed to get tweets by username: %w", err)
	}
	return conv.FromTweetModelToDomainList(tweetModels), nil
}
