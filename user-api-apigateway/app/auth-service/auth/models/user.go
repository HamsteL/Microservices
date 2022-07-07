package models

import (
	"auth-service/auth/middlewares"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

type User struct {
	ID           uint      `gorm:"primary_key" json:"id"`
	FirstName    string    `gorm:"size:255;" json:"first_name"`
	LastName     string    `gorm:"size:255;" json:"last_name"`
	Email        string    `gorm:"size:255; not null;" json:"email"`
	PasswordHash string    `gorm:"size:255; not null;" json:"password_hash"`
	UpdateDate   time.Time `json:"update_date"`
}

type Password struct {
	RawPassword  string
	PasswordHash string
}

func FindUserByEmail(host, port, email string) (*User, error) {
	const endpoint string = "getUserByEmail"
	var urlE string = fmt.Sprintf("http://%s:%s/%s?", host, port, endpoint)

	payload := url.Values{}
	payload.Add("email", email)
	req, err := http.NewRequest("GET", urlE+payload.Encode(), nil)
	if err != nil {
		_ = fmt.Errorf("Error while create request", err)
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		_ = fmt.Errorf("Error while get user", err)
		return nil, err
	}

	var user User
	errParse := json.NewDecoder(resp.Body).Decode(&user)
	if errParse != nil {
		errParse = fmt.Errorf("No Body(FindUserByEmail)", errParse)
		return nil, errParse
	}

	return &user, nil
}

func RegisterUser(host, port, email, passwordHash string) (*User, error) {
	const endpoint string = "users"
	var url string = fmt.Sprintf("http://%s:%s/%s", host, port, endpoint)

	postBody, _ := json.Marshal(map[string]string{
		"email":         email,
		"password_hash": passwordHash,
	})
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(url, "application/json", responseBody)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	var user User
	errParse := json.NewDecoder(resp.Body).Decode(&user)
	if errParse != nil {
		errParse = fmt.Errorf("No Body(RegisterUser)", errParse)
		return nil, errParse
	}

	return &user, nil
}

func IsPasswordsMatch(password, userPasswordHash string) bool {
	passwordHash := middlewares.GetStringHash(password)
	if passwordHash != userPasswordHash {
		return false
	}

	return true
}
