package grpc

// Config contains a host and a port in order to start a gRPC server for each node.
// Moreover, it has request and majority timeoutes in milliseconds.
type Config struct {
	Host            string `koanf:"host"`
	Port            int    `koanf:"port"`
	RequestTimeout  int    `koanf:"request_timeout"`
	MajorityTimeout int    `koanf:"majority_timeout"`
}
