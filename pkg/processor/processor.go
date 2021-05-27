package processor

import (
	"github.com/heroku/docker-registry-client/registry"
	"github.com/pkg/errors"
	containershipappv1beta2 "github.com/relativitydev/containership/api/v1beta2"
)

func Run(images []containershipappv1beta2.Image, registries []RegistryCredentials) error {
	for _, imageConfig := range images {
		for _, creds := range registries {
			// 1. List tags
			tags, err := listTags(imageConfig.TargetRepository, creds)
			if err != nil {
				return errors.Wrap(err, "Failed to list tags")
			}

			println(tags)
		}
	}

	return nil
}

func listTags(repository string, creds RegistryCredentials) ([]string, error) {
	client, err := registry.New(creds.LoginURI, creds.Username, creds.Password)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create docker client")
	}

	return client.Tags(repository)
}
