package qbot

import "log/slog"

type Bot struct {
	*Config
}

type Config struct {
	Logger *slog.Logger
	Dialector
}

func (c *Config) Apply(config *Config) error {
	if config != c {
		*config = *c
	}
	return nil
}

type Option interface {
	Apply(config *Config) error
}

func New(dialector Dialector, opts ...Option) (*Bot, error) {
	config := &Config{}
	for _, opt := range opts {
		if opt != nil {
			if confErr := opt.Apply(config); confErr != nil {
				return nil, confErr
			}
		}
	}

	if config.Logger == nil {
		config.Logger = slog.Default()
	}
	config.Dialector = dialector

	return &Bot{
		Config: config,
	}, nil
}
