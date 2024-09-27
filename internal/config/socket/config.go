package socket

// Config of the node's socket interface.
// This interface is being used for making rpcs.
type Config struct {
	// Port is the socket port of the interface.
	Port int `koanf:"port"`
	// Timeout is the socket timeout value.
	Timeout int `koanf:"timeout"`
}
