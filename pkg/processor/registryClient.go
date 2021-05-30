package processor

import (
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

type RegistryClient interface {
	listTags(repository string, creds RegistryCredentials) ([]string, error)
}

type registryClientImpl struct {
}

func NewRegistryClient() RegistryClient {
	return registryClientImpl{}
}

// listTags returns the list of tags in a repository
func (c registryClientImpl) listTags(repository string, creds RegistryCredentials) ([]string, error) {
	ref, err := name.NewRepository(creds.Hostname + "/" + repository)
	if err != nil {
		return nil, err
	}

	tags, err := remote.List(ref, remote.WithAuth(
		&authn.Basic{
			Username: creds.Username,
			Password: creds.Password,
		},
	))
	if err != nil {
		return nil, err
	}

	return tags, nil
}
