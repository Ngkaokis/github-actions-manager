package api

import "github.com/oursky/github-actions-manager/pkg/utils/defaults"

type Config struct {
	Addr *string `toml:"addr,omitempty" validate:"omitempty,tcp_addr"`
}

func (c *Config) GetAddr() string {
	return defaults.Value(c.Addr, "127.0.0.1:8002")
}
