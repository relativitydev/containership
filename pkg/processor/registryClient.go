package processor

import (
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

type RegistryClient interface {
	listTags(repository string, creds RegistryCredentials) ([]string, error)
	pull(source string, creds RegistryCredentials) (v1.Image, error)
	push(imageFullyQualifiedName string, img v1.Image, creds RegistryCredentials) error
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

// pull retrieves an image from a given image and credentials
// imageFullyQualifiedName: (ex: gcr.io/google-containers/pause:v1.2.3)
func (c registryClientImpl) pull(imageFullyQualifiedName string, creds RegistryCredentials) (v1.Image, error) {
	ref, err := name.ParseReference(imageFullyQualifiedName)
	if err != nil {
		return nil, err
	}

	img, err := remote.Image(ref, remote.WithAuth(
		&authn.Basic{
			Username: creds.Username,
			Password: creds.Password,
		},
	))
	if err != nil {
		return nil, err
	}

	return img, nil
}

// push uploads an image to the tagged resgitry
// imageFullyQualifiedName: (ex: gcr.io/google-containers/pause:v1.2.3)
func (c registryClientImpl) push(imageFullyQualifiedName string, img v1.Image, creds RegistryCredentials) error {
	ref, err := name.ParseReference(imageFullyQualifiedName)
	if err != nil {
		return err
	}

	return remote.Write(ref, img, remote.WithAuth(
		&authn.Basic{
			Username: creds.Username,
			Password: creds.Password,
		},
	))
}
