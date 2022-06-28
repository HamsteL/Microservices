package controllers

import (
	"net/http"
	"user-api/app/responses"
)

func (server *Server) HealthCheck(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "status:OK")
}
