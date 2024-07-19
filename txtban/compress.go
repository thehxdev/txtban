package txtban

import (
	"compress/gzip"
	"net/http"
)

func sendCompressed(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Encoding", "gzip")

	writer := gzip.NewWriter(w)
    defer writer.Close()

	_, err := writer.Write(data)
	if err != nil {
		sendError(w, errInternalServerError, http.StatusInternalServerError)
	}

	err = writer.Flush()
	if err != nil {
		sendError(w, errInternalServerError, http.StatusInternalServerError)
	}
}

func acceptsGzip(header []string) bool {
	for _, v := range header {
		if v == "gzip" {
			return true
		}
	}
	return false
}
