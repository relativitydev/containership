package processor

import (
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
)

type RegistryClient interface {
	delete(destFQN string, creds RegistryCredentials) error
	copy(sourceFQN string, destFQN string, creds RegistryCredentials) error
	listTags(repository string, creds RegistryCredentials) ([]string, error)
}

type registryClientImpl struct {
}

func NewRegistryClient() RegistryClient {
	return registryClientImpl{}
}

// listTags returns the list of tags in a repository
func (c registryClientImpl) listTags(repository string, creds RegistryCredentials) ([]string, error) {
	return crane.ListTags(creds.Hostname+"/"+repository, crane.WithAuth(
		&authn.Basic{
			Username: creds.Username,
			Password: creds.Password,
		},
	))
}

// Copies an image from one remote registry to the other
func (c registryClientImpl) copy(sourceFQN string, destFQN string, creds RegistryCredentials) error {
	return crane.Copy(sourceFQN, destFQN, crane.WithAuth(
		&authn.Basic{
			Username: creds.Username,
			Password: creds.Password,
		},
	))
}

func (c registryClientImpl) delete(destFQN string, creds RegistryCredentials) error {
	return crane.Delete(destFQN, crane.WithAuth(
		&authn.Basic{
			Username: creds.Username,
			Password: creds.Password,
		},
	))
}
