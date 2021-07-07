package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/nanaagyirbrown/memrizr/handler"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main(){
	log.Println("Starting server...")

	router := gin.Default()

	handler.NewHandler(&handler.Config{
		R: router,
	})

	srv := &http.Server{
		Addr: ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to initialize server: %v\n", err)
		}
	}()

	log.Printf("Listening on port %v\n", srv.Addr)

	//Wait for kill signal of channel
	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	//This blocks until a signal is passed into the quit channel
	<-quit

	//The context is used to inform the server it has 5 seconds to finish the request its currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	// Shutdown server
	log.Printf("Shutting down server....")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}
}
