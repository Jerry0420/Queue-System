package domain

import (
	"context"
	"time"
)

type signKeyTypes struct{ SIGNIN, EMAIL string }

var SignKeyTypes signKeyTypes = signKeyTypes{SIGNIN: "signin", EMAIL: "email"}

type SignKey struct {
	ID int `json:"id"`
	StoreId     int       `json:"store_id"`
	SignKey     string    `json:"sign_key"`
	SignKeyType string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
}

type SignKeyRepositoryInterface interface {
	Create(ctx context.Context, signKey SignKey) (signKeyID int, err error)
}
