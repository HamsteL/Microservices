package controllers

import (
	"auth-service/auth/middlewares"
	"auth-service/auth/models"
	"auth-service/auth/responses"
	"fmt"
	"net/http"
)

func (server *Server) SignUp(w http.ResponseWriter, r *http.Request) {
	parseErr := r.ParseForm()
	if parseErr != nil {
		http.Error(w, "Please pass the data as URL form encoded", http.StatusBadRequest)
		return
	}
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	_, findErr := models.FindUserByEmail(server.AppSettings.ApiHost, server.AppSettings.ApiPort, email)
	if findErr == nil {
		responses.ERROR(w, http.StatusBadRequest, fmt.Errorf("user already exists", findErr))
		return
	}

	passwordHash := middlewares.GetStringHash(password)
	user, createErr := models.RegisterUser(server.AppSettings.ApiHost, server.AppSettings.ApiPort, email, passwordHash)
	if createErr != nil {
		responses.ERROR(w, http.StatusInternalServerError, createErr)
		return
	}

	errSes := server.Storage.CreateSession(user.Email)
	if errSes != nil {
		responses.ERROR(w, http.StatusInternalServerError, errSes)
		return
	}

	responses.JSON(w, http.StatusOK, map[string]string{
		"status":    "User registered",
		"sessionId": server.Storage.Sessions[email],
		"email":     email,
	})
}

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	// ParseForm parses the raw query from the URL and updates r.Form
	parseErr := r.ParseForm()
	if parseErr != nil {
		http.Error(w, "Please pass the data as URL form encoded", http.StatusBadRequest)
		return
	}
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user, findErr := models.FindUserByEmail(server.AppSettings.ApiHost, server.AppSettings.ApiPort, email)
	if findErr != nil {
		responses.ERROR(w, http.StatusBadRequest, findErr)
		return
	}

	passwordHash := middlewares.GetStringHash(password)
	if passwordHash != user.PasswordHash {
		responses.ERROR(w, http.StatusForbidden, fmt.Errorf("Incorrect password"))
		return
	}

	if server.Storage.SessionExistsByEmail(email) {
		responses.JSON(w, http.StatusOK, map[string]string{
			"status":    "User already authorized",
			"sessionId": server.Storage.Sessions[email],
			"email":     email,
		})
		return
	}

	err := server.Storage.CreateSession(email)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, map[string]string{
		"sessionId": server.Storage.Sessions[email],
		"email":     email,
	})
}

func (server *Server) Logout(w http.ResponseWriter, r *http.Request) {
	parseErr := r.ParseForm()
	if parseErr != nil {
		http.Error(w, "Please pass the data as URL form encoded", http.StatusBadRequest)
		return
	}
	sessionId := r.Form.Get("session_id")

	server.Storage.DeleteSession(sessionId)

	responses.JSON(w, http.StatusOK, map[string]string{
		"status": fmt.Sprintf("Session %s was deleted", sessionId),
	})
}

func (server *Server) Auth(w http.ResponseWriter, r *http.Request) {
	parseErr := r.ParseForm()
	if parseErr != nil {
		http.Error(w, "Please pass the data as URL form encoded", http.StatusBadRequest)
		return
	}
	sessionId := r.Form.Get("sessionId")
	reqUrl := r.Form.Get("req_url")

	if !server.Storage.SessionExists(sessionId) {
		responses.JSON(w, http.StatusUnauthorized, nil)
		return
	}

	if len(reqUrl) < 1 {
		responses.JSON(w, http.StatusBadRequest, "Empty URL")
		return
	}

	http.Redirect(w, r, reqUrl, http.StatusSeeOther)
}
