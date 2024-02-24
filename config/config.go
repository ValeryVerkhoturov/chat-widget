package config

import (
	"github.com/joho/godotenv"
	"os"
	"strings"
)

const (
	SessionName    = "chat-widget"
	SessionUserKey = "userId"
)

var (
	Port          string
	Host          string
	PublicUrl     string
	TelegramUrl   string
	SessionSecret string
	MongoDBUri    string
	Origins       []string
)

func init() {
	var err error

	if os.Getenv("GO_ENV") != "production" {
		if err = godotenv.Load(); err != nil {
			panic(err)
		}
	}

	Port = os.Getenv("PORT")
	Host = os.Getenv("HOST")
	PublicUrl = os.Getenv("PUBLIC_URL")
	TelegramUrl = os.Getenv("TELEGRAM_URL")
	SessionSecret = os.Getenv("SESSION_SECRET")
	MongoDBUri = os.Getenv("MONGODB_URI")
	Origins = strings.Split(os.Getenv("ORIGINS"), ",")

	if len(Origins) == 0 || len(Port) == 0 || len(Host) == 0 || len(PublicUrl) == 0 ||
		len(TelegramUrl) == 0 || len(SessionSecret) == 0 || len(MongoDBUri) == 0 {
		panic("Invalid env variables")
	}
}
