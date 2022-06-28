package controllers

import (
	"errors"
	"math/rand"
	"net/http"
	"time"
	"user-api/app/responses"
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
