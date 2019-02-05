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
	return fmt.Sprintf("%v/.cache/wallhaven/", home)
}

func getConfigPath(flag string) string {
	if flag != "" {
		return flag
	}
	return fmt.Sprintf("%v/.config/wallhaven/", home)
}
