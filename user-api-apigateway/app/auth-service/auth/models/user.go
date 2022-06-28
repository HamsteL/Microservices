package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	const endpoint string = "userByEmail"
	var url string = fmt.Sprintf("http://%s:%s/%s/%s", host, port, endpoint, email)

	resp, err := http.Get(url)
	if err != nil {
		_ = fmt.Errorf("No User", err)
		return nil, err
	}
	defer resp.Body.Close()

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
		"email":    email,
		"password": passwordHash,
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
