package main

import (
	"fmt"
	"os"
)

var home = os.Getenv("HOME")

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
