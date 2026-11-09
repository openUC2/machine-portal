package templates

import (
	"strings"

	"github.com/sargassum-world/godest/env"
)

const envPrefix = "TEMPLATES_"

type Config struct {
	Paths []string
}

func GetConfig() (c Config, err error) {
	// This is a file path specific to ImSwitch OS
	const defaultPath = ""
	c.Paths = strings.Split(env.GetString(envPrefix+"PATH", defaultPath), ":")
	if len(c.Paths) == 1 && c.Paths[0] == "" {
		c.Paths = nil
	}

	return c, nil
}
