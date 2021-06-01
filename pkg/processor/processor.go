package processor

import (
	"fmt"
	"net/http"

	"github.com/ahmetb/go-linq"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote/transport"
	"github.com/pkg/errors"
	containershipappv1beta2 "github.com/relativitydev/containership/api/v1beta2"
)

func Run(client RegistryClient, images []containershipappv1beta2.Image, registries []RegistryCredentials) error {
	for _, imageConfig := range images {
		currentSourceRepo := imageConfig.SourceRepository

		err := setTargetRepository(&imageConfig.TargetRepository, currentSourceRepo)
		if err != nil {
			return errors.Wrapf(err, "Source repository name parsing error %s", currentSourceRepo)
		}

		for _, creds := range registries {
			// List tags
			targetTags, err := client.listTags(imageConfig.TargetRepository, creds)
			if err != nil {
				v, ok := err.(*transport.Error)
				if ok && v.StatusCode == http.StatusNotFound {
					// repository isn't found, just move on. We'll create it later.
				} else {
					return errors.Wrap(err, "Failed to list tags")
				}
			}

			// Determine which tags should be imported and which should be deleted
			tagsToDelete, tagsToImport := populateTagArrays(targetTags, imageConfig.SupportedTags)

			// Copy the supported tags from source to destination
			for _, tag := range tagsToImport {
				imageSourceFQN := currentSourceRepo + ":" + tag

				imageDestinationFQN := fmt.Sprintf("%s/%s:%s", creds.Hostname, imageConfig.TargetRepository, tag)

				err = client.copy(imageSourceFQN, imageDestinationFQN, creds)
				if err != nil {
					return errors.Wrapf(err, "Failed to copy image from %s to %s", imageSourceFQN, imageDestinationFQN)
				}
			}

			// Delete the unsupported tags
			for _, tag := range tagsToDelete {
				imageDestinationFQN := fmt.Sprintf("%s/%s:%s", creds.Hostname, imageConfig.TargetRepository, tag)

				err = client.delete(imageDestinationFQN, creds)
				if err != nil {
					return errors.Wrapf(err, "Failed to delete image from %s", imageDestinationFQN)
				}
			}

			// Set the next registry hop for the next loop
			currentSourceRepo = creds.Hostname + "/" + imageConfig.TargetRepository
		}
	}

	return nil
}

func setTargetRepository(targetRepository *string, sourceRepository string) error {
	// Set the target repository to match source if empty
	if *targetRepository == "" {
		repo, err := name.NewRepository(sourceRepository)
		if err != nil {
			return errors.Wrapf(err, "Source repository name parsing error %s", sourceRepository)
		}

		*targetRepository = repo.RepositoryStr()
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
