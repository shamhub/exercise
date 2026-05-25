package envconfigreader

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/matryer/resync"
	"github.com/shamhub/exercise/types"
)

// Package loads environment variable from configs/*.env files
// and override by variables set through shell(CD phase)
type envConfig struct {
	defaultFile  string
	overrideFile string
}

func NewEnvConfig(confLocation string) *envConfig {
	overridefile, defaultFile := readConfig(confLocation)
	return &envConfig{
		overrideFile: overridefile,
		defaultFile:  defaultFile,
	}
}

var once resync.Once

func readConfig(confLocation string) (overrideFile, defaultFile string) {
	once.Do(func() {
		checkPath(confLocation)
		overrideFile, defaultFile := getFileLocations(confLocation)

		err := godotenv.Load(defaultFile)
		if err != nil {
			panic(err)
		}

		err = godotenv.Overload(overrideFile)
		if err != nil {
			panic(err)
		}
	})
	return
}

func checkPath(confLocation string) {
	if _, err := os.Stat(confLocation); os.IsNotExist(err) {
		panic(err.Error())
	}
}

func getFileLocations(confLocation string) (overrideFile, defaultFile string) {
	env := os.Getenv(types.DEPLOYMENT_ENV_VAR)
	switch env {
	case types.LOCAL_ENV,
		types.DEV_ENV,
		types.TEST_ENV,
		types.STAGE_ENV,
		types.PRODUCTION_ENV:
		defaultFile = confLocation + "/.env"
		overrideFile = confLocation + "/." + env + ".env"
	default:
		err_msg := fmt.Sprintf("BALANCER_ENV value %q is invalid", env)
		panic(err_msg)
	}
	return
}
