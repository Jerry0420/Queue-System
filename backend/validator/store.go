package validator

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jerry0420/queue-system/backend/domain"
)

func StoreOpen(r *http.Request) (domain.Store, error) {
	var jsonBody map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&jsonBody)
	if err != nil {
		return domain.Store{}, domain.ServerError40001
	}
	name, ok := jsonBody["name"].(string)
	if !ok || name == "" {
		return domain.Store{}, domain.ServerError40001
	}
	email, ok := jsonBody["email"].(string)
	if !ok || email == "" {
		return domain.Store{}, domain.ServerError40001
	}
	password, ok := jsonBody["password"].(string)
	if !ok || password == "" {
		return domain.Store{}, domain.ServerError40001
	}
	store := domain.Store{Name: name, Email: email, Password: password}
	return store, nil
}

func StoreSignin(r *http.Request) (domain.Store, error) {
	var jsonBody map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&jsonBody)
	if err != nil {
		return domain.Store{}, domain.ServerError40001
	}
	email, ok := jsonBody["email"].(string)
	if !ok || email == "" {
		return domain.Store{}, domain.ServerError40001
	}
	password, ok := jsonBody["password"].(string)
	if !ok || password == "" {
		return domain.Store{}, domain.ServerError40001
	}
	store := domain.Store{Email: email, Password: password}
	return store, nil
}

func StoreTokenRefresh(r *http.Request) (*http.Cookie, error) {
	encryptedRefreshToken, err := r.Cookie(domain.SignKeyTypes.REFRESH)
	if err != nil || len(encryptedRefreshToken.Value) == 0 {
		return nil, domain.ServerError40102
	}
	return encryptedRefreshToken, nil
}

func StoreClose(r *http.Request) (domain.TokenClaims, error) {
	tokenClaims := r.Context().Value(domain.SignKeyTypes.NORMAL).(domain.TokenClaims)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id != tokenClaims.StoreID {
		return domain.TokenClaims{}, domain.ServerError40004
	}
	return tokenClaims, nil
}

func StorePasswordForgot(r *http.Request) (domain.Store, error) {
	var jsonBody map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&jsonBody)
	if err != nil {
		return domain.Store{}, domain.ServerError40001
	}
	email, ok := jsonBody["email"].(string)
	if !ok || email == "" {
		return domain.Store{}, domain.ServerError40001
	}
	store := domain.Store{Email: email}
	return store, nil
}

func StorePasswordUpdate(r *http.Request) (map[string]string, int, error) {
	var jsonBody map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&jsonBody)
	if err != nil {
		return map[string]string{}, -1, domain.ServerError40001
	}
	passwordToken, ok := jsonBody["passwordToken"].(string)
	if !ok || passwordToken == "" {
		return map[string]string{}, -1, domain.ServerError40001
	}
	password, ok := jsonBody["password"].(string)
	if !ok || password == "" {
		return map[string]string{}, -1, domain.ServerError40001
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return map[string]string{}, -1, domain.ServerError40001
	}
	
	body := map[string]string{"passwordToken": passwordToken, "password": password} 
	return body, id, nil
}