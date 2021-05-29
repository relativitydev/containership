package processor

import (
	"fmt"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
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
	ref, err := name.NewRepository(repository)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Failed to parse repository %s", repository))
	}

	tags, err := remote.List(ref, remote.WithAuth(
		&authn.Basic{
			Username: creds.Username,
			Password: creds.Password,
		},
	))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Failed to list tags in repository %s", repository))
	}

	return tags, nil
}
