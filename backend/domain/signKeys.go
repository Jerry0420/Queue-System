package domain

import (
	"time"
)

type signKeyTypes struct{ NORMAL, PASSWORD, REFRESH string }

var SignKeyTypes signKeyTypes = signKeyTypes{NORMAL: "normal", PASSWORD: "password", REFRESH: "refresh"}

type SignKey struct {
	ID          int       `json:"id"`
	StoreId     int       `json:"store_id"`
	SignKey     string    `json:"sign_key"`
	SignKeyType string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
}
