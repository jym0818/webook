package metric

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type MiddlewareBuilder struct {
	Namespace  string
	Subsystem  string
	Name       string
	Help       string
	InstanceID string
}

func (m *MiddlewareBuilder) Build() gin.HandlerFunc {
	summary := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:   m.Namespace,
		Subsystem:   m.Subsystem,
		Name:        m.Name + "_resp_time",
		Help:        m.Help,
		ConstLabels: map[string]string{"instance_id": m.InstanceID},
		Objectives:  map[float64]float64{0.5: 0.01, 0.9: 0.01, 0.99: 0.001, 0.999: 0.0001},
	}, []string{"method", "pattern", "status"})

	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: m.Namespace,
		Subsystem: m.Subsystem,
		Name:      m.Name + "_active_req",
		Help:      m.Help,
		ConstLabels: map[string]string{
			"instance_id": m.InstanceID,
		},
	})

	prometheus.MustRegister(summary)
	prometheus.MustRegister(gauge)
	return func(c *gin.Context) {
		start := time.Now()
		gauge.Inc()
		defer func() {

			duration := time.Since(start)
			gauge.Dec()
			pattern := c.FullPath()
			if pattern == "" {
				pattern = "unknown"
			}
			summary.WithLabelValues(c.Request.Method, pattern, strconv.Itoa(c.Writer.Status())).Observe(float64(duration.Milliseconds()))
		}()
		c.Next()
	}
}
