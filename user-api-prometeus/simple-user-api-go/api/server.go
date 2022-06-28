package api

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"simple-crud-app-go/simple-user-api-go/api/controllers"
	"simple-crud-app-go/simple-user-api-go/api/models"
	"strings"
)

var server = controllers.Server{}

func Run() {
	appSettings := GetEnvVars()
	server.Initialize(appSettings)
	fmt.Printf("Listening to port %s\n", appSettings.ApiPort)
	server.Run(appSettings.ApiPort)
}

func GetEnvVars() models.AppSettings {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env variables, %v", err)
	}

	var appSettings models.AppSettings
	appSettings.DbDriver = TrimStr(os.Getenv("DB_DRIVER"))
	appSettings.DbUser = TrimStr(os.Getenv("DB_USER"))
	appSettings.DbPassword = TrimStr(os.Getenv("DB_PASSWORD"))
	appSettings.DbPort = TrimStr(os.Getenv("DB_PORT"))
	appSettings.DbHost = TrimStr(os.Getenv("DB_HOST"))
	appSettings.DbName = TrimStr(os.Getenv("DB_NAME"))
	appSettings.ApiPort = fmt.Sprintf(":%s", TrimStr(os.Getenv("API_PORT")))

	return appSettings
}

func TrimStr(str string) string {
	return strings.TrimSuffix(str, "\n")
}
