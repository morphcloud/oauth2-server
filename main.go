package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	rest_v1 "github.com/morphcloud/oauth2-server/internal/handlers/http/rest/v1"
)

var (
	appName, lisAddr string
)

func configureEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln(err)
	}

	appName = os.Getenv("APP_NAME")
	if appName == "" {
		appName = "order-service"
	}

	lisAddr = os.Getenv("PORT")
	if lisAddr == "" {
		log.Fatalln("Port is not set")
	} else {
		lisAddr = ":"+lisAddr
	}
}

func waitForShutdown(srv http.Server, l *log.Logger) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	l.Println("Graceful shutdown:", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func main() {
	configureEnv()

	l := log.New(os.Stdout, os.Getenv("APP_NAME") + " ", log.LstdFlags)

	router := mux.NewRouter()

	router.HandleFunc("/oauth/login", rest_v1.LoginHandler).Methods("POST")

	srv := http.Server{
		Addr:              lisAddr,
		Handler:           router,
		ErrorLog:          l,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       20 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			l.Fatalln(err)
		}
	}()
	l.Printf("%s has been started on %s\n", appName, lisAddr)

	waitForShutdown(srv, l)
}
