package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/chinathaip/chatify/config"
	"github.com/chinathaip/chatify/hub"
	"github.com/chinathaip/chatify/router"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("---Hello Chatify---")
	cfg := config.All()
	ctx, cancel := context.WithCancel(context.Background())
	h := hub.New()
	go h.Init(ctx)

	db, err := gorm.Open(postgres.Open(cfg.DBConnection), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to db: %v", err.Error())
	}
	e := router.RegRoute(h, db)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	srv := http.Server{
		Addr:    cfg.Server.Port,
		Handler: e,
	}
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	<-signals
	cancel() //close hub
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Error Shutting down! %s", err.Error())
	}
	log.Println("Shutdown successfully")
}
