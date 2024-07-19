package txtban

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/thehxdev/txtban/tberr"
)

// general json data type for handling json request body
type JsonData struct {
	UserId      string `json:"uuid,omitempty"`
	Password    string `json:"password,omitempty"`
	Name        string `json:"name,omitempty"`
	OldPassword string `json:"old_password,omitempty"`
	NewPassword string `json:"new_password,omitempty"`
}

func getHeaderValue(headers map[string][]string, name string, idx int) (string, error) {
	header := headers[name]
	headerLen := len(header)

	if headerLen == 0 {
		return "", tberr.New(fmt.Sprintf("%s header is empty", name))
	}

	if idx < 0 || idx >= headerLen {
		return "", tberr.New("invalid index value")
	}

	return header[idx], nil
}

func getAuthKey(headers map[string][]string) (string, error) {
	return getHeaderValue(headers, "Authorization", 0)
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
	w.Write(data)
}
