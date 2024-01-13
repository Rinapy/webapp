package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"webapp/internal/config"
)

type storer interface {
	Get(key string) (data []byte, found bool)
	GetAll() (keys []string)
}

type WebServer struct {
	http.Server
}

func InitServer(cfg config.Config, s storer) *WebServer {
	h := newHandler(s)
	r := newRouter(h)

	srv := WebServer{
		http.Server{
			Addr:    cfg.Addr,
			Handler: r}}

	return &srv
}

func (srv *WebServer) ShutdownOnSignal() {
	shutDownSignal, _ := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM, os.Interrupt,
	)

	<-shutDownSignal.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	log.Println("gracefully shutting down the web-server..")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("failed to shutdown the server gracefully, forcing exit", err)
	}
}

func (srv *WebServer) Launch() {
	log.Println("starting web-server on address:", srv.Server.Addr)
	if err := srv.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("failed to start listening", err)
	}
}
