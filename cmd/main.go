package main

import (
	"context"
	"github.com/caarlos0/env/v11"
	"github.com/crazybolillo/ldapb/internal"
	"github.com/glauth/ldap"
	"gopkg.in/yaml.v3"
	"io"
	"log/slog"
	"os"
)

type appConfig struct {
	Config string `env:"CONFIG" envDefault:"/etc/ldapb/ldapb.yaml"`
	Listen string `env:"LISTEN" envDefault:":389"`
}

func main() {
	os.Exit(run(context.Background()))
}

func run(ctx context.Context) int {
	cfg, err := env.ParseAs[appConfig]()
	if err != nil {
		slog.Error("Failed to retrieve environment variables", slog.String("reason", err.Error()))
		return 1
	}

	err = serve(ctx, cfg)
	if err != nil {
		return 1
	}

	return 0
}

func readConfig(path string) (*internal.Config, error) {
	rawConf, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer rawConf.Close()

	content, err := io.ReadAll(rawConf)
	if err != nil {
		return nil, err
	}

	var config internal.Config
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func serve(ctx context.Context, config appConfig) error {
	cfg, err := readConfig(config.Config)
	if err != nil {
		slog.Error("Failed to read configuration file", "reason", err)
		return err
	}

	handler := internal.Handler{
		Config:  cfg,
		Session: internal.NewSession(),
	}

	server := ldap.NewServer()
	server.BindFunc("", &handler)
	server.UnbindFunc("", &handler)
	server.SearchFunc("", &handler)
	server.AbandonFunc("", &handler)

	slog.Info("Starting server", "addr", config.Listen)

	return server.ListenAndServe(config.Listen)
}
