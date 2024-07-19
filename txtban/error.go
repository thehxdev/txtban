package txtban

import (
	"encoding/json"
	"net/http"

	"github.com/thehxdev/txtban/tberr"
)

var (
	errEmptyAuthorizationHeader = tberr.New("Bad request", "set 'Authorization' header")
	errUnauthorized             = tberr.New("Unauthorized")
	errInternalServerError      = tberr.New("Internal Server Error")
	errEmptyTxtID               = tberr.New("txt id is empty", "txt id could not be empty")
	errEmptyTxtName             = tberr.New("txt name is empty", "txt name could not be empty")
	errBadJsonData              = tberr.New("failed to parse request json data")
)

func sendError(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)

	errData, err := json.Marshal(err)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errData = []byte(`{"error": "Internal Server Error"}`)
	}

	w.Write(errData)
}
