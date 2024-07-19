package txtban

import (
	"encoding/json"
	"io"
	"net/http"
)

// general json data type for handling json request body
type JsonData struct {
	UserId      string `json:"uuid,omitempty"`
	Password    string `json:"password,omitempty"`
	Name        string `json:"name,omitempty"`
	OldPassword string `json:"old_password,omitempty"`
	NewPassword string `json:"new_password,omitempty"`
}

func parseJsonBody(body []byte, v any) error {
	return json.Unmarshal(body, v)
}

func readRequestBody(body io.Reader) ([]byte, error) {
	return io.ReadAll(body)
}

func sendJson(w http.ResponseWriter, v any) {
	data, err := json.Marshal(v)
	if err != nil {
		sendError(w, errInternalServerError, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
