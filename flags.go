package ant

// Flags available command flags
type Flags struct {
	Configfile string
	GRPC       int
	HTTP       int
	TLSCA      string
	TLSCrt     string
	TLSKey     string
	Version    bool
}
