package health

import "net/http"

func HealthHandler(svc Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(svc.Check()))
	})
}
