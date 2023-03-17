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
	"github.com/chinathaip/chatify/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("---Hello Chatify---")
	cfg := config.All()

	db, err := gorm.Open(postgres.Open(cfg.DBConnection), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to db: %v", err.Error())
	}

	msgService := &service.MessageModel{DB: db}
	chatService := &service.ChatModel{DB: db}
	userService := &service.UserModel{DB: db}
	participantService := &service.ParticipantModel{DB: db}

	hub := hub.New(chatService, msgService, userService, participantService)
	handler := router.NewHandler(chatService, msgService, userService, participantService)

	ctx, cancel := context.WithCancel(context.Background())
	go hub.Init(ctx)

	e := router.RegRoute(hub, handler)

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
