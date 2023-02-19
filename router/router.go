package router

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Router struct {
	echo *echo.Echo
}

func New() *Router {
	return &Router{echo: echo.New()}
}

func (r *Router) Serve(port string) {

	r.echo.Use(middleware.Logger())
	r.echo.GET("/ws", handleSocket)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	srv := http.Server{
		Addr:    port,
		Handler: r.echo,
	}
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	<-signals
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Error Shutting down! %s", err.Error())
	}
	log.Println("Shutdown successfully")
}
