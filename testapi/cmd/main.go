package main

import (
	"net/http"
	"os"
	"os/signal"

	"github.com/getapid/apid-cli/common/log"
	"github.com/getapid/apid-cli/testapi/internal/handler"
	"github.com/gin-gonic/gin"
)

const (
	listenAddr = "0.0.0.0:80"
)

func main() {
	log.Init(-1)
	defer log.L.Sync()

	router := gin.Default()
	bind(handler.NewGinHandler(), router)

	server := http.Server{
		Addr:    listenAddr,
		Handler: router,
	}

	go log.L.Error(server.ListenAndServe())
	defer server.Close()
	log.L.Info("listening on ", listenAddr)

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
	log.L.Info("shutting down...")
}

func bind(h handler.GinHandler, r gin.IRouter) {
	r.GET("/health", h.HandleHealthCheck)
	r.POST("/echo", h.HandleEcho)
	r.POST("/auth", h.HandleLogin)

	authenticatedBeers := r.Group("beer")
	authenticatedBeers.Use(h.AuthMiddleware)
	authenticatedBeers.GET("", h.ListBeers)
	authenticatedBeers.GET("/:id", h.GetBeer)
}
