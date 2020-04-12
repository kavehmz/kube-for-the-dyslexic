package main

import (
	"context"
	"fmt"
	"io/ioutil"
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

func relay(w http.ResponseWriter, req *http.Request) {
	var netClient = &http.Client{
		Timeout: time.Second * 2,
	}

	resp, err := netClient.Get(fmt.Sprintf("http://%s/echo?message=%s", req.FormValue("echo_server"), req.FormValue("message")))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(fmt.Sprintf("error: %v", err)))
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(fmt.Sprintf("body error: %v", err)))
		return
	}

	hostname, _ := os.Hostname()
	fmt.Fprintf(w, "Relay-Host: %s\nRelay-To: %s\nReply:\n%s", hostname, req.FormValue("echo_server"), string(body))
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/relay", relay)
	router.HandleFunc("/ready", readinessProbe)
	router.HandleFunc("/alive", livenessProbe)

	srv := &http.Server{
		Addr:    ":8081",
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
