package bot

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	ReqCounter prometheus.Counter
	ErrCounter prometheus.Counter
}

func newMetrics() *Metrics {
	return &Metrics{
		ReqCounter: promauto.NewCounter(prometheus.CounterOpts{
			Name: "resonator_commands_invoked_total",
			Help: "The total number of invoked commands",
		}),
		ErrCounter: promauto.NewCounter(prometheus.CounterOpts{
			Name: "resonator_command_errors_total",
			Help: "The total number of commands errors",
		}),
	}
}
