package domain

type Queue struct {
	ID int `json:"id"`
	Name string `json:"name"`
	StoreID int `json:"store_id"`
}
