package metrics

// Metrics is a module that stores a node performance.
type Metrics struct {
	latency    float64 // time spent for processing a single transaction in milliseconds
	throughput float64 // transactions per second parameter
}

// NewMetrics returns a new metrics instance.
func NewMetrics() *Metrics {
	return &Metrics{
		latency:    0,
		throughput: 0,
	}
}

// ObserveLatency sets a new record for latency by averaging two values.
func (m *Metrics) ObserveLatency(latency float64) {
	if m.latency == 0 {
		m.latency = latency
	}

	m.latency = float64((m.latency + latency) / 2)
}

// ObserveThroughput sets a new record for throughput by averaging two values.
func (m *Metrics) ObserveThroughput(throughput float64) {
	if m.throughput == 0 {
		m.throughput = throughput
	}

	m.throughput = float64((m.throughput + throughput) / 2)
}

// GetValues is used to export the current metrics. (latency, throughput)
func (m *Metrics) GetValues() (float64, float64) {
	return m.latency, m.throughput
}
