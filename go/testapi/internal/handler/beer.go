package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/getapid/apid/common/log"
)

type beer struct {
	Id    int     `json:"id"`
	Name  string  `json:"name"`
	Type  string  `json:"type"`
	Price float32 `json:"price"`
}

type beerHandler struct {
	beers  []beer
	mBeers map[int]beer
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

	h.mBeers = make(map[int]beer, len(h.beers))
	for _, b := range h.beers {
		if _, ok := h.mBeers[b.Id]; ok {
			log.L.Panicw("found two beers in config with same id; are you drunk?",
				"id", b.Id,
			)
		}

		h.mBeers[b.Id] = b
	}
}

func (h beerHandler) ListBeers(c *gin.Context) {
	response := struct {
		Beers []beer `json:"beers"`
	}{h.beers}

	c.JSON(http.StatusOK, response)
}

func (h beerHandler) GetBeer(c *gin.Context) {
	id := c.Param("id")

	i, err := strconv.Atoi(id)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	b, exists := h.mBeers[i]
	if !exists {
		c.Status(http.StatusNotFound)
		return
	}

	response := struct {
		Beers beer `json:"beer"`
	}{b}

	c.JSON(http.StatusOK, response)
}
