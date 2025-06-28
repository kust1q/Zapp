package domain

import "time"

type Users struct {
	Id         int       `json:"-"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Bio        string    `json:"bio"`
	Avatar_url string    `json:"avatar_url"`
	Created_at time.Time `json:"created_at"`
}

type Follows struct {
	Follower_id  int       `json:"follower_id"`
	Following_id int       `json:"following_id"`
	Created_at   time.Time `json:"created_at"`
}
