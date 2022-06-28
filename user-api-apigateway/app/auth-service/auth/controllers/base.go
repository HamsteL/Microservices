package controllers

import (
	"auth-service/auth/models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"
)

type Server struct {
	DB          *gorm.DB
	Router      *mux.Router
	AppSettings *models.AppSettings
	Storage     *models.CookieStorage
}

func (server *Server) Initialize(settings models.AppSettings) {
	server.Router = mux.NewRouter()
	server.Storage = models.SetupStorage(settings.SignSecret)
	server.AppSettings = &settings
	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
