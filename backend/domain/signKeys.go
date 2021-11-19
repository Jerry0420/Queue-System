package domain

type signKeyTypes struct{ NORMAL, PASSWORD, REFRESH, SESSION string }

var SignKeyTypes signKeyTypes = signKeyTypes{NORMAL: "normal", PASSWORD: "password", REFRESH: "refresh", SESSION: "session"}

type SignKey struct {
	ID          int    `json:"id"`
	StoreId     int    `json:"store_id"`
	SignKey     string `json:"sign_key"`
	SignKeyType string `json:"type"`
}
