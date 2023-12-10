package routes

import (
	"context"
	"encoding/json"
	firebase "firebase.google.com/go/v4"
	"github.com/go-chi/chi"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"log"
	"main-api/db"
	"net/http"
	"os"
	"time"
)

func ApiRouter() *chi.Mux {
	r := chi.NewRouter()

	db.ConnectDatabase()

	gcpServiceAccount := map[string]string{
		"type":                        "service_account",
		"project_id":                  os.Getenv("GCP_PROJECT_ID"),
		"private_key_id":              os.Getenv("GCP_PRIVATE_KEY_ID"),
		"private_key":                 os.Getenv("GCP_PRIVATE_KEY"),
		"client_email":                os.Getenv("GCP_CLIENT_EMAIL"),
		"client_id":                   os.Getenv("GCP_CLIENT_ID"),
		"client_x509_cert_url":        os.Getenv("GCP_CLIENT_X509_CERT_URL"),
		"auth_uri":                    "https://accounts.google.com/o/oauth2/auth",
		"token_uri":                   "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"universe_domain":             "googleapis.com",
	}

	gcpServiceAccountJson, err := json.Marshal(gcpServiceAccount)

	if err != nil {
		log.Fatalf("Error marshalling fb config json: %s", err)
	}

	credentials, _ := google.CredentialsFromJSON(context.Background(), gcpServiceAccountJson, []string{"https://www.googleapis.com/auth/cloud-platform"}...)

	config := &firebase.Config{ProjectID: os.Getenv("FB_PROJECT_ID")}

	app, err := firebase.NewApp(context.Background(), config, option.WithCredentials(credentials))
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	ctx := context.Background()

	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	type LoginRequest struct {
		Token string `json:"token"`
	}

	r.Post("/v1/auth/login", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var requestBody LoginRequest

		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		expiresIn := time.Hour * 24 * 14

		token, err := client.SessionCookie(r.Context(), requestBody.Token, expiresIn)
		if err != nil {
			log.Printf("error verifying ID token: %v\n\n", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Write([]byte(token))
	})

	r.Get("/v1/auth/user", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		session := r.Header.Get("session")

		if err != nil {
			log.Printf("error getting cookie: %v\n\n", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		token, err := client.VerifySessionCookieAndCheckRevoked(r.Context(), session)
		if err != nil {
			log.Printf("error verifying session cookie: %v\n\n", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := client.GetUser(r.Context(), token.UID)

		// Create an instance of UserResponse and populate it
		response := map[string]string{
			"user_id": user.UID,
			"email":   user.Email,
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(response)
	})

	return r
}
