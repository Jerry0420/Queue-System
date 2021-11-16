package presenter

import (
	"encoding/json"
	"time"

	"github.com/jerry0420/queue-system/backend/domain"
)

func StoreForResponse(store domain.Store) map[string]interface{} {
	var storeJson []byte
	var storeMap map[string]interface{}
	storeJson, _ = json.Marshal(store)
	json.Unmarshal(storeJson, &storeMap)
	delete(storeMap, "password")
	return storeMap
}

func StoreToken(store domain.Store, token string, tokenExpiresAt time.Time) map[string]interface{} {
	var storeJson []byte
	var storeMap map[string]interface{}
	storeJson, _ = json.Marshal(store)
	json.Unmarshal(storeJson, &storeMap)
	delete(storeMap, "password")
	storeMap["token"] = token
	storeMap["token_expires_at"] = tokenExpiresAt.Unix()
	return storeMap
}
