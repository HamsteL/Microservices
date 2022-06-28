package controllers

import (
	"errors"
	"math/rand"
	"net/http"
	"simple-crud-app-go/simple-user-api-go/api/responses"
	"time"
)

func (server *Server) GetRandResponseStatus(w http.ResponseWriter, r *http.Request) {
	// random generator
	rand.Seed(time.Now().UnixNano())
	var isError bool = (rand.Intn(2)) == 0
	if isError {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Error response"))
		return
	}

	responses.JSON(w, http.StatusOK, "status:OK")
}
