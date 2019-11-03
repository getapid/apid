package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type beer struct {
	Id    int     `json:"id"`
	Name  string  `json:"name"`
	Type  string  `json:"type"`
	Price float32 `json:"price"`
}

type beerHandler struct {
	beers []beer
}

func (h *beerHandler) init() {
	h.beers = []beer{
		{
			Id:    1,
			Name:  "Stella Artois",
			Type:  "pilsner",
			Price: 2.5,
		},
		{
			Id:    2,
			Name:  "Guinness",
			Type:  "Irish dry stout",
			Price: 3.5,
		},
		{
			Id:    3,
			Name:  "Zagorka",
			Type:  "lager",
			Price: 1.5,
		},
	}
}

func (h beerHandler) ListBeers(c *gin.Context) {
	response := struct {
		Beers []beer `json:"beers"`
	}{h.beers}

	c.JSON(http.StatusOK, response)
}
