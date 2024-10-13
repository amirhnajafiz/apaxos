package grpc

// Config contains a host and a port
// in order to start a gRPC server for each node.
// Moreover, it has a waiting_timeout and majority_timeout in milliseconds
// that waits for a response from gRPC calls to other nodes.
type Config struct {
	Host            string `koanf:"host"`
	Port            int    `koanf:"port"`
	WaitingTimeout  int    `koanf:"waiting_timeout"`
	MajorityTimeout int    `koanf:"majority_timeout"`
}
