package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type EnvPipelineConfig struct {
	StateBucket    string `required:"true"`
	StateKeyPrefix string `required:"true"`
	StateLockTable string `required:"true"`
	StateRegion    string `required:"true"`
	Cluster        string `required:"true"`
	RepoPath       string `required:"true"`
	NumRoutines    int    `default:"2"`
}

// NewEnvPipelineConfig sets values for EnvPipelineConfig by fetching it from environment variable
// and returns a pointer to a the config
func NewEnvPipelineConfig() *EnvPipelineConfig {
	var env EnvPipelineConfig
	err := envconfig.Process("pipeline", &env)
	if err != nil {
		log.Fatal(err.Error())
	}
	return &env
}
