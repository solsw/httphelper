package rest

import (
	"net/http"
)

// IsNotStatusOK determines whether r.StatusCode is not [http.StatusOK].
func IsNotStatusOK(r *http.Response) bool {
	return r.StatusCode != http.StatusOK
}
