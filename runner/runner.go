package runner

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"GoRestServer/model"
	"GoRestServer/pkg/config"
)

func Run(engine *gin.Engine) {
	addr := fmt.Sprintf("%s:%d", config.Server.ListenAddr, config.Server.Port)

	s := &http.Server{
		Addr:           addr,
		Handler:        engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		// service connections
		if err := s.ListenAndServe(); err != nil {
			log.Fatal().Msgf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Debug().Msg("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	model.Close()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal().Msgf("Server Shutdown:", err)
	}
	log.Debug().Msg("Server exiting ...")
}
