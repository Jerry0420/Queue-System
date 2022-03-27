package domain

import "time"

type customerStatus struct{ WAITING, PROCESSING, DONE, DELETE string }

var CustomerStatus customerStatus = customerStatus{WAITING: "waiting", PROCESSING: "processing", DONE: "done", DELETE: "delete"}

type Customer struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	QueueID   int       `json:"-"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
