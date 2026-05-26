package config

import (
	"flag"
	"os"
	"testing"
)

// resetEnvVars is convenience function to delete environment variables.
// Without running this in each test function, the config will persist between tests.
func resetEnvVars() {
	envVars := []string{"PROMINVCFG", "PROMINVQUERY", "PROMINVURL"}
	for _, v := range envVars {
		os.Unsetenv(v)
	}
}

func TestFlagsDefaults(t *testing.T) {
	resetEnvVars()
	flag.Parse()
	cfgVal := getCfgItem(flags.Config, "PROMINVCFG", "", defaultConfig)
	if cfgVal != defaultConfig {
		t.Errorf("Expected --config to contain \"%v\" but got \"%v\".", defaultConfig, cfgVal)
	}
	queryVal := getCfgItem(flags.Query, "PROMINVQUERY", "", defaultQuery)
	if queryVal != defaultQuery {
		t.Errorf("Expected --query to contain \"%v\" but got \"%v\".", defaultQuery, queryVal)
	}
	urlVal := getCfgItem(flags.URL, "PROMINVURL", "", "")
	if urlVal != "" {
		t.Errorf("Expected --url to be empty but it contained \"%s\"", urlVal)
	}
}

func TestFlagsEnv(t *testing.T) {
	resetEnvVars()
	expectedConfig := "foobarbaz"
	expectedQuery := "up{job=\"foo\"}"
	expectedURL := "http://prometheus.notmydomain:9090"
	// This needs to be set prior to doing ParseFlags()
	os.Setenv("PROMINVCFG", expectedConfig)
	os.Setenv("PROMINVQUERY", expectedQuery)
	os.Setenv("PROMINVURL", expectedURL)
	flag.Parse()
	outConfig := getCfgItem(flags.Config, "PROMINVCFG", "", defaultConfig)
	if outConfig != expectedConfig {
		t.Errorf("Expected --config to contain \"%s\" but got \"%s\".", expectedConfig, outConfig)
	}
	outQuery := getCfgItem(flags.Query, "PROMINVQUERY", "", defaultQuery)
	if outQuery != expectedQuery {
		t.Errorf("Expected --query to contain \"%s\" but got \"%s\".", expectedQuery, outQuery)
	}
	outURL := getCfgItem(flags.URL, "PROMINVURL", "", "")
	if outURL != expectedURL {
		t.Errorf("Expected --url to be %s but got %s", expectedURL, outURL)
	}
}

func TestEmptyConfig(t *testing.T) {
	resetEnvVars()
	testFile, err := os.CreateTemp("", "testcfg")
	if err != nil {
		t.Fatalf("Unable to create TempFile: %v", err)
	}
	defer os.Remove(testFile.Name())
	cfgFile := new(Config)
	cfgFile.WriteConfig(testFile.Name())
	// An instance of Flags is required to feed ParseConfig
	cfg, err := ParseConfig(testFile.Name())
	if err == nil {
		t.Error("ParseConfig should have returned errNoPromURL")
	} else if err != errNoPromURL {
		t.Errorf("Expected errNoPromURL but got %v", err)
	}
	if cfg != nil {
		t.Errorf("Config is populated despite an error")
	}
}

func TestConfig(t *testing.T) {
	resetEnvVars()
	testFile, err := os.CreateTemp("", "testcfg")
	if err != nil {
		t.Fatalf("Unable to create TempFile: %v", err)
	}
	defer os.Remove(testFile.Name())
	expectedQuery := "metric{job=\"fakeJob\"}"
	expectedURL := "http://prometheus.notmydomain.com:9090"
	cfgFile := new(Config)
	cfgFile.Query = expectedQuery
	cfgFile.URL = expectedURL
	cfgFile.WriteConfig(testFile.Name())
	// An instance of Flags is required to feed ParseConfig
	cfg, err := ParseConfig(testFile.Name())
	if err != nil {
		t.Fatalf("Error parsing config file: %v", err)
	}
	if cfg.Query != expectedQuery {
		t.Errorf("Unexpected Config.Query.  Expected=%s, Got=%s", expectedQuery, cfg.Query)
	}
	if cfg.URL != expectedURL {
		t.Errorf("Unexpected Config.URL.  Expected=%s, Got=%s", expectedURL, cfg.URL)
	}
}
