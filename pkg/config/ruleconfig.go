package config

import (
	"fmt"
	"os"
)

func ReadRuleConfig() []byte {
	env := os.Getenv("APP_ENV")
	if env == "" {
		panic("Provide alue for APP_ENV")
	}

	path := fmt.Sprintf("configs/%s.json", env)
	file, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("failed to load rules for %s : %v", env, err))
	}
	return file
}
