package processor

import (
	"net/http"
	"strings"

	"github.com/ahmetb/go-linq"
	"github.com/google/go-containerregistry/pkg/v1/remote/transport"
	"github.com/pkg/errors"
	containershipappv1beta2 "github.com/relativitydev/containership/api/v1beta2"
)

func Run(client RegistryClient, images []containershipappv1beta2.Image, registries []RegistryCredentials) error {
	for _, imageConfig := range images {
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

			// determine which tags should be imported and which should be deleted
			tagsToDelete, tagsToImport := populateTagArrays(targetTags, imageConfig.SupportedTags)

			println(strings.Join(tagsToImport, ","))
			println(strings.Join(tagsToDelete, ","))
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
