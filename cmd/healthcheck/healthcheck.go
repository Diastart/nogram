package healthcheck

import (
	"net/http"
)


func HandleHealth(response http.ResponseWriter, request *http.Request) {
    response.Header().Set("Content-Type", "application/json")
    response.WriteHeader(http.StatusOK)
    response.Write([]byte(`{"status": "ok"}`))
}