package middlewares

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

//Collect metric "Total requests per endpoint"
var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of requests.",
	},
	[]string{"path"},
)

//Collect metric "Response status on request"
var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "response_status",
		Help: "Status of HTTP response",
	},
	[]string{"path", "status"},
)

//Collect metric "Duration of response on requests"
var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_response_time_seconds",
	Help: "Duration of HTTP requests.",
}, []string{"path"})

func init() {
	var err error
	err = prometheus.Register(totalRequests)
	if err != nil {
		fmt.Printf("Error while register totalRequests metric %s\n", err)
	}
	err = prometheus.Register(httpDuration)
	if err != nil {
		fmt.Printf("Error while register httpDuration metric %s\n", err)
	}
	err = prometheus.Register(responseStatus)
	if err != nil {
		fmt.Printf("Error while register pesponseStatus metric %s\n", err)
	}
}
