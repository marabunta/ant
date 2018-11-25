package ant

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-yaml/yaml"
)

// Parse command line options and configuration file
type Parse struct {
	Flags
}

// Parse parse the command line flags
func (p *Parse) Parse(fs *flag.FlagSet) (*Flags, error) {
	fs.BoolVar(&p.Flags.Version, "v", false, "Print version")
	fs.StringVar(&p.Flags.Configfile, "c", "", "`ant.yml` configuration file")
	fs.StringVar(&p.Flags.ID, "id", "", "optional `ID` to be used instead of an UUID")
	fs.StringVar(&p.Flags.Home, "home", "", "path to `directory` to store the configuration (default $HOME/.marabunta)")
	fs.IntVar(&p.Flags.GRPC, "grpc", 1415, "marabunta gRPC `port` (default 1415)")
	fs.IntVar(&p.Flags.HTTP, "http", 8000, "marabunta HTTP `port` (default 8000)")
	fs.StringVar(&p.Flags.TLSCA, "tls.ca", "", "Path to TLS Certificate Authority (`CA`)")
	fs.StringVar(&p.Flags.TLSCrt, "tls.crt", "", "Path to TLS `certificate`")
	fs.StringVar(&p.Flags.TLSKey, "tls.key", "", "Path to TLS `private key`")

	err := fs.Parse(os.Args[1:])
	if err != nil {
		return nil, err
	}
	return &p.Flags, nil
}

func (p *Parse) parseYml(file string) (*Config, error) {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	// set defaults
	var cfg = Config{
		HTTPPort: 8000,
		GRPCPort: 1415,
	}
	if err := yaml.Unmarshal(f, &cfg); err != nil {
		return nil, fmt.Errorf("unable to parse YAML file %q %s", file, err)
	}
	return &cfg, nil
}

// Usage prints to standard error a usage message
func (p *Parse) Usage(fs *flag.FlagSet) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options...]\n\n", os.Args[0])
		var flags []string
		fs.VisitAll(func(f *flag.Flag) {
			flags = append(flags, f.Name)
		})
		for _, v := range flags {
			f := fs.Lookup(v)
			s := fmt.Sprintf("  -%s", f.Name)
			name, usage := flag.UnquoteUsage(f)
			if len(name) > 0 {
				s += " " + name
			}
			if len(s) <= 4 {
				s += "\t"
			} else {
				s += "\n    \t"
			}
			s += usage
			fmt.Fprintf(os.Stderr, "%s\n", s)
		}
	}
}

// ParseArgs parse command arguments
func (p *Parse) ParseArgs(fs *flag.FlagSet) (*Config, error) {
	flags, err := p.Parse(fs)
	if err != nil {
		return nil, err
	}

	// if -v
	if flags.Version {
		return nil, nil
	}

	// if true, create a certificate and ask marabunta server to sign it
	var needCertificate bool

	home := flags.Home
	if home == "" {
		// if no home path defined use $HOME/.marabunta
		home, err = GetUserSdir()
		if err != nil {
			return nil, err
		}
		if err := os.MkdirAll(home, os.ModePerm); err != nil {
			return nil, err
		}
	}

	// if -c
	if flags.Configfile != "" {
		if !isFile(flags.Configfile) {
			return nil, fmt.Errorf("cannot read file: %q, use (\"%s -h\") for help", flags.Configfile, os.Args[0])
		}

		// parse the `run.yml` file
		cfg, err := p.parseYml(flags.Configfile)
		if err != nil {
			return nil, err
		}

		// TLS CA
		if cfg.TLS.CA != "" {
			if !isFile(cfg.TLS.CA) {
				return nil, fmt.Errorf("cannot read TLS CA file: %q, use (\"%s -h\") for help", cfg.TLS.CA, os.Args[0])
			}
		} else {
			needCertificate = true
		}

		// TLS certificate
		if cfg.TLS.Crt != "" {
			if !isFile(cfg.TLS.Crt) {
				return nil, fmt.Errorf("cannot read TLS crt file: %q, use (\"%s -h\") for help", cfg.TLS.Crt, os.Args[0])
			}
		} else {
			needCertificate = true
		}

		// TLS KEY
		if cfg.TLS.Key != "" {
			if !isFile(cfg.TLS.Key) {
				return nil, fmt.Errorf("cannot read TLS Key file: %q, use (\"%s -h\") for help", cfg.TLS.Key, os.Args[0])
			}
		} else {
			needCertificate = true
		}

		if needCertificate {
			err := createCertificate(home)
			if err != nil {
				return nil, err
			}
		}

		return cfg, nil
	}

	if fs.NFlag() < 1 {
		// TODO
		// check how to deal with the needCertificate and createCertificate to avoid duplicating code
		return nil, fmt.Errorf("missing options, use (\"%s -h\") for help", os.Args[0])
	}

	// create new cfg if not using -c
	cfg := new(Config)

	if flags.GRPC != 0 {
		cfg.GRPCPort = flags.GRPC
	}

	if flags.HTTP != 0 {
		cfg.HTTPPort = flags.HTTP
	}

	tls := TLS{}

	// TLS CA
	if flags.TLSCA != "" {
		if !isFile(flags.TLSCA) {
			return nil, fmt.Errorf("cannot read file: %q, use (\"%s -h\") for help", flags.TLSCA, os.Args[0])
		}
		tls.CA = flags.TLSCA
	} else {
		needCertificate = true
	}

	// TLS certificate
	if flags.TLSCrt != "" {
		if !isFile(flags.TLSCrt) {
			return nil, fmt.Errorf("cannot read file: %q, use (\"%s -h\") for help", flags.TLSCrt, os.Args[0])
		}
		tls.Crt = flags.TLSCrt
	} else {
		needCertificate = true
	}

	// TLS KEY
	if flags.TLSKey != "" {
		if !isFile(flags.TLSKey) {
			return nil, fmt.Errorf("cannot read file: %q, use (\"%s -h\") for help", flags.TLSKey, os.Args[0])
		}
		tls.Key = flags.TLSKey
	} else {
		needCertificate = true
	}

	cfg.TLS = tls

	return cfg, nil
}
