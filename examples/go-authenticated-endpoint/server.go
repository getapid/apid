package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var secretKey = []byte("very-secret-key")

var users = map[string]string{
	"john.doe": "Pa55word",
}

type claims struct {
	// we embed the JWT standard claims for fields like expiry time and subject
	jwt.StandardClaims
	Username string `json:"username"`
}

type beer struct {
	Id    int     `json:"id"`
	Name  string  `json:"name"`
	Type  string  `json:"type"`
	Price float32 `json:"price"`
}

var beers = []beer{
	{
		Id:    1,
		Name:  "Stella Artois",
		Type:  "pilsner",
		Price: 2.50,
	},
	{
		Id:    2,
		Name:  "Guinness",
		Type:  "Irish dry stout",
		Price: 3.50,
	},
}

func Authenticated(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// extract the token from the request's headers
		header := r.Header.Get("Authorization")
		header = strings.TrimPrefix(header, "Bearer ")

		// parse the header they sent using our secret key
		token, err := jwt.Parse(header, func(*jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		handler(w, r)
	}
}

func ListBeers(w http.ResponseWriter, r *http.Request) {
	bytes, _ := json.Marshal(beers)
	w.Write(bytes)
}

func Login(w http.ResponseWriter, r *http.Request) {
	// grab the basic auth from the request
	username, providedPassword, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// verify that the user exists and the provided password is correct
	actualPassword, ok := users[username]
	if !ok || actualPassword != providedPassword {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// create the userClaims that we are going to issue to our user
	// this will be the payload of our JWT token
	userClaims := claims{
		StandardClaims: jwt.StandardClaims{
			// we give an expiration of 1 hour
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
		Username: username,
	}

	// create the header and payload of the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, userClaims)
	// sign the token with our secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// wrap the token in a response so that we can return valid JSON
	response := struct {
		Token string `json:"access_token"`
	}{tokenString}

	// serialize the response into JSON bytes
	serializedToken, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// send back the token
	w.Write(serializedToken)
}

func main() {
	// wire the ListBeers function to handle requests sent to /beers
	http.HandleFunc("/beer", Authenticated(ListBeers))
	http.HandleFunc("/login", Login)
	// start the HTTP server
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
