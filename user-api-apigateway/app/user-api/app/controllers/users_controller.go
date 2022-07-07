package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"user-api/app/models"
	"user-api/app/responses"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

const sessionIdCookieName string = "session_id"

func (server *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	userCreated, err := user.SaveUser(server.DB)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userCreated.ID))
	responses.JSON(w, http.StatusCreated, userCreated)
}

func (server *Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	users, err := user.FindAllUsers(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, users)
}

func (server *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	user := models.User{}
	userGotten, err := user.FindUserByID(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	if userAuthorized, err := server.IsAuthorizedUser(r, user.Email); err != nil || !userAuthorized {
		responses.ERROR(w, http.StatusForbidden, err)
		return
	}

	responses.JSON(w, http.StatusOK, map[string]string{
		"FirstName": userGotten.FirstName,
		"LastName":  userGotten.LastName,
		"Email":     userGotten.Email,
	})
}

func (server *Server) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	// ParseForm parses the raw query from the URL and updates r.Form
	parseErr := r.ParseForm()
	if parseErr != nil {
		http.Error(w, "Please pass the data as URL form encoded", http.StatusBadRequest)
		return
	}
	email := r.Form.Get("email")

	if len(email) == 0 {
		responses.ERROR(w, http.StatusBadRequest, fmt.Errorf("No email in request"))
		return
	}
	user := models.User{}
	userGotten, err := user.FindUserByEmail(server.DB, email)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, userGotten)
}

func (server *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	userFromDb := models.User{}
	userGotten, err := userFromDb.FindUserByID(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	if userAuthorized, err := server.IsAuthorizedUser(r, userGotten.Email); err != nil || !userAuthorized {
		responses.ERROR(w, http.StatusForbidden, err)
		return
	}

	updatedUser, err := user.UpdateUser(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, map[string]string{
		"FirstName": updatedUser.FirstName,
		"LastName":  updatedUser.LastName,
		"Email":     updatedUser.Email,
	})
}

func (server *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := models.User{}
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	_, err = user.DeleteUser(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	if userAuthorized, err := server.IsAuthorizedUser(r, user.Email); err != nil || !userAuthorized {
		responses.ERROR(w, http.StatusForbidden, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, "")
}

func (server *Server) getEmailFromCookie(r *http.Request) (string, error) {
	sessionIdCookie, err := r.Cookie(sessionIdCookieName)
	if err != nil {
		fmt.Printf("Cant find auth cookie\r\n")
		return "", fmt.Errorf("Cant find auth cookie")
	}

	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(sessionIdCookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(server.signSecret), nil
	})
	if err != nil {
		return "Error while decode token", err
	}

	for key, val := range claims {
		fmt.Printf("Key: %v, value: %v\n", key, val)
		if key == "Email" {
			return val.(string), nil
		}
	}

	return "", fmt.Errorf("Email not found in Token")
}

func (server *Server) IsAuthorizedUser(r *http.Request, userEmail string) (bool, error) {
	emailFromCookie, err := server.getEmailFromCookie(r)
	if err != nil {
		return false, err
	}

	if userEmail != emailFromCookie {
		return false, fmt.Errorf("Forbiden content")
	}

	return true, nil
}
