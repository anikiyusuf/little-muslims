package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yusufaniki/muslim_tech/internal/boostrap"
	"github.com/yusufaniki/muslim_tech/internal/server"
)

const (
	shutdownTimeout = 5 * time.Second
)
func main() {
	app, err := boostrap.InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	srv := server.NewServer(app)

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	go func() {

		app.Logger.Info("Starting server on port ...")
		if err := srv.Start(); err != nil  && err != http.ErrServerClosed {
			app.Logger.Fatalf("Server encountered an error.", map[string]interface{}{"error": err})

		}
	}()
		<-shutdownChan
		app.Logger.Info("Shutdown signal received, shutting down server...")
        
		if err := srv.Shutdown(shutdownTimeout); err != nil {
			app.Logger.Errorf("Error during server shutdown.", map[string]interface{}{"error": err})
		}else {
			app.Logger.Info("Server shutdown completed gracefully.")
		
	}
	app.Logger.Info("Application stopped.")
}