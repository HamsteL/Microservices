package models

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Token struct {
	Email string
	*jwt.StandardClaims
}

type CookieStorage struct {
	Sessions   map[string]string
	signSecret string
}

func SetupStorage(signSecret string) *CookieStorage {
	storage := new(CookieStorage)
	storage.signSecret = signSecret
	storage.Sessions = make(map[string]string)

	return storage
}

func (storage *CookieStorage) CreateSession(email string) error {
	expiresAt := time.Now().Add(time.Minute * 100000).Unix()
	tk := &Token{
		Email: email,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, error := token.SignedString([]byte(storage.signSecret))
	if error != nil {
		return fmt.Errorf("Error while create session", error)
	}
	storage.Sessions[email] = tokenString

	return nil
}

func (storage *CookieStorage) DeleteSession(sessionId string) error {
	if len(storage.Sessions) < 1 {
		return fmt.Errorf("session is not exists")
	}

	var sessionWasDeleted bool = false
	for email, session := range storage.Sessions {
		if session == sessionId {
			delete(storage.Sessions, email)
			sessionWasDeleted = true
			break
		}
	}

	if !sessionWasDeleted {
		return fmt.Errorf("session is not exists")
	}
	return nil
}
