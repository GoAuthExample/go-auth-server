package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Group(func(r chi.Router) {
		r.Use(userMiddleware)
		r.Get("/user", s.fetchUserData)
	})

	r.Get("/auth/{provider}/callback", s.getAuthCallbackFunction)
	r.Post("/auth/logout", s.Logout)
	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}

func (s *Server) getAuthCallbackFunction(w http.ResponseWriter, r *http.Request) {

	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

	if user, err := gothic.CompleteUserAuth(w, r); err == nil {
		log.Printf("User", user)

		_, err := s.db.Login(&user)
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
		}

		session, err := gothic.Store.Get(r, "session")
		if err != nil {
			log.Printf("Error making new session: %v", err)
			return
		}
		session.Values["user"] = user
		err = session.Save(r, w)
		if err != nil {
			log.Printf("Error saving new session: %v", err)
			return
		}

		http.Redirect(w, r, "http://localhost:8080/dashboard", http.StatusFound)
	} else {
		gothic.BeginAuthHandler(w, r)
	}

}

func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {

	session, err := gothic.Store.Get(r, "session")
	if err != nil {
		log.Printf("Error getting session for logout: %v", err)
		http.Error(w, "Logout error", http.StatusInternalServerError)
		return
	}

	// Destroy the session by setting its MaxAge to a negative value
	// and clearing its values. This also tells the browser to delete the cookie.
	session.Options.MaxAge = -1
	session.Values = make(map[interface{}]interface{}) // Clear stored values
	err = session.Save(r, w)                           // Save the modified session to apply changes (delete cookie)

	if err != nil {
		log.Printf("Error saving session during logout: %v", err)
		http.Error(w, "Logout error", http.StatusInternalServerError)
		return
	}

	log.Printf("Logout printing session: %v", session.Values)
	render.Render(w, r, PostResponseRender("Logged out successfully"))
}

func (s *Server) fetchUserData(w http.ResponseWriter, r *http.Request) {

	gothUser := r.Context().Value("user").(goth.User)

	// It's not really necessary to make a DB call here but I'm just simulating an actual DB operation
	user, err := s.db.GetUser(gothUser.UserID)
	if err != nil {
		render.Render(w, r, ErrServerError(err))
		return
	}

	log.Printf("User JSON: %v", user)
	render.Render(w, r, UserResponseRender(user))

	return
}

func userMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s, err := gothic.Store.Get(r, "session")

		if err != nil || s == nil || s.Values["user"] == nil {
			render.Render(w, r, ErrInvalidRequest(fmt.Errorf("No user session")))
			return
		}

		if s.Options.MaxAge < 0 {
			render.Render(w, r, ErrInvalidRequest(fmt.Errorf("No user session")))
			return
		}

		u := s.Values["user"].(goth.User)
		ctx := context.WithValue(r.Context(), "user", u)
		next.ServeHTTP(w, r.WithContext(ctx)) // Hand off the request with the user attached
	})
}
