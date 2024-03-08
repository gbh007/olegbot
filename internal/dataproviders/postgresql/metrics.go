package postgresql

import (
	"app/internal/common"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	subsystemName = "repository"
)

var (
	quoteCount = promauto.With(common.DefaultRegistry).NewGauge(prometheus.GaugeOpts{
		Namespace: common.MetricsNamespace,
		Subsystem: subsystemName,
		Name:      "quote_count",
		Help:      "Количество цитат",
	})
)
