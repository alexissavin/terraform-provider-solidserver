package solidserver

import (
)

type Config struct {
  Host       string
  Username   string
  Password   string
  SSLVerify  bool
}

func (c *Config) APIClient() (*Config, error) {

  return c, nil
}
