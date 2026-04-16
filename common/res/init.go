package res

import (
	"net/http"

	"github.com/atharvyadav96k/gcp/common"
)

func Send(w http.ResponseWriter, status int, message string, data interface{}) {
	response := Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
	// Convert response to JSON and send it back to the client
	json, _ := common.ToJSON(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}
