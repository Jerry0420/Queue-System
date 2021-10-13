package domain

import "time"

type Store struct {
	ID int `json:"id"`
	Email string `json:"email"`
	Password string `json:"password"`
	Name string `json:"name"`
	Description string `json:"description"`
	CreatedAt time.Time `json:"created_at"`
	Status string `json:"status"`
	SessionID string `json:"session_id"`
}

type StoreRepositoryInterface interface {

}

type StoreUsecaseInterface interface {

}

// query := `INSERT INTO stores (email, password, name, description, status) VALUES ($1, $2, $3, $4, $5) `
// stmt, err := db.PrepareContext(r.Context(), query)
// _, err = stmt.ExecContext(r.Context(), "jeerywa@gmail.com", "im password", "jerry", "im description", "open")
// utils.JsonResponseOK(w, nil)

// query := `UPDATE stores SET name=$1,session_id=uuid_generate_v4()`
// stmt, err := db.PrepareContext(r.Context(), query)
// _, err = stmt.ExecContext(r.Context(), "new_jerry")

// query = `SELECT id,name,created_at,status,session_id FROM stores`
// rows, err := db.QueryContext(r.Context(), query)
// defer func() {
// 	errRow := rows.Close()
// }()
// result := make([]domain.Store, 0)
// for rows.Next() {
// 	store := domain.Store{}
// 	err = rows.Scan(&store.ID, &store.Name, &store.CreatedAt, &store.Status, &store.SessionID)
// 	result = append(result, store)
// }