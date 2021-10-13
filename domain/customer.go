package domain

import "time"

type Customer struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Phone string `json:"phone"`
	QueueID int `json:"queue_id"`
	Status string `json:"status"`
	CreatedAt time.Time `json:"created_at"`

}

type CustomerRepositoryInterface interface {

}

type CustomerUsecaseInterface interface {
	
}