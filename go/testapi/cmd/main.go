package main

import (
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/iv-p/apid/common/log"
	"github.com/iv-p/apid/testapi/internal/handler"
)

const (
	listenAddr = "localhost:8080"
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

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
	log.L.Info("shutting down...")
}

func bind(h handler.GinHandler, r gin.IRouter) {
	r.GET("/health", h.HandleHealthCheck)
	r.POST("/echo", h.HandleEcho)
}
