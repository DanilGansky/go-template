package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sirupsen/logrus"
)

// Run runs server.ListenAndServe and then stops it gracefully
func Run(s *http.Server, timeout time.Duration, log *logrus.Logger) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go func() {
		log.Infof("server started at: %s", s.Addr)
		if err := s.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatalf("server error: %s", err.Error())
			}
		}
	}()

	<-signals
	log.Info("shutting down...")
	return s.Shutdown(ctx)
}
