package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds all Prometheus metrics
type Metrics struct {
	// HTTP metrics
	HTTPRequestsTotal   *prometheus.CounterVec
	HTTPRequestDuration *prometheus.HistogramVec
	HTTPRequestSize     *prometheus.HistogramVec
	HTTPResponseSize    *prometheus.HistogramVec

	// Database metrics
	DBConnectionsOpen     prometheus.Gauge
	DBConnectionsInUse    prometheus.Gauge
	DBConnectionsIdle     prometheus.Gauge
	DBQueryDuration       *prometheus.HistogramVec
	DBTransactionsTotal   *prometheus.CounterVec

	// Cache metrics
	CacheHitsTotal   *prometheus.CounterVec
	CacheMissesTotal *prometheus.CounterVec
	CacheOperations  *prometheus.HistogramVec

	// Queue metrics
	QueueMessagesProduced *prometheus.CounterVec
	QueueMessagesConsumed *prometheus.CounterVec
	QueueMessagesDuration *prometheus.HistogramVec

	// Business metrics
	CampaignsTotal     prometheus.Counter
	DonationsTotal     *prometheus.CounterVec
	DonationAmount     *prometheus.HistogramVec
	UsersTotal         prometheus.Gauge
	ActiveSessionsTotal prometheus.Gauge
}

// New creates a new Metrics instance
func New(namespace string) *Metrics {
	return &Metrics{
		// HTTP metrics
		HTTPRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "http_requests_total",
				Help:      "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
		HTTPRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "http_request_duration_seconds",
				Help:      "HTTP request duration in seconds",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
		HTTPRequestSize: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "http_request_size_bytes",
				Help:      "HTTP request size in bytes",
				Buckets:   prometheus.ExponentialBuckets(100, 10, 8),
			},
			[]string{"method", "path"},
		),
		HTTPResponseSize: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "http_response_size_bytes",
				Help:      "HTTP response size in bytes",
				Buckets:   prometheus.ExponentialBuckets(100, 10, 8),
			},
			[]string{"method", "path"},
		),

		// Database metrics
		DBConnectionsOpen: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "db_connections_open",
				Help:      "Number of open database connections",
			},
		),
		DBConnectionsInUse: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "db_connections_in_use",
				Help:      "Number of database connections in use",
			},
		),
		DBConnectionsIdle: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "db_connections_idle",
				Help:      "Number of idle database connections",
			},
		),
		DBQueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "db_query_duration_seconds",
				Help:      "Database query duration in seconds",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"operation"},
		),
		DBTransactionsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "db_transactions_total",
				Help:      "Total number of database transactions",
			},
			[]string{"status"},
		),

		// Cache metrics
		CacheHitsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "cache_hits_total",
				Help:      "Total number of cache hits",
			},
			[]string{"cache_name"},
		),
		CacheMissesTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "cache_misses_total",
				Help:      "Total number of cache misses",
			},
			[]string{"cache_name"},
		),
		CacheOperations: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "cache_operation_duration_seconds",
				Help:      "Cache operation duration in seconds",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"operation", "cache_name"},
		),

		// Queue metrics
		QueueMessagesProduced: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "queue_messages_produced_total",
				Help:      "Total number of messages produced to queue",
			},
			[]string{"topic"},
		),
		QueueMessagesConsumed: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "queue_messages_consumed_total",
				Help:      "Total number of messages consumed from queue",
			},
			[]string{"topic", "status"},
		),
		QueueMessagesDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "queue_message_processing_duration_seconds",
				Help:      "Queue message processing duration in seconds",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"topic"},
		),

		// Business metrics
		CampaignsTotal: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "campaigns_total",
				Help:      "Total number of campaigns created",
			},
		),
		DonationsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "donations_total",
				Help:      "Total number of donations",
			},
			[]string{"status", "payment_method"},
		),
		DonationAmount: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "donation_amount",
				Help:      "Donation amount distribution",
				Buckets:   prometheus.ExponentialBuckets(10000, 2, 15), // 10k to ~300M IDR
			},
			[]string{"payment_method"},
		),
		UsersTotal: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "users_total",
				Help:      "Total number of registered users",
			},
		),
		ActiveSessionsTotal: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "active_sessions_total",
				Help:      "Total number of active sessions",
			},
		),
	}
}

// RecordHTTPRequest records an HTTP request
func (m *Metrics) RecordHTTPRequest(method, path string, status int, duration time.Duration, requestSize, responseSize int64) {
	m.HTTPRequestsTotal.WithLabelValues(method, path, string(rune(status))).Inc()
	m.HTTPRequestDuration.WithLabelValues(method, path).Observe(duration.Seconds())
	m.HTTPRequestSize.WithLabelValues(method, path).Observe(float64(requestSize))
	m.HTTPResponseSize.WithLabelValues(method, path).Observe(float64(responseSize))
}

// RecordDBQuery records a database query
func (m *Metrics) RecordDBQuery(operation string, duration time.Duration) {
	m.DBQueryDuration.WithLabelValues(operation).Observe(duration.Seconds())
}

// RecordCacheOperation records a cache operation
func (m *Metrics) RecordCacheOperation(cacheName, operation string, hit bool, duration time.Duration) {
	if hit {
		m.CacheHitsTotal.WithLabelValues(cacheName).Inc()
	} else {
		m.CacheMissesTotal.WithLabelValues(cacheName).Inc()
	}
	m.CacheOperations.WithLabelValues(operation, cacheName).Observe(duration.Seconds())
}

// RecordQueueMessageProduced records a message produced to queue
func (m *Metrics) RecordQueueMessageProduced(topic string) {
	m.QueueMessagesProduced.WithLabelValues(topic).Inc()
}

// RecordQueueMessageConsumed records a message consumed from queue
func (m *Metrics) RecordQueueMessageConsumed(topic, status string, duration time.Duration) {
	m.QueueMessagesConsumed.WithLabelValues(topic, status).Inc()
	m.QueueMessagesDuration.WithLabelValues(topic).Observe(duration.Seconds())
}
