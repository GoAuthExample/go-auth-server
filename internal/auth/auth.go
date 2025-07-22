package auth

import (
	"log"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

var (
	key    = os.Getenv("SESSION_SECRET")
	maxAge = 8640 // 1 day
	isProd = false
)

func NewAuth() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env file")
	}

	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	store := sessions.NewCookieStore([]byte(key))

	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = false // HttpOnly should always be enabled
	store.Options.Secure = isProd

	gothic.Store = store
	goth.UseProviders(google.New(googleClientId, googleClientSecret, "http://localhost:3000/auth/google/callback"))

}
