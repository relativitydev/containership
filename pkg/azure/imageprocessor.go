package azure

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"

	linq "github.com/ahmetb/go-linq"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/relativitydev/containership/pkg/gates"
	"github.com/relativitydev/containership/pkg/utils"
)

var (
	log                              = ctrl.Log.WithName("azure/imageprocessor")
	allowedDestinations              = strings.Split(os.Getenv("ALLOWED_DESTINATIONS"), ",")
	handler             ImageHandler = new(ImageHandlerImpl)
	alertMessage        string
)

// ImageProcessor processes changes to Azure images
func ImageProcessor(ctx context.Context, promotion ImagePromotion) (events []corev1.Event) {
	if promotion.TargetRepository == "" {
		promotion.TargetRepository = utils.After(promotion.SourceImage, "/")
	}

	destinations, sameRingError := reorderAzureDestinations(promotion.Destinations)
	if sameRingError != nil {
		alertMessage = fmt.Sprintf("Two or more rings have the same value for the following repository: %s", promotion.TargetRepository)
		log.Error(sameRingError, alertMessage)

		events = append(events, corev1.Event{
			Type:    "Warning",
			Reason:  "DuplicateRings",
			Message: alertMessage,
		})
	}

	for i, destination := range destinations {
		if !isAllowedDestination(destination.Name) {
			continue
		}

		// Get tags from Azure
		destinationTagsList, tagError := handler.ListImageTags(ctx, "https://"+destination.Name, promotion.TargetRepository)

		if tagError != nil {
			// Catch repository not found error
			if skipError(tagError.Error(), promotion.TargetRepository) {
				destinationTagsList = make([]string, 0)
			} else {
				alertMessage = fmt.Sprintf("Error getting tags from Azure for the following repository: %s", promotion.TargetRepository)
				log.Error(tagError, alertMessage)

				events = append(events, corev1.Event{
					Type:    "Warning",
					Reason:  "AzureTagRetrieval",
					Message: alertMessage,
				})
				continue
			}
		}

		// Helper arrays to distinguish tags to be removed from tags to be imported
		toDelete, toAdd := populateTagArrays(destinationTagsList, promotion.SupportedTags)

		log.Info(fmt.Sprintf("Target Repository: %s\n\t\t\t\t\t\t\t\tDestination: %s\n\t\t\t\t\t\t\t\tExisting Tags: %s\n\t\t\t\t\t\t\t\tTags to Add: %s\n\t\t\t\t\t\t\t\tTags to Delete: %s", promotion.TargetRepository, destination.Name, strings.Join(destinationTagsList, ", "), strings.Join(toAdd, ", "), strings.Join(toDelete, ", ")))

		if i != 0 {
			promotion.SourceImage = destinations[i-1].Name + "/" + promotion.TargetRepository
		}

		events = append(events, removeTags(ctx, promotion.TargetRepository, destination, toDelete)...)
		events = append(events, importTags(ctx, promotion.SourceImage, promotion.TargetRepository, toAdd, destination)...)
	}

	return events
}

// reorderAzureDestinations reorders azure destination array based on destination.Ring
func reorderAzureDestinations(azureDestinations []PromotionDestination) ([]PromotionDestination, error) {
	var sameRingError error

	sort.SliceStable(azureDestinations, func(i, j int) bool {
		if azureDestinations[i].Ring == azureDestinations[j].Ring && azureDestinations[i].Ring > 0 {
			sameRingError = errors.New("Invalid destination priorities. Duplicate Ring value of " + strconv.Itoa(azureDestinations[i].Ring) + ".")
		}
		return azureDestinations[i].Ring < azureDestinations[j].Ring
	})

	return azureDestinations, errors.Wrap(sameRingError, "Failed to reorder Azure destinations")
}

func isAllowedDestination(name string) bool {
	if len(allowedDestinations) == 0 || allowedDestinations[0] == "" {
		return true
	}

	return linq.From(allowedDestinations).Contains(name)
}

// populateTagArrays creates helper arrays to distinguish tags to be removed from tags to be imported
func populateTagArrays(destinationTagsList []string, supportedTags []string) ([]string, []string) {
	toAdd := make([]string, 0)
	linq.From(supportedTags).Except(linq.From(destinationTagsList)).ToSlice(&toAdd)

	toDelete := make([]string, 0)
	linq.From(destinationTagsList).Except(linq.From(supportedTags)).ToSlice(&toDelete)

	return toDelete, toAdd
}

/*
Remove tags from Azure
Note: this helper function may log an error but will not return it
*/
func removeTags(ctx context.Context, targetImageName string, destination PromotionDestination, tags []string) (events []corev1.Event) {
	for _, tag := range tags {
		// Check if the image is deployed
		pass, err := isDeployedImage(ctx, destination.Name, targetImageName, tag)
		if err != nil {
			log.Error(err, "Could not evaluate deployed image gate")
		}

		if !pass {
			alertMessage = fmt.Sprintf("Evaluation of deployed image gate failed. Image is still being used.\nregistry: %s\nrepository: %s:%s", destination.Name, targetImageName, tag)
			log.Info(alertMessage)

			events = append(events, corev1.Event{
				Type:    "Warning",
				Reason:  "ImageStillDeployed",
				Message: alertMessage,
			})

			break // Move onto the next tag.
		}

		deleteTagError := handler.DeleteImage(ctx, "https://"+destination.Name, targetImageName, tag)
		if deleteTagError != nil {
			alertMessage = fmt.Sprintf("Error deleting %s from %s", tag, targetImageName)
			log.Error(deleteTagError, alertMessage)

			events = append(events, corev1.Event{
				Type:    "Warning",
				Reason:  "TagDeletionFailure",
				Message: alertMessage,
			})
		} else {
			log.Info("Removed " + tag + " from " + targetImageName)
		}
	}

	return events
}

