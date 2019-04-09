package main

import (
	"fmt"
	"os"

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
		return fmt.Sprintf("%v/wallhaven/", cache)
	}
	return fmt.Sprintf("%v/.cache/wallhaven/", home)
}

func getConfigPath(flag string) string {
	if flag != "" {
		return flag
	}
	config := os.Getenv("XDG_CONFIG_HOME")
	if config != "" {
		return fmt.Sprintf("%v/wallhaven/", config)
	}
	return fmt.Sprintf("%v/.config/wallhaven/", home)
}
