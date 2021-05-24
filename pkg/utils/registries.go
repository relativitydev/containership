package utils

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/relativitydev/containership/pkg/config"
)

// GetRegistryConfig returns the configuration for a container registry
func GetRegistryConfig(registryName string) (s config.Registry, err error) {
	err = envconfig.Process(registryName, &s)
	return
}
