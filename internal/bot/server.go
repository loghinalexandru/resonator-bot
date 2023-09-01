package bot

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	shutdownTimeout   = 5 * time.Second
	readHeaderTimeout = 3 * time.Second
)

type shutdownFunc func()

func StartMetricsServer(botContext *Context) shutdownFunc {
	mux := http.NewServeMux()
	reg := prometheus.NewRegistry()

	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		botContext.Err,
		botContext.Req,
	)

	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	server := &http.Server{
		Addr:              ":8081",
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			botContext.Logger.Error("Unexpected error on metrics server", "err", err)
		}
	}()

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		err := server.Shutdown(ctx)

		if err != nil {
			botContext.Logger.Error("Unexpected error on closing metrics server", "err", err)
		}
	}
}
