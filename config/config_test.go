package config_test

import (
	"os"
	"testing"
	"time"
	"tixgo/config"
)

func writeTempFile(dir, name, content string) error {
	path := dir + "/" + name
	return os.WriteFile(path, []byte(content), 0644)
}

func withTempDir(t *testing.T, fn func(tmpDir string)) {
	t.Helper()
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer os.Chdir(origDir)
	fn(tmpDir)
}

func TestLoadConfig(t *testing.T) {
	validConfig := `
app:
  name: tixgo
  environment: dev
  debug_mode: true
server:
  host: localhost
  port: 8080
  read_timeout: 10s
  write_timeout: 10s
  idle_timeout: 10s
database:
  type: postgres
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  name: tixgo_dev
  ssl_mode: disable
  max_open_conns: 10
  max_idle_conns: 5
  max_lifetime: 3600s
  max_idle_time: 3600s
`
	invalidYaml := `app: [name: tixgo` // malformed YAML
	invalidValues := `
app:
  name: ""
  environment: "invalid"
  debug_mode: true
server:
  host: ""
  port: 70000
  read_timeout: 0
  write_timeout: 0
  idle_timeout: 0
database:
  type: ""
  host: ""
  port: 0
  user: ""
  password: ""
  name: ""
  ssl_mode: "invalid"
  max_open_conns: 0
  max_idle_conns: 0
  max_lifetime: 0
  max_idle_time: 0
`

	t.Run("success", func(t *testing.T) {
		withTempDir(t, func(tmpDir string) {
			err := writeTempFile(tmpDir, "config.yaml", validConfig)
			if err != nil {
				t.Fatalf("write config: %v", err)
			}
			cfg, err := config.LoadConfig()
			if err != nil {
				t.Fatalf("LoadConfig failed: %v", err)
			}
			if cfg.App.Name != "tixgo" || cfg.Server.Port != 8080 || cfg.Database.Type != "postgres" {
				t.Errorf("unexpected config: %+v", cfg)
			}
			if cfg.Server.ReadTimeout != 10*time.Second {
				t.Errorf("expected 10s, got %v", cfg.Server.ReadTimeout)
			}
		})
	})

	t.Run("missing config file", func(t *testing.T) {
		withTempDir(t, func(_ string) {
			_, err := config.LoadConfig()
			if err == nil {
				t.Error("expected error for missing config file, got nil")
			}
		})
	})

	t.Run("invalid yaml", func(t *testing.T) {
		withTempDir(t, func(tmpDir string) {
			err := writeTempFile(tmpDir, "config.yaml", invalidYaml)
			if err != nil {
				t.Fatalf("write config: %v", err)
			}
			_, err = config.LoadConfig()
			if err == nil {
				t.Error("expected error for invalid yaml, got nil")
			}
		})
	})

	t.Run("validation error", func(t *testing.T) {
		withTempDir(t, func(tmpDir string) {
			err := writeTempFile(tmpDir, "config.yaml", invalidValues)
			if err != nil {
				t.Fatalf("write config: %v", err)
			}
			_, err = config.LoadConfig()
			if err == nil {
				t.Error("expected validation error, got nil")
			}
		})
	})

	t.Run("env override loads config.env.yaml", func(t *testing.T) {
		withTempDir(t, func(tmpDir string) {
			t.Setenv("APP_ENV", "testenv")
			// base config is valid, but override file is invalid
			err := writeTempFile(tmpDir, "config.yaml", validConfig)
			if err != nil {
				t.Fatalf("write config: %v", err)
			}
			err = writeTempFile(tmpDir, "config.testenv.yaml", invalidValues)
			if err != nil {
				t.Fatalf("write config.env: %v", err)
			}
			_, err = config.LoadConfig()
			if err == nil {
				t.Error("expected validation error from env override, got nil")
			}
		})
	})

	t.Run("env override fallback to base config if env file missing", func(t *testing.T) {
		withTempDir(t, func(tmpDir string) {
			t.Setenv("APP_ENV", "notfound")
			err := writeTempFile(tmpDir, "config.yaml", validConfig)
			if err != nil {
				t.Fatalf("write config: %v", err)
			}
			cfg, err := config.LoadConfig()
			if err != nil {
				t.Fatalf("LoadConfig failed: %v", err)
			}
			if cfg.App.Name != "tixgo" {
				t.Errorf("expected app name 'tixgo', got %s", cfg.App.Name)
			}
		})
	})

	t.Run("overide config with env variable", func(t *testing.T) {
		withTempDir(t, func(tmpDir string) {
			err := writeTempFile(tmpDir, "config.yaml", validConfig)
			if err != nil {
				t.Fatalf("write config: %v", err)
			}

			t.Setenv("APP_SERVER_PORT", "8081")
			cfg, err := config.LoadConfig()
			if err != nil {
				t.Fatalf("LoadConfig failed: %v", err)
			}
			if cfg.Server.Port != 8081 {
				t.Errorf("expected port 8081, got %d", cfg.Server.Port)
			}
		})
	})
}
