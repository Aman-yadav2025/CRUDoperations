package main

import (
	"CRUDapp/internal/config"
	"CRUDapp/internal/http/handlers/student"
	"CRUDapp/internal/storage/sqlite"
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//step1 load the config
	cfg := config.MustLoad()
	//step2 stepup dataBase

	storage, err := sqlite.New(cfg)

	if err != nil {
		log.Fatalf("database is fucked %s", err)
	}
	slog.Info("storage Intialized",slog.String("env",cfg.Env),slog.String("version","1.0.0"))

	//setup routers
	router := http.NewServeMux() // intiate router

	//create endpoints
	router.HandleFunc("POST /api/students", student.New(storage))

	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))

	router.HandleFunc("GET /api/students", student.GetList(storage))

	router.HandleFunc("PATCH /api/students/{id}",student.UpdateStudent(storage))
	//setup server
	//start the server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	//start listening

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		slog.Info("Starting Server on", slog.String("address", cfg.HttpServer.Addr))
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("server intialzation failed")
		}
	}()

	<-done // till there is no value in done channel stop here indifinatly
	//stop the server

	slog.Info("Shuting down the server at :", slog.String("address", cfg.HttpServer.Addr))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdown server ", slog.String("error", err.Error()))
	}

	slog.Info("Server ShitDown Successfully")
}
