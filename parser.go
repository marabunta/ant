package ant

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-yaml/yaml"
)

// Parse command line options and configuration file
type Parse struct {
	Flags
}

// Parse parse the command line flags
func (p *Parse) Parse(fs *flag.FlagSet) (*Flags, error) {
	fs.BoolVar(&p.Flags.Version, "v", false, "Print version")
	fs.BoolVar(&p.Flags.Start, "start", false, "Start client")
	fs.StringVar(&p.Flags.Configfile, "c", "", "`ant.yml` configuration file")
	fs.StringVar(&p.Flags.S3, "s3", "", "`access:secret` keys if empty it will use environment variables ANT_S3_ACCESS_KEY and ANT_S3_SECRET_KEY, <command>")

	err := fs.Parse(os.Args[1:])
	if err != nil {
		return nil, err
	}
	return &p.Flags, nil
}

func (p *Parse) parseYml(file string, cfg *Config) (*Config, error) {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(f, &cfg); err != nil {
		return nil, fmt.Errorf("unable to parse YAML file %q %s", file, err)
	}
	return cfg, nil
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
	var (
		home            string
		needCertificate bool
		cfg             = &Config{
			Marabunta: "marabunta.host",
			HTTPPort:  8000,
			GRPCPort:  1415,
			TLS: TLS{
				ServerName: "marabunta",
			},
		}
	)

	// if -c
	if flags.Configfile != "" {
		if !isFile(flags.Configfile) {
			return nil, fmt.Errorf("cannot read file: %q, use (\"%s -h\") for help", flags.Configfile, os.Args[0])
		}

		// parse the `run.yml` file
		cfg, err := p.parseYml(flags.Configfile, cfg)
		if err != nil {
			return nil, err
		}

		// Home
		if cfg.Home == "" {
			home, err := GetHome()
			if err != nil {
				return nil, err
			}
			cfg.Home = home
		}

		// TLS certificate
		if cfg.TLS.Crt != "" {
			if !isFile(cfg.TLS.Crt) {
				return nil, fmt.Errorf("cannot read TLS crt file: %q, use (\"%s -h\") for help", cfg.TLS.Crt, os.Args[0])
			}
		} else {
			cfg.TLS.Crt = filepath.Join(cfg.Home, "ant.crt")
			needCertificate = true
		}

		// TLS KEY
		if cfg.TLS.Key != "" {
			if !isFile(cfg.TLS.Key) {
				return nil, fmt.Errorf("cannot read TLS Key file: %q, use (\"%s -h\") for help", cfg.TLS.Key, os.Args[0])
			}
		} else {
			cfg.TLS.Key = filepath.Join(cfg.Home, "ant.key")
			needCertificate = true
		}

		if needCertificate {
			err := createCertificate(cfg)
			if err != nil {
				return nil, err
			}
		}

		return cfg, nil
	}

	home, err = GetHome()
	if err != nil {
		return nil, err
	}

	cfg.Home = home

	cfg.TLS.Crt = filepath.Join(cfg.Home, "ant.crt")
	if !isFile(cfg.TLS.Crt) {
		needCertificate = true
	}
	cfg.TLS.Key = filepath.Join(cfg.Home, "ant.key")
	if !isFile(cfg.TLS.Key) {
		needCertificate = true
	}

	if needCertificate {
		if err := createCertificate(cfg); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}
