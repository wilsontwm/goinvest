package utils

import (
	"net/http"
)

// Get the param key in the URL
func GetParam(r *http.Request, key string) (result string) {
	keys, ok := r.URL.Query()[key]
	if ok && len(keys[0]) > 0 {
		result = keys[0]
	}

	return
}
