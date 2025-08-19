package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

const (
	localHost        = "127.0.0.1"
	defaultProxyPort = 8080
	defaultSocksPort = 9050
)

func main() {

	runtime.MemProfileRate = 0

	log.SetOutput(os.Stdout)

	proxyPort := flag.Int("proxyport", defaultProxyPort, "Proxy port")

	flag.Parse()

	proxyAddr := localHost + ":" + strconv.Itoa(*proxyPort)

	httpTransport := getTransport()

	client := &http.Client{
		Transport: httpTransport,
		Timeout:   120 * time.Second,
	}

	proxyHandler := getHandler(client)

	server := &http.Server{
		Addr:    proxyAddr,
		Handler: http.HandlerFunc(proxyHandler),
	}

	go func() {
		log.Printf("Starting proxy on %s", proxyAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server exited: %v", err)
		}
	}()

	createPidFile()
	defer removePidFile()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-stop
	log.Println("Shutting down proxy...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Shutdown failed: %v", err)
	}

	log.Println("Proxy stopped")

}

func createPidFile() {
	err := PidFileCreate()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "ERR Failed to create pid file: %v\n", err)
	}
}

func removePidFile() {
	err := PidFileRemove()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "ERR Failed to remove pid file: %v\n", err)
	}
}
