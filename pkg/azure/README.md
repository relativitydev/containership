# Azure

The image import and deletion logic for Azure Container Registries

## Environmental Variables

| Variable | Description | Default Value |
|----------|-------------|---------------|
| `AZURE_ENVIRONMENT` | `AzureUSGovernmentCloud` or `AzurePublicCloud` | `AzurePublicCloud` |
| `AZURE_TENANT_ID` | Required for authentication | `` |
| `AZURE_CLIENT_ID` | Required for authentication | `` |
| `AZURE_CLIENT_SECRET` | Required for authentication | `<super-secret>` |
| `AZURE_REGISTRY_SUBSCRIPTION_ID` | Global subscription to be used by default for all ACRs | `` |
| `AZURE_REGISTRY_RESOURCE_GROUP` | Global resource group to be used by default for all ACRs | `` |
| `ALLOWED_DESTINATIONS` | Comma separated list of destination images can be imported into | `acrdev.azurecr.io,acrtest.azurecr.io,acrreg.azurecr.io` |

You can include these envvars in your `.env` file at the project root.
