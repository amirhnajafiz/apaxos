package metrics

// Metrics is a module that stores a node performance.
type Metrics struct {
	latency    float64 // time spent for processing a single transaction in milliseconds
	throughput float64 // transactions per second parameter

	lNumObservations int // to count the number of latency records
	tNumObservations int // to count the number of throughput records
}

// NewMetrics returns a new metrics instance.
func NewMetrics() *Metrics {
	return &Metrics{
		latency:    0,
		throughput: 0,
	}
}

// ObserveLatency adds a new latency record and updates the average latency.
func (m *Metrics) ObserveLatency(latency float64) {
	if m.tNumObservations == 0 {
		m.latency = latency
	} else {
		// calculate cumulative moving average
		m.latency = (m.latency*float64(m.tNumObservations) + latency) / float64(m.tNumObservations+1)
	}
	m.tNumObservations++
}

// ObserveThroughput adds a new throughput record and updates the average throughput.
func (m *Metrics) ObserveThroughput(throughput float64) {
	if m.lNumObservations == 0 {
		m.throughput = throughput
	} else {
		// calculate cumulative moving average
		m.throughput = (m.throughput*float64(m.lNumObservations) + throughput) / float64(m.lNumObservations+1)
	}
	m.lNumObservations++
}

// GetValues is used to export the current metrics. (latency, throughput)
func (m *Metrics) GetValues() (float64, float64) {
	return m.latency, m.throughput
}
