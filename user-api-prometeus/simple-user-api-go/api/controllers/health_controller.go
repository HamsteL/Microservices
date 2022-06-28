package controllers

import (
	"net/http"
	"simple-crud-app-go/simple-user-api-go/api/responses"
)

func (server *Server) HealthCheck(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "status:OK")
}
