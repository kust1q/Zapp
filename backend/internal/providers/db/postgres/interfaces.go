package postgres

import "context"

type (
	UserCache interface {
		Exists(ctx context.Context, dataType, data string) (bool, error)
	}
)
