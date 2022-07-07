package controllers

import (
	"auth-service/auth/middlewares"
	"auth-service/auth/models"
	"auth-service/auth/responses"
	"fmt"
	"net/http"
	"time"
)

const sessionIdCookieName string = "session_id"

func (server *Server) SignUp(w http.ResponseWriter, r *http.Request) {
	parseErr := r.ParseForm()
	if parseErr != nil {
		http.Error(w, "Please pass the data as URL form encoded", http.StatusBadRequest)
		return
	}
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	u, findErr := models.FindUserByEmail(server.AppSettings.ApiHost, server.AppSettings.ApiPort, email)
	if findErr != nil || u != nil {
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
	sessionId, err := getAuthCookie(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	if sessionId != "" {
		responses.JSON(w, http.StatusBadRequest, "You already logged in. Log out first")
		return
	}

	// ParseForm parses the raw query from the URL and updates r.Form
	parseErr := r.ParseForm()
	if parseErr != nil {
		http.Error(w, "Please pass the data as URL form encoded", http.StatusBadRequest)
		return
	}
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	fmt.Printf("Login request for %s\n", email)

	user, findErr := models.FindUserByEmail(server.AppSettings.ApiHost, server.AppSettings.ApiPort, email)
	if findErr != nil {
		responses.ERROR(w, http.StatusBadRequest, findErr)
		return
	}

	if !models.IsPasswordsMatch(password, user.PasswordHash) {
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

	err = server.Storage.CreateSession(email)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	setAuthCookie(&w, server.Storage.Sessions[email])
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
	sessionId, err := getAuthCookie(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	server.Storage.DeleteSession(sessionId)
	deleteAuthCookie(&w)

	responses.JSON(w, http.StatusOK, map[string]string{
		"status": fmt.Sprintf("Logout successful", sessionId),
	})
}

func (server *Server) Auth(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Auth start\n")
	parseErr := r.ParseForm()
	if parseErr != nil {
		http.Error(w, "Please pass the data as URL form encoded", http.StatusBadRequest)
		return
	}

	sessionId, err := getAuthCookie(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	reqUrl := r.Form.Get("req_url")

	fmt.Printf("Auth params: sessionId(from cookie) = %s, req_url = %v\n", sessionId, reqUrl)
	sessionExists, email := server.Storage.SessionExists(sessionId)
	if !sessionExists {
		fmt.Printf("Session not exists 401\n")
		responses.JSON(w, http.StatusUnauthorized, nil)
		return
	}
	fmt.Printf("Session exists OK\n")

	w.Header().Set("x-username", email)
	w.Header().Set("x-auth-token", sessionId)
	responses.JSON(w, http.StatusOK, nil)
	return
}

func setAuthCookie(w *http.ResponseWriter, sessionId string) {
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{
		Name:    sessionIdCookieName,
		Value:   sessionId,
		Expires: expiration,
		Path:    "/",
	}
	http.SetCookie(*w, &cookie)
}

func deleteAuthCookie(w *http.ResponseWriter) {
	expiration := time.Unix(0, 0)
	cookie := http.Cookie{
		Name:    sessionIdCookieName,
		Value:   "",
		Expires: expiration,
		Path:    "/",
	}
	http.SetCookie(*w, &cookie)
}

func getAuthCookie(r *http.Request) (string, error) {
	sessionId, err := r.Cookie(sessionIdCookieName)
	if err != nil {
		fmt.Printf("Cant find auth cookie\r\n")
		return "", fmt.Errorf("Cant find auth cookie")
	}

	return sessionId.Value, nil
}
