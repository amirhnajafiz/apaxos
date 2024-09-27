package http

// Config of the node's HTTP handler.
type Config struct {
	// Port will be used for the node's HTTP listener.
	Port int `koanf:"port"`
}
