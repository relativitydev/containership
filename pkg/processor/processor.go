package processor

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ahmetb/go-linq"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote/transport"
	"github.com/pkg/errors"
	containershipappv1beta2 "github.com/relativitydev/containership/api/v1beta2"
)

func Run(client RegistryClient, images []containershipappv1beta2.Image, registries []RegistryCredentials) error {
	for _, imageConfig := range images {
		currentSourceRepo := imageConfig.SourceRepository

		// Set the target repository to match source if empty
		if imageConfig.TargetRepository == "" {
			repo, err := name.NewRepository(currentSourceRepo)
			if err != nil {
				return errors.Wrapf(err, "Source repository name parsing error %s", currentSourceRepo)
			}

			imageConfig.TargetRepository = repo.RepositoryStr()
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

			for _, tag := range tagsToImport {
				imageSourceFQN := currentSourceRepo + ":" + tag

				imageDestinationFQN := fmt.Sprintf("%s/%s:%s", creds.Hostname, imageConfig.TargetRepository, tag)

				err = client.copy(imageSourceFQN, imageDestinationFQN, creds)
				if err != nil {
					return errors.Wrapf(err, "Failed to copy image from %s to %s", imageSourceFQN, imageDestinationFQN)
				}
			}

			println(strings.Join(tagsToDelete, ","))

			// Set the next registry hop for the next loop
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
