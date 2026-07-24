package models

import "time"

type RefreshToken struct {
	id         string
	token_hash string
	user_id    string
	created_at time.Time
	revoked_at time.Time
	expired_at time.Time
}
