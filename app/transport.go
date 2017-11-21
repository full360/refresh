package app

import (
	"encoding/json"
	"net/http"
)

func RefreshHandler(svc Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		_, err := svc.Refresh()
		if err != nil {
			encodeError(err, w)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	})
}

// encode errors from business-logic
func encodeError(err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
