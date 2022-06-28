package auth

import (
	"auth-service/auth/controllers"
	"auth-service/auth/models"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

var server = controllers.Server{}

func Run() {
	appSettings := GetEnvVars()
	server.Initialize(appSettings)
	fmt.Printf("Listening to port %s\n", appSettings.AuthPort)
	server.Run(appSettings.AuthPort)
}

func GetEnvVars() models.AppSettings {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env variables, %v", err)
	}

	var appSettings models.AppSettings
	appSettings.ApiHost = TrimStr(os.Getenv("API_HOST"))
	appSettings.ApiPort = TrimStr(os.Getenv("API_PORT"))
	appSettings.AuthPort = fmt.Sprintf(":%s", TrimStr(os.Getenv("AUTH_PORT")))
	appSettings.SignSecret = TrimStr(os.Getenv("SIGN_SECRET"))

	return appSettings
}

func TrimStr(str string) string {
	return strings.TrimSuffix(str, "\n")
}
