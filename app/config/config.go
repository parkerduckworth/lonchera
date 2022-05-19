// Package config provides application configuration for all required environments,
// structured as Go types, read in from YAML files. Absolutely no secrets are to be
// stored in these files, or in source control in general. All secrets should be
// loaded into an environment, via remote manager, securely stored in a hosted
// platform, or simply managed locally outside of source control. Each env may have
// its variables injected differently, based on the environment's specific needs/usage.
package config

import (
	"fmt"
	"os"

	"github.com/semi-technologies/weaviate-go-client/v4/weaviate"
	"github.com/spf13/viper"
)

var Conf Config

// Config structures the environment configuration which is read
// in from a YAML file. The file contents should match the structure
// of this type
type Config struct {
	Env    string
	Server struct {
		HTTPPort     string
		ReadTimeout  string
		WriteTimeout string
	}
	Logger struct {
		Level string
	}
	Weaviate weaviate.Config
}

// Setup reads the environment file based on the application env,
// and populates a Config instance. Otherwise this function kills
// the running process if any errors occur
func Setup() {
	buildEnv := os.Getenv("GO_ENV")
	if len(buildEnv) == 0 {
		buildEnv = "local"
	}

	viper.SetConfigName(buildEnv + ".config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./env")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("config file read error")
		fmt.Println(err)
		os.Exit(1)
	}

	if err := viper.Unmarshal(&Conf); err != nil {
		fmt.Println("config file Unmarshal error")
		fmt.Println(err)
		os.Exit(1)
	}
}
