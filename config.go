package ant

// Config yaml/command line configuration
type Config struct {
	Home     string `yaml:"home"`
	HTTPPort int    `yaml:"marabunta_http_port"`
	GRPCPort int    `yaml:"marabunta_grpc_port"`
	TLS      `yaml:"tls"`
}

// TLS configuration options
type TLS struct {
	Crt string `yaml:"crt"`
	Key string `yaml:"key"`
	CA  string `yaml:"ca"`
}
