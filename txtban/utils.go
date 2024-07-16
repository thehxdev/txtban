package txtban

import (
	"github.com/thehxdev/txtban/tberr"
)

func getAuthKey(headers map[string][]string) (string, error) {
	authHeader := headers["Authorization"]
	if len(authHeader) == 0 {
		return "", tberr.New("authorization header is empty")
	}
	return authHeader[0], nil
}
