package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	HttpRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"endpoints", "method", "status"},
	)

	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_requests_duration_seconds",
			Help: "Request latency in seconds",
		},
		[]string{"endpoint"},
	)

	CacheHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total cache hits",
		},
	)

	CacheMiss = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_miss_total",
			Help: "Total cache miss",
		},
	)

	CacheDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "cache_duration_seconds",
			Help: "Cache access duration",
		},
		[]string{"result"}, //hit or miss
	)

	CacheWrites = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_writes_total",
			Help: "Total cache writes",
		},
	)

	ExternalAPIDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name: "external_apo_duration_seconds",
			Help: "External API latency",
		},
	)

	RedisHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "redis_hits_total",
			Help: "Total Redis cache hits",
		},
	)

	RedisMiss = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "redis_miss_total",
			Help: "Total Redis cache miss",
		},
	)

	RedisDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "redis_duration_seconds",
			Help: "Redis access duration",
		},
		[]string{"result"}, //hit or miss
	)

	RedisWrites = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "redis_writes_total",
			Help: "Total Redis cache writes",
		},
	)
)

func Init() {
	prometheus.MustRegister(
		HttpRequestTotal,
		HttpRequestDuration,
		CacheHits,
		CacheMiss,
		CacheDuration,
		CacheWrites,
		ExternalAPIDuration,
		RedisHits,
		RedisMiss,
		RedisDuration,
		RedisWrites,
	)
}
