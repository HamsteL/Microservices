package controllers

import (
	"auth-service/auth/responses"
	"net/http"
)

func (server *Server) HealthCheck(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, map[string]string{
		"status": "OK",
	})
}
