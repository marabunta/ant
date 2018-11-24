package ant

// Config yaml/command line configuration
type Config struct {
	HTTPPort int `yaml:"http_port"`
	GRPCPort int `yaml:"grpc_port"`
	TLS      `yaml:"tls"`
}

// TLS configuration options
type TLS struct {
	Crt string `yaml:"crt"`
	Key string `yaml:"key"`
	CA  string `yaml:"ca"`
}
