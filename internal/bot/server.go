package bot

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	shutdownTimeout   = 5 * time.Second
	readHeaderTimeout = 3 * time.Second
)

type shutdownFunc func()

func StartMetricsServer(logger Logger, reg *prometheus.Registry) shutdownFunc {
	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	server := &http.Server{
		Addr:              ":8081",
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Unexpected error on metrics server", "err", err)
		}
	}()

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			logger.Error("Unexpected error on closing metrics server", "err", err)
		}
	}
}
