package controllers

import (
	"fmt"
	"log"
	"net/http"
	"user-api/app/models"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Server struct {
	DB         *gorm.DB
	Router     *mux.Router
	signSecret string
}

func (server *Server) Initialize(settings models.AppSettings) {
	ConnectToDB(server, &settings)
	server.Router = mux.NewRouter()
	server.signSecret = settings.SignSecret
	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, server.Router))
}

func BuildConnectionString(settings *models.AppSettings) string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		settings.DbHost,
		settings.DbPort,
		settings.DbUser,
		settings.DbName,
		settings.DbPassword)
}

func ConnectToDB(server *Server, settings *models.AppSettings) {
	connectionString := BuildConnectionString(settings)
	var err error
	server.DB, err = gorm.Open(settings.DbDriver, connectionString)
	if err != nil {
		log.Fatal("Error while connecting to DB:", err)
	}
}
