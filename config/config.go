// config provides flag and configuration file reading for satinv
package config

import (
	"flag"
)

// Flags are the command line flags
type Flags struct {
	Config  string
	Debug   bool
	List    bool
	Refresh bool
}

// ParseFlags transcribes command line flags into a struct
func ParseFlags() *Flags {
	f := new(Flags)
	// Config file
	flag.StringVar(&f.Config, "config", "", "Config file")
	flag.BoolVar(&f.Debug, "debug", false, "Write logoutput to stderr")
	flag.BoolVar(&f.List, "list", false, "Produce a full inventory to stdout")
	flag.BoolVar(&f.Refresh, "refresh", false, "Force a cache refresh")
	flag.Parse()
	return f
}
