# Containership

**NOTE: This version only supports promotion to Azure container registries and using Palo Alto's Prisma for vulnerability scanning. A new, cloud agnostic, version is in development. It will work with any OCI compliant container registry and support a variety of promotion gates.**

A kubernetes operator automating container image promotion and deletion across registries. It includes gates for vulnerability scanning and protection against deleting images still being used. 

Built with Operator-SDK. For information on how to use Operator-Sdk, visit their [website](https://sdk.operatorframework.io/docs/).

Containership can be used to:
- Define images configuration as code
- Import images from external and internal sources to Azure Container Registries (ACR)
- Conditionally promote images based on gates
  - Use Prisma to only promote images with no security vulnerabilities
  - Only delete images if they are not being run in a cluster
- Supports Azure consumer and government regsitries

### Table of Contents
- [Functionality and Usage](./docs/Usage.md)
- [Gate Directory](./pkg/gates/README.md)

## Testing

### Using a .env file

Environment variables are defined in `.env` in the project root. This centralized location is used for debugging and testing. It is not used in deployments.

```
GOPATH=/home/vscode/go
GO111MODULE=on
WATCH_NAMESPACE=
AZURE_GO_SDK_LOG_LEVEL=
AZURE_TENANT_ID=
AZURE_CLIENT_ID=
AZURE_CLIENT_SECRET=<super-secret>
AZURE_REGISTRY_SUBSCRIPTION_ID=
AZURE_REGISTRY_RESOURCE_GROUP=

PRISMA_VULNERABILITY_LEVEL=low

ALLOWED_DESTINATIONS="registryname.azurecr.io"

REGISTRYNAME_USERNAME=
REGISTRYNAME_PASSWORD=<super-secret>
REGISTRYNAME_PRISMA_URL=prisma.example.com
REGISTRYNAME_PRISMA_USERNAME=containership-user
REGISTRYNAME_PRISMA_PASSWORD=<super-secret>

```

### Run locally outside the cluster

Run the operator locally while hitting a real cluster's APIs. `kubectl config current-context` is the cluster you will use.

Install your CRDs into the cluster
```
make install
```

Apply your CR to the cluster
```
kubectl apply -f config/samples/containership_v1beta1_containermanagementobject.yaml
```

Run the operator against the cluster
```
make run
```

### Debugging in VSCode

You can also debug in VSCode with the ability to set breakpoints. There is already a debug configuration setup called "Containership" in `launch.json`.

Again, you need to install your CRD and create your CRs. See above.

### Unit Testing

Business logic should live in the `/pkg` directory. Each must contain unit tests. All unit tests are called via
```
make test-pkg
```

It is also possible to run/debug individual tests, files and packages using VSCode.


### E2E Testing

To start up the kind test cluster run in terminal:
```
make kind-start
```
To run all the end to end tests (and unit tests) run in terminal:
```
make test
```
Once you are done testing, close your kind cluster by running in terminal:
```
make kind-stop
```

_Note: these tests typically take about 200s (Breakdown: 115s - E2E tests, 85s - pkg unit tests)._