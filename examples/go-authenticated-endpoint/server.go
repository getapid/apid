package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var secretKey = []byte("very-secret-key")

var users = map[string]string{
	"john.doe": "Pa55word",
}

type Claims struct {
	// we embed the jwt standard claims for fields like expiry time and subject
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

func ListBeers(w http.ResponseWriter, r *http.Request) {
	bytes, _ := json.Marshal(beers)
	w.Write(bytes)
}

func Login(w http.ResponseWriter, r *http.Request) {
	// grab the basic auth from the request
	username, providedPassword, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// verify that the user exists and the provided password is correct
	actualPassword, ok := users[username]
	if !ok || actualPassword != providedPassword {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
		Username: username,
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func main() {
	// wire the ListBeers function to handle requests sent to /beers
	http.HandleFunc("/beers", ListBeers)
	// start the HTTP server
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
