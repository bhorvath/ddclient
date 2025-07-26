package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Service provides application configuration handling operations.
type Service interface {
	// BuildConfig returns a config.App.
	// Returns an error if the config file path has been specified,
	// but cannot be read.
	BuildConfig() (*App, error)
	// SaveConfig persists config.App.
	// Returns an error if no config file path has been specified.
	SaveConfig() error
}

type service struct {
	args *Args
}

// NewService returns a service for handling application configuration.
func NewService(a *Args) Service {
	return &service{args: a}
}

func (s *service) BuildConfig() (*App, error) {
	cfg := &App{}
	if s.args.ConfigFilePath != "" {
		data, err := ioutil.ReadFile(s.args.ConfigFilePath)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(data, cfg)
	}

	if s.args.Domain != "" {
		cfg.Domain = s.args.Domain
	}
	if s.args.Type != "" {
		cfg.Type = s.args.Type
	}
	if s.args.Name != "" {
		cfg.Name = s.args.Name
	}
	if s.args.APIKey != "" {
		cfg.APIKey = s.args.APIKey
	}
	if s.args.SecretKey != "" {
		cfg.SecretKey = s.args.SecretKey
	}
	if err := s.validateConfig(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (s *service) validateConfig(cfg *App) error {
	var e []string
	if cfg.Domain == "" {
		e = append(e, "domain not set")
	}
	if cfg.Type == "" {
		e = append(e, "type not set")
	}
	if cfg.Name == "" {
		fmt.Println("Name not set - modifying root domain record")
	}
	if cfg.APIKey == "" {
		e = append(e, "apikey not set")
	}
	if cfg.SecretKey == "" {
		e = append(e, "secretkey not set")
	}
	if e != nil {
		return fmt.Errorf("Validation failed: %s", strings.Join(e, ", "))
	}
	return nil
}

func (s *service) SaveConfig() error {
	if s.args.ConfigFilePath == "" {
		return errors.New("No config filename specified")
	}

	// Read existing config (ignore file not found)
	data, err := ioutil.ReadFile(s.args.ConfigFilePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	savedCfg := App{}
	json.Unmarshal(data, &savedCfg)

	// Add any given options
	if s.args.Domain != "" {
		savedCfg.Domain = s.args.Domain
	}
	if s.args.Type != "" {
		savedCfg.Type = s.args.Type
	}
	if s.args.Name != "" {
		savedCfg.Name = s.args.Name
	}
	if s.args.APIKey != "" {
		savedCfg.APIKey = s.args.APIKey
	}
	if s.args.SecretKey != "" {
		savedCfg.SecretKey = s.args.SecretKey
	}

	// Save the updated configs
	d, _ := json.Marshal(savedCfg)
	ioutil.WriteFile(s.args.ConfigFilePath, d, 0644)
	return nil
}
