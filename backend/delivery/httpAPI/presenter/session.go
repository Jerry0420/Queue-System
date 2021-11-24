package presenter

import (
	"bytes"
	"encoding/json"

	"github.com/jerry0420/queue-system/backend/domain"
)

func SessionCreate(session domain.StoreSession) string {
	// TODO: combine url and qrcode
	var flushedData bytes.Buffer
	json.NewEncoder(&flushedData).Encode(session)
	return flushedData.String()
}
