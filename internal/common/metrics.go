package common

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	MetricsNamespace = "olegbot"

	resultOK   = "ok"
	resultFail = "fail"
)

var (
	DefaultRegistry = prometheus.DefaultRegisterer
)

func ConvertOk(ok bool) string {
	if ok {
		return resultOK
	}

	return resultFail
}

const (
	subsystemName       = "controller"
	endpointNameLabel   = "endpoint"
	endpointStatusLabel = "status"
)

var (
	UpdateCount = promauto.With(DefaultRegistry).NewCounter(prometheus.CounterOpts{
		Namespace: MetricsNamespace,
		Subsystem: subsystemName,
		Name:      "update_count",
		Help:      "Количество обновлений (сообщений)",
	})
	HandleCount = promauto.With(DefaultRegistry).NewCounterVec(prometheus.CounterOpts{
		Namespace: MetricsNamespace,
		Subsystem: subsystemName,
		Name:      "handle_count",
		Help:      "Количество обработанных сообщений",
	}, []string{endpointNameLabel, endpointStatusLabel})
)
