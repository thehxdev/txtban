package txtban

import "errors"

func getAuthKey(headers map[string][]string) (string, error) {
	authHeader := headers["Authorization"]
	if len(authHeader) == 0 {
		return "", errors.New("Authorization header is empty")
	}
	return authHeader[0], nil
}
