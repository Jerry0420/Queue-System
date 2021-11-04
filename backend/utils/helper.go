package utils

func CheckKeysInMap(data map[string]interface{}, keys []string) bool {
	for _, key := range keys {
		if _, ok := data[key]; !ok {
			return false
		}
	}
	return true
}