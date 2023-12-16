package common

import (
	"github.com/prometheus/client_golang/prometheus"
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
