package domain

type Queue struct {
	ID int `json:"id"`
	Name string `json:"name"`
	StoreID int `json:"store_id"`
}

type QueueRepository interface {

}

type QueueUsecase interface {
	
}