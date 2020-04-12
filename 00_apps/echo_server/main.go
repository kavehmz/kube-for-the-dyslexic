package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var readinessState bool = true

func readinessProbe(w http.ResponseWriter, req *http.Request) {
	log.Println("readinessProbe received")
	if req.FormValue("readiness") == "false" {
		readinessState = false
	}
	if req.FormValue("readiness") == "true" {
		readinessState = true
	}

	if !readinessState {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("not ready"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ready"))
}

var livenessState bool = true

func livenessProbe(w http.ResponseWriter, req *http.Request) {
	log.Println("livenessProbe received")
	if req.FormValue("liveness") == "false" {
		livenessState = false
	}
	if req.FormValue("liveness") == "true" {
		livenessState = true
	}

	if !livenessState {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("not alive"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("alive"))
}

func echo(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "echo: %s\n", req.FormValue("message"))
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/echo", echo)
	router.HandleFunc("/ready", readinessProbe)
	router.HandleFunc("/alive", livenessProbe)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Panicf("echo exited with error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Panicf("shutdown failed: %v", err)
	}
}
