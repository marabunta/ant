package ant

// Flags available command flags
type Flags struct {
	Configfile string
	GRPC       int
	HTTP       int
	Home       string
	ID         string
	Start      bool
	TLSCA      string
	TLSCrt     string
	TLSKey     string
	Version    bool
}
