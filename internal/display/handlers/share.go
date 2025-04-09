package handlers

import "net/http"

func writeResponse(w http.ResponseWriter, statusCode int, contentType, body string) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	w.Write([]byte(body))
}
