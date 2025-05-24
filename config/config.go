package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func LoadConfig() (*AppConfig, error) {
	v := viper.New()
	setupViper(v)

	if err := loadConfigurations(v); err != nil {
		return nil, err
	}

	config, err := unmarshalConfig(v)
	if err != nil {
		return nil, err
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func loadConfigurations(v *viper.Viper) error {
	if err := loadBaseConfig(v); err != nil {
		return err
	}

	if err := loadEnvConfig(v); err != nil {
		return err
	}

	setupEnvVars(v)

	return nil
}

func setupViper(v *viper.Viper) {
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath("/app/config")
	v.SetConfigName("config")
}

func loadBaseConfig(v *viper.Viper) error {
	if err := v.ReadInConfig(); err != nil {
		if ok := errors.As(err, &viper.ConfigFileNotFoundError{}); !ok {
			return fmt.Errorf("config file not found: %w", err)
		}

		fmt.Println("Base config file not found, replying on environment variables")
	} else {
		fmt.Println("Base config file is loaded")
	}

	return nil
}

func loadEnvConfig(v *viper.Viper) error {
	env := getEnvironment()
	v.Set("environment", env)

	v.SetConfigName("config." + env)
	if err := v.MergeInConfig(); err != nil {
		if ok := errors.As(err, &viper.ConfigFileNotFoundError{}); !ok {
			return fmt.Errorf("config file not found: %w", err)
		}

		fmt.Printf("Environment config file config.%s.yaml not found, replying on base config\n", env)
	} else {
		fmt.Printf("Merge config in file config.%s.yaml to base config\n", env)
	}

	return nil
}

func getEnvironment() string {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	return env
}

func setupEnvVars(v *viper.Viper) {
	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
}

func unmarshalConfig(v *viper.Viper) (*AppConfig, error) {
	var config AppConfig
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return &config, nil
}
