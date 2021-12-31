package domain

type storeSessionStatus struct{ NORMAL, SCANNED, USED string }

var StoreSessionStatus storeSessionStatus = storeSessionStatus{NORMAL: "normal", SCANNED: "scanned", USED: "used"}
const StoreSessionString string = "session"

type StoreSession struct {
	ID                 string `json:"id"`
	StoreId            int    `json:"-"`
	StoreSessionStatus string `json:"status"`
}
