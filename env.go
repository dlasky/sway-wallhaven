package main

import (
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

var home = os.Getenv("HOME")

func getCachePathFromCtx(c *cli.Context) string {
	cache := c.String("cache")
	return getCachePath(cache)
}

func getConfigPathFromCtx(c *cli.Context) string {
	config := c.String("config")
	return getConfigPath(config)
}

func getCachePath(flag string) string {
	if flag != "" {
		return flag
	}
	cache := os.Getenv("XDG_CACHE_HOME")
	if cache != "" {
		return filepath.Join(cache, "wallhaven/")
	}
	return filepath.Join(home, ".cache/wallhaven/")
}

func getConfigPath(flag string) string {
	if flag != "" {
		return flag
	}
	config := os.Getenv("XDG_CONFIG_HOME")
	if config != "" {
		return filepath.Join(config, "wallhaven/")
	}
	return filepath.Join(home, ".config/wallhaven/")
}
