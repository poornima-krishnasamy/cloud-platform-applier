package config_test

import (
	"os"
	"testing"

	"github.com/poornima-krishnasamy/cloud-platform-applier/pkg/config"
)

func TestNewEnvPipelineTestConfig(t *testing.T) {
	defer os.Unsetenv("PIPELINE_STATEREGION")
	os.Setenv("PIPELINE_STATEREGION", "eu-west-2")
	config := config.NewEnvPipelineConfig()

	if config.StateRegion != "eu-west-2" {
		t.Errorf("Expected value for environment variable PIPELINE_STATEREGION not set")
	}

}