// Add tags to Azure
// Note: this helper function may log an error but will not return it
func importTags(ctx context.Context, sourceImage string, targetRepository string, tags []string, destination PromotionDestination) (events []corev1.Event) {
	if destination.SubscriptionID == "" {
		destination.SubscriptionID = os.Getenv("AZURE_REGISTRY_SUBSCRIPTION_ID")
	}

	if destination.ResourceGroup == "" {
		destination.ResourceGroup = os.Getenv("AZURE_REGISTRY_RESOURCE_GROUP")
	}

	var wg sync.WaitGroup

	wg.Add(len(tags))

	for _, tag := range tags {
		go func(tag string) {
			importSourceImage := utils.After(sourceImage, "/") + ":" + tag
			importSourceURI := utils.Before(sourceImage, "/")

			internalImport := strings.Contains(importSourceURI, "azurecr")
			pass := true

			if internalImport {
				// Scan the image for vulnerabilities
				var err error
				pass, err = isVulnerableImage(ctx, importSourceURI, targetRepository, tag)

				if err != nil {
					log.Error(err, "Could not evaluate prisma gate")
				}
			}

			if pass {
				// Import the image
				_, importTagError := handler.ImportImage(ctx, destination.SubscriptionID, destination.ResourceGroup, utils.Before(destination.Name, "."), importSourceURI, importSourceImage, targetRepository, []string{tag})
				if importTagError != nil {
					alertMessage = fmt.Sprintf("Import error\ndest acr %s\nimportsource %s\nsource image %s\ntarget repository %s", utils.Before(destination.Name, "."), importSourceURI, importSourceImage, targetRepository)
					log.Error(importTagError, alertMessage)

					events = append(events, corev1.Event{
						Type:    "Warning",
						Reason:  "ImageImportError",
						Message: alertMessage,
					})
				} else {
					log.Info("Successful import",
						"source", importSourceImage,
						"destination", destination.Name)
				}
			} else {
				alertMessage = fmt.Sprintf("Evaluation of prisma gate failed.\nregistry: %s\nrepository: %s", importSourceURI, importSourceImage)
				log.Info(alertMessage)

				events = append(events, corev1.Event{
					Type:    "Warning",
					Reason:  "PrismaGateFailed",
					Message: alertMessage,
				})
			}

			// letting the WaitGroup know that the thread is complete
			defer wg.Done()
		}(tag)
	}

	// wait for Done() to be executed on all threads in the WaitGroup
	wg.Wait()

	return events
}

func isVulnerableImage(ctx context.Context, registry string, repository string, tag string) (bool, error) {
	// Scan the image for vulnerabilities
	gateHandler := gates.NewGateHandler()
	gate, err := gateHandler.GetGate("prisma", map[string]string{
		"registry":   registry,
		"repository": repository,
		"tag":        tag,
	})

	if err != nil {
		log.Error(err, "Could not create prisma gate")
	}

	result, err := gate.Evaluate(ctx)

	return result, errors.Wrap(err, "Prisma gate evaluation failed from an error")
}

func isDeployedImage(ctx context.Context, registry string, repository string, tag string) (bool, error) {
	// Scan the image for vulnerabilities
	gateHandler := gates.NewGateHandler()
	gate, err := gateHandler.GetGate("deployedImage", map[string]string{
		"registry":   registry,
		"repository": repository,
		"tag":        tag,
	})

	if err != nil {
		log.Error(err, "Could not create deployed image gate")
	}

	result, err := gate.Evaluate(ctx)

	return result, errors.Wrap(err, "Deployed image gate evaluation failed from an error")
}

// skipError returns true if error message matches preTargetImageMessage + targetImageName + postTargetImageMessage
func skipError(err string, targetImageName string) bool {
	return strings.Contains(err, "repository \\\""+targetImageName+"\\\" is not found")
}

// Helper functions for used only in testing

// TagsMatch returns true if the tags in Azure match the supported tags
func TagsMatch(ctx context.Context, targetImageName string, supportedTags []string, destinations []PromotionDestination) (bool, error) {
	for _, destination := range destinations {
		destinationTagsList, err := handler.ListImageTags(ctx, "https://"+destination.Name, targetImageName)
		if err != nil {
			return false, errors.Wrap(err, "Unable to list tags from Azure")
		}

		sort.Strings(destinationTagsList)
		sort.Strings(supportedTags)

		if reflect.DeepEqual(destinationTagsList, supportedTags) {
			return true, nil
		}
	}

	return false, nil
}

// DoesTagExistAzure returns true if the targetTag exists in Azure
func DoesTagExistAzure(ctx context.Context, targetImageName string, destination PromotionDestination, targetTag string) (bool, error) {
	destinationTagsList, err := handler.ListImageTags(ctx, "https://"+destination.Name, targetImageName)
	if err != nil {
		return false, errors.Wrap(err, "Unable to list tags from Azure")
	}

	return linq.From(destinationTagsList).Contains(targetTag), nil
}
