package azure

import (
	"context"
	"fmt"
	"strings"

	acrMgmtPlane "github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/2019-05-01/containerregistry"
	acrDataPlane "github.com/Azure/azure-sdk-for-go/services/preview/containerregistry/runtime/2019-08-15-preview/containerregistry"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	linq "github.com/ahmetb/go-linq"
	retry "github.com/cenkalti/backoff/v4"
	"github.com/pkg/errors"
)

// ImageHandler is an interface that allows for the creation of methods to mock those in this file
type ImageHandler interface {
	ImportImage(ctx context.Context, subscriptionID string, resourceGroup string, acrName string, externalRegistryURI string, sourceImageName string, targetImageName string, targetTags []string) (ar autorest.Response, err error)
	DeleteImage(ctx context.Context, loginURI string, imageName string, tag string) error
	ListImageTags(ctx context.Context, loginURI string, imageName string) ([]string, error)
}

// ImageHandlerImpl is the standard implementation of image handler. Any time you would like to use the following functions externally, you must first create a new ImageHandlerImpl instance.
type ImageHandlerImpl struct {
}

// ImportImage pulls an image from the externalRegistryURI into an Azure Container Registry. Only one image tag is pulled at a time and must be
// appended to sourceImageName.
// - context: httpContext
// - subscriptionID: azure subscription ID of the ACR
// - resourceGroup: resource group name of the ACR (ex: k8s-acr)
// - acrName: name of the ACR (ex: acrdev)
// - externalRegistryURI: where image is to be imported from (ex: quay.io)
// - sourceImageName: image to pull from external registry (w/ namespace & tag) (ex: prometheus/prometheus:v1.15.2)
// - targetImageName: image to push from destination registry (w/ namespace) (ex: cloud-engineering/core/prometheus/prometheus)
// - targetTags: array of tags to use for the imported image. Note: this is only applied to the destination ACR. (ex: ["latest", "v1.15.2"])
func (handler *ImageHandlerImpl) ImportImage(ctx context.Context, subscriptionID string, resourceGroup string, acrName string, externalRegistryURI string, sourceImageName string, targetImageName string, targetTags []string) (ar autorest.Response, err error) {
	auth, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return ar, errors.Wrap(err, "Failed to authenticate using environment variables")
	}

	managementURI := "https://management.azure.com"
	if strings.Contains(acrName, "gov") {
		managementURI = "https://management.usgovcloudapi.net"
	}

	client := acrMgmtPlane.NewRegistriesClientWithBaseURI(managementURI, subscriptionID)
	client.RetryAttempts = 1
	client.Authorizer = auth

	importSource := acrMgmtPlane.ImportSource{
		RegistryURI: &externalRegistryURI,
		SourceImage: &sourceImageName,
	}

	importSource.Credentials = GetImportSourceCredentials(externalRegistryURI)

	// append tag to image name. This is how Azure wants it.
	requestTags := make([]string, len(targetTags))
	for i, tag := range targetTags {
		requestTags[i] = fmt.Sprintf("%s:%s", targetImageName, tag)
	}

	// retry the image import with an exponential backoff
	var future acrMgmtPlane.RegistriesImportImageFuture

	err = retry.Retry(
		func() error {
			future, err = client.ImportImage(ctx, resourceGroup, acrName, acrMgmtPlane.ImportImageParameters{
				Source:     &importSource,
				TargetTags: &requestTags,
			})

			if err != nil {
				errorType := err.(autorest.DetailedError)
				statusCode := errorType.StatusCode

				// only retrying on 429's for now.
				// otherwise, Permanent will notify Retry to stop retrying
				if statusCode != 429 {
					err = retry.Permanent(err)
				}
			}

			return errors.Wrap(err, "Azure import image error.")
		}, retry.NewExponentialBackOff())

	if err != nil {
		return ar, errors.Wrap(err, "failed to import image")
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return ar, fmt.Errorf("cannot get the image import response: %v", err)
	}

	return future.Result(client)
}

// DeleteImage removes an image from Azure Container Registry
// - context: httpContext
// - loginURI: the acr login uri (ex: https://acrdev.azurecr.io)
// - imageName: name of the image
// - tag: image tag
func (handler *ImageHandlerImpl) DeleteImage(ctx context.Context, loginURI string, imageName string, tag string) error {
	auth := NewACRAuthorizer()
	client := acrDataPlane.NewTagClient(loginURI)
	client.RetryAttempts = 1
	client.Authorizer = auth

	_, err := client.Delete(ctx, imageName, tag)

	return errors.Wrap(err, "failed to get delete image")
}

// ListImageTags returns a list of all the tags of an image in the ACR
// - context: httpContext
// - loginURI: the acr login uri (ex: https://acrdev.azurecr.io)
// - imageName: name of the image
func (handler *ImageHandlerImpl) ListImageTags(ctx context.Context, loginURI string, imageName string) ([]string, error) {
	// Note: we have to make two azure calls. 1. to get the tags 2. to get the manifests. We compare the two to make sure the tags actually have a manifest. Sometimes, a tag exists without a manifest, which returns 404 when pulling the image. We can't only use show-manifests because sometimes the manifest response doesn't have all the tags listed. This error was caught by integration tests for containership-e2e-test/busybox:1.32.0
	auth := NewACRAuthorizer()
	tagsClient := acrDataPlane.NewTagClient(loginURI)
	tagsClient.RetryAttempts = 1
	tagsClient.Authorizer = auth

	tagsResponse, err := tagsClient.GetList(ctx, imageName, "", nil, "", "")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get list of tags")
	}

	var tags []string
	if tagsResponse.Tags != nil {
		linq.From(*tagsResponse.Tags).SelectT(func(t acrDataPlane.TagAttributesBase) string {
			return *t.Name
		}).ToSlice(&tags)
	}

	manifestsClient := acrDataPlane.NewManifestsClient(loginURI)
	manifestsClient.RetryAttempts = 1
	manifestsClient.Authorizer = auth

	manifestsResponse, err := manifestsClient.GetList(ctx, imageName, "", nil, "")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get list of tags")
	}

	if len(*manifestsResponse.ManifestsAttributes) > 0 {
		for _, manifest := range *manifestsResponse.ManifestsAttributes {
			if manifest.Tags != nil && len(*manifest.Tags) > 0 {
				tags = append(tags, *manifest.Tags...)
			}
		}
	}

	var mergedTags []string

	linq.From(tags).Distinct().ToSlice(&mergedTags)

	return mergedTags, nil
}
