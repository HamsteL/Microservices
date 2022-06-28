package controllers

import (
	"auth-service/auth/middlewares"
)

func (server *Server) initializeRoutes() {
	server.Router.Use(middlewares.SetupMiddleware)

	//Auth routes
	server.Router.HandleFunc("/signup", server.SignUp).Methods("POST")
	server.Router.HandleFunc("/login", server.Login).Methods("GET")
	server.Router.HandleFunc("/logout", server.Logout).Methods("GET")

	server.Router.HandleFunc("/health", server.HealthCheck).Methods("GET")
}
