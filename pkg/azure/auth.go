package azure

import (
	"os"
	"strings"

	acrMgmtPlane "github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/2019-05-01/containerregistry"
	"github.com/Azure/go-autorest/autorest"
	"github.com/relativitydev/containership/pkg/utils"
)

// NewACRAuthorizer creates an Authorizer configured for AzureContainerRegistry
func NewACRAuthorizer() autorest.Authorizer {
	return autorest.NewBasicAuthorizer(os.Getenv("AZURE_CLIENT_ID"), os.Getenv("AZURE_CLIENT_SECRET"))
}

// GetImportSourceCredentials gets the source credentials for imports from Azure
func GetImportSourceCredentials(loginURI string) *acrMgmtPlane.ImportSourceCredentials {
	var username string

	var password string

	registryName := utils.Before(loginURI, ".")
	config, _ := utils.GetRegistryConfig(registryName)

	if (config.Username == "" || config.Password == "") && strings.Contains(loginURI, "azurecr") { //nolint
		username = os.Getenv("AZURE_CLIENT_ID")
		password = os.Getenv("AZURE_CLIENT_SECRET")
	} else if config.Username != "" && config.Password != "" {
		username = config.Username
		password = config.Password
	} else {
		return nil
	}

	return &acrMgmtPlane.ImportSourceCredentials{
		Username: &username,
		Password: &password,
	}
}
