package controllers

import (
	"simple-crud-app-go/simple-user-api-go/api/middlewares"
)

func (server *Server) initializeRoutes() {
	//Users routes
	server.Router.HandleFunc("/users", middlewares.AddHeaderToRequest(server.CreateUser)).Methods("POST")
	server.Router.HandleFunc("/users", middlewares.AddHeaderToRequest(server.GetUsers)).Methods("GET")
	server.Router.HandleFunc("/users/{id}", middlewares.AddHeaderToRequest(server.GetUser)).Methods("GET")
	server.Router.HandleFunc("/users/{id}", middlewares.AddHeaderToRequest(middlewares.AddHeaderToRequest(server.UpdateUser))).Methods("PUT")
	server.Router.HandleFunc("/users/{id}", middlewares.AddHeaderToRequest(server.DeleteUser)).Methods("DELETE")
	server.Router.HandleFunc("/health", middlewares.AddHeaderToRequest(server.HealthCheck)).Methods("GET")
}
