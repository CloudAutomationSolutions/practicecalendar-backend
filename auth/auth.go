package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"path"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
)

type Response struct {
	Message string `json:"message"`
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

// GetJWTMiddleware ...
// Create a *JWTMiddleware object to use as a middleware for verifying incoming requests and the token used with them.
// aud := "YOUR_API_IDENTIFIER" should be the API ID from the Auth0 application we have defined.
// iss := "https://practicecalendar.eu.auth0.com/"
func GetJWTMiddleware(aud, iss string) (*jwtmiddleware.JWTMiddleware, error) {
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			// Verify 'aud' claim
			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
			if !checkAud {
				return token, errors.New("invalid audience")
			}
			// Verify 'iss' claim
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errors.New("invalid issuer")
			}

			u, err := url.Parse(iss)
			if err != nil {
				return token, errors.New("cannot parse iss url passed")
			}
			// Generate the path required to fetch the public key: "https://practicecalendar.eu.auth0.com/.well-known/jwks.json"
			u.Path = path.Join(u.Path, "/.well-known/jwks.json")
			certificateStringURL := u.String()

			cert, err := getPemCert(certificateStringURL, token)
			if err != nil {
				panic(err.Error())
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})

	return jwtMiddleware, nil
}

func getPemCert(certificateStringURL string, token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get(certificateStringURL)

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("unable to find appropriate key")
		return cert, err
	}

	return cert, nil
}
