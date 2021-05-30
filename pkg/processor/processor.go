package processor

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ahmetb/go-linq"
	"github.com/google/go-containerregistry/pkg/v1/remote/transport"
	"github.com/pkg/errors"
	containershipappv1beta2 "github.com/relativitydev/containership/api/v1beta2"
)

func Run(client RegistryClient, images []containershipappv1beta2.Image, registries []RegistryCredentials) error {
	for _, imageConfig := range images {
		currentSourceRepo := imageConfig.SourceRepository

		for _, creds := range registries {
			// List tags
			targetTags, err := client.listTags(imageConfig.TargetRepository, creds)
			if err != nil {
				v, ok := err.(*transport.Error)
				if ok && v.StatusCode == http.StatusNotFound {
					// if error is because the repository isn't found, just move on. We'll create it later.
				} else {
					return errors.Wrap(err, "Failed to list tags")
				}
			}

			// Determine which tags should be imported and which should be deleted
			tagsToDelete, tagsToImport := populateTagArrays(targetTags, imageConfig.SupportedTags)

			for _, tag := range tagsToImport {
				imageSourceFQN := currentSourceRepo + ":" + tag

				imageDestinationFQDN := fmt.Sprintf("%s/%s:%s", creds.Hostname, imageConfig.TargetRepository, tag)

				// Pull images to add
				img, err := client.pull(imageSourceFQN, creds)
				if err != nil {
					return errors.Wrapf(err, "Failed to pull image %s", imageSourceFQN)
				}

				// Push to destination
				if err := client.push(imageDestinationFQDN, img, creds); err != nil {
					return errors.Wrapf(err, "Failed to push image %s", imageDestinationFQDN)
				}
			}

			println(strings.Join(tagsToDelete, ","))

			// Set the next regsitry hop for the next loop
			currentSourceRepo = creds.Hostname + "/" + imageConfig.TargetRepository
		}
	}

	return nil
}

// populateTagArrays creates returns tags to be removed and tags to be imported
func populateTagArrays(targetTags []string, supportedTags []string) ([]string, []string) {
	toAdd := make([]string, 0)
	linq.From(supportedTags).Except(linq.From(targetTags)).ToSlice(&toAdd)

	toDelete := make([]string, 0)
	linq.From(targetTags).Except(linq.From(supportedTags)).ToSlice(&toDelete)

	return toDelete, toAdd
}
