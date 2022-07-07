package controllers

import (
	"auth-service/auth/middlewares"
	"net/http"
)

func (server *Server) initializeRoutes() {
	server.Router.Use(middlewares.SetupMiddleware)

	//Auth routes
	server.Router.HandleFunc("/signup", server.SignUp).Methods(http.MethodPost)
	server.Router.HandleFunc("/login", server.Login).Methods(http.MethodGet)
	server.Router.HandleFunc("/logout", server.Logout).Methods(http.MethodGet)

	server.Router.HandleFunc("/auth{req_url:/?.*}", server.Auth).Methods(http.MethodGet)

	server.Router.HandleFunc("/health", server.HealthCheck).Methods(http.MethodGet)
}
