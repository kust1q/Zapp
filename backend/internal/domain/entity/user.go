package entity

import (
	"time"
)

type User struct {
	ID          int       `db:"id"`
	Username    string    `db:"username"`
	Email       string    `db:"email"`
	Password    string    `db:"password"`
	Bio         string    `db:"bio"`
	Gen         string    `db:"gen"`
	AvatarURL   string    `db:"avatar_url"`
	CreatedAt   time.Time `db:"created_at"`
	IsSuperuser bool      `db:"is_superuser"`
}

type Follow struct {
	ID          int       `db:"id"`
	FollowerID  int       `db:"follower_id"`
	FollowingID int       `db:"following_id"`
	CreatedAt   time.Time `db:"created_at"`
}
