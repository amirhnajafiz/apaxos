package metrics

// Metrics is a struct that is used to track
// system's performance for each node.
type Metrics struct {
	tps      float64   // number of transactions per second
	latency  []float64 // a single transaction handling time
	failures int       // number of failed transactions
}

// NewMetrics return a new metrics instance.
func NewMetrics() *Metrics {
	return &Metrics{
		tps:      0,
		failures: 0,
		latency:  make([]float64, 0),
	}
}

// Metrics methods.
func (m *Metrics) ObserveTPS(period float64) {
	if m.tps == 0 {
		m.tps = period
	}

	tmp := float64(1 / period)
	m.tps = float64((m.tps + tmp) / 2)
}

func (m *Metrics) ObserveLatency(period float64) {
	m.latency = append(m.latency, period)
}

func (m *Metrics) IncFailure() {
	m.failures++
}

func (m *Metrics) Export() map[string]interface{} {
	return map[string]interface{}{
		"latency":                 m.latency,
		"failures":                m.failures,
		"transactions per second": m.tps,
	}
}
