package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type beer struct {
	Id    int     `json:"id"`
	Name  string  `json:"name"`
	Type  string  `json:"type"`
	Price float32 `json:"price"`
}

type handler struct {
	inventory []beer
}

func (s handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bytes, _ := json.Marshal(s.inventory)
	w.Write(bytes)
}

func main() {
	// create the handler with an inventory of two beers
	beerServer := handler{
		inventory: []beer{
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
		},
	}

	// start the HTTP server
	log.Fatal(http.ListenAndServe("localhost:8080", beerServer))
}
