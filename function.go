package function

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/CloudAutomationSolutions/practicecalendar-backend/auth"
)

// F ...
// Wrapper to add middlelayer check: https://auth0.com/docs/quickstart/backend/golang/01-authorization#validate-access-tokens
func F(w http.ResponseWriter, r *http.Request) {

	iss := os.Getenv("ISS")
	if iss == "" {
		log.Printf("Cannot find ISS configuration")
		http.Error(w, "Cannot find ISS configuration", http.StatusInternalServerError)
	}

	aud := os.Getenv("AUD")
	if aud == "" {
		log.Printf("Cannot find AUD configuration")
		http.Error(w, "Cannot find AUD configuration", http.StatusInternalServerError)
	}

	jwtMiddleware, err := auth.GetJWTMiddleware(aud, iss)
	if err != nil {
		log.Printf("Cannot create JWT authentication middleware: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	token := r.Header.Get("JTW-Token")
	if token == "" {
		log.Println("Token Empty. Unauthorised!")
		http.Error(w, "JWT Token not provided!", http.StatusUnauthorized)
	}

	jwtMiddleware.HandlerWithNext(w, r, actualHTTPHandlerFunction)

}

func actualHTTPHandlerFunction(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Cannot read body: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write([]byte(body))
}
