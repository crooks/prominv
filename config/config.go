// config provides flag and configuration file reading for satinv
package config

import (
	"errors"
	"flag"
	"os"

	"gopkg.in/yaml.v3"
)

type errorConst string

const (
	defaultQuery  string = "up{job=\"node\"}"
	defaultConfig string = "/etc/prominv/config.yml"
	errNoPromURL         = errorConst("prometheus API URL has not been defined")
)

func (e errorConst) Error() string {
	return string(e)
}

// Flags are the command line flags
type Flags struct {
	Config string
	List   bool
	Query  string
	URL    string
}

type Config struct {
	Query  string `yaml:"promql_query"`
	URL    string `yaml:"prometheus_url"`
	Labels struct {
		Delete  []string `yaml:"delete"`
		GroupBy string   `yaml:"group_by"`
	} `yaml:"labels"`
}

var flags *Flags

// ParseFlags transcribes command line flags into a struct
func init() {
	flags = new(Flags)
	// Config file
	flag.StringVar(&flags.Config, "config", "", "Config file")
	flag.StringVar(&flags.Query, "query", "", "The PromQL query to be executed")
	flag.StringVar(&flags.URL, "url", "", "A URL for the Prometheus API.  Something like \"https://prometheus.myorg:9090\".")
	flag.BoolVar(&flags.List, "list", false, "Produce a full inventory to stdout")
}

func DoList() bool {
	flag.Parse()
	return flags.List
}

// getCfgItem takes four inputs: A flag, an environment variable name, a config setting and a default.
// It returns the highest priority input string that is populated.  Priorities are:
// 1) Flag, 2) Environment Variable, 3) Config Setting, 4) Default
func getCfgItem(flagVal, envName, cfgVal, defaultVal string) string {
	if flagVal == "" {
		if os.Getenv(envName) != "" {
			return os.Getenv(envName)
		} else if cfgVal != "" {
			return cfgVal
		} else {
			return defaultVal
		}
	}
	return flagVal
}

func GetConfigFilename() string {
	filename := getCfgItem(flags.Config, "PROMINVCFG", "", defaultConfig)
	return filename
}

// ParseConfig expects a YAML formatted config file and populates a Config struct
func ParseConfig(fileName string) (*Config, error) {
	config := new(Config)
	file, err := os.Open(fileName)
	if err == nil {
		// The config file exists and has been opened
		defer file.Close()
		y := yaml.NewDecoder(file)
		errDecode := y.Decode(&config)
		if errDecode != nil {
			return nil, errDecode
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		// An unexpected error occurred. i.e. There was an error and it wasn't "file not found".
		return nil, err
	}
	defer file.Close()
	flag.Parse()
	// Set config defaults here
	config.Query = getCfgItem(flags.Query, "PROMINVQUERY", config.Query, defaultQuery)
	config.URL = getCfgItem(flags.URL, "PROMINVURL", config.URL, "")
	if len(config.Labels.Delete) == 0 {
		config.Labels.Delete = []string{"__name__", "instance"}
	}

	// Config failures
	if config.URL == "" {
		return nil, errNoPromURL
	}
	return config, nil
}

// WriteConfig will create a YAML formatted config file from a Config struct
func (c *Config) WriteConfig(filename string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, data, 0640)
	if err != nil {
		return err
	}
	return nil
}
