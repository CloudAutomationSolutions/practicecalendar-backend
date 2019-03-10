package function

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
)

func init() {
	iss = os.Getenv("ISS")
	aud = os.Getenv("AUD")
	projectID = os.Getenv("GCP_PROJECT")

	var err error

	dsClient, err = NewDB(context.Background(), projectID)
	if err != nil {
		log.Fatalf("Firestore.NewClient: %v", err)
	}
}

var (
	iss       string
	aud       string
	projectID string

	dsClient Datastore
)

// F ...
// Wrapper to add middlelayer check: https://auth0.com/docs/quickstart/backend/golang/01-authorization#validate-access-tokens
func F(w http.ResponseWriter, r *http.Request) {
	// allow cross domain AJAX requests. Make sure to add practicecalendar.com when deploying
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	jwtMiddleware, err := GetJWTMiddleware(aud, iss)
	if err != nil {
		log.Printf("Cannot create JWT authentication middleware: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "OPTIONS":
		return
	case "GET":
		jwtMiddleware.HandlerWithNext(w, r, getHTTPHandler)
	case "POST":
		jwtMiddleware.HandlerWithNext(w, r, postHTTPHandler)
	}
}

func getHTTPHandler(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user")
	if user == nil {
		http.Error(w, "Cannot find subject in the request. Was the JWT Token sent in correctly?", http.StatusNotFound)
		return
	}
	sub := user.(*jwt.Token).Claims.(jwt.MapClaims)["sub"]

	userDS, err := dsClient.GetUser(r.Context(), sub.(string))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching entry for subject %s: %s", sub, err), http.StatusNotFound)
		return
	}

	jsonEntry, err := json.Marshal(&userDS)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching entry for subject %s: %s", sub, err), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonEntry)
}

func postHTTPHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Cannot read body: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := r.Context().Value("user")
	if user == nil {
		http.Error(w, "Cannot find subject in the request. Was the JWT Token sent in correctly?", http.StatusUnauthorized)
		return
	}
	sub := user.(*jwt.Token).Claims.(jwt.MapClaims)["sub"]
	fmt.Fprintf(w, "Subject: %s", sub)

	var projects []Project
	err = json.Unmarshal(body, &projects)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot set user: %s", err), http.StatusBadRequest)
		return
	}
	userEntry := User{sub.(string), projects}

	dsClient.SetUser(r.Context(), &userEntry)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User Entry Updated"))
}
