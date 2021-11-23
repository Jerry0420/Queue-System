package domain

type storeSessionStatus struct{ NORMAL, SCANNED, USED string }

var StoreSessionStatus storeSessionStatus = storeSessionStatus{NORMAL: "normal", SCANNED: "scanned", USED: "used"}

type StoreSession struct {
	ID      string `json:"id"`
	StoreId int    `json:"store_id"`
	StoreSessionStatus  string `json:"status"`
}
