# Containership

Containership is a kubernetes operator that automates image management responsibilities. 

Features include:
- pushing and pulling images into multiple container registries
- delete old image tags

Coming soon:
- conditionally promote image based on gates
- scan images for security vulnerabilities

## Table of Contents
- [Getting Started](#Getting-Started)
- [Deploying Containership](#deploying-containership)
- [Releases](#Releases)
- [Contributing](#contributing)
   - [Building & deploying locally](#building--deploying-locally)
   - [Testing](#Testing)
   
## Getting Started

Containership requires two custom resources - ContainerManagementObject (CMO) and RegistriesConfig. First, let's configure the RegistriesConfig. 
```yaml
apiVersion: containership.app/v1beta2
kind: RegistriesConfig
metadata:
  name: registriesconfig-sample
spec:
  registries:
    - name: dockerhub-relativitydev
      hostname: index.docker.io
      secretName: registries-secret
    - name: gcr-helloworld
      hostname: gcr.io
      secretName: registries-secret
```
There are two registries defined, each with a unique name `dockerhub-relativitydev` and `gcr-helloworld`. The `hostname` is where the registry is hosted. `secretName` references the name of Kubernetes secret where the registry's authentication credentials can be found. **The order the registries are listed is the same order images will be promoted.**

Next, we'll make a CMO.
```yaml
apiVersion: containership.app/v1beta2
kind: ContainerManagementObject
metadata:
  name: containermanagementobject-sample
spec:
  images:
    - sourceRepository: busybox # if domain and namespace aren't specified, "docker.io/library" is default
      targetRepository: relativitydev/busybox
      supportedTags:
        - glibc
        - latest
    - sourceRepository: gcr.io/google_containers/pause
      supportedTags:
        - 3.2
        - latest
```
There are two images to be managed, _busybox_ (or _docker.io/library/busybox_) and *gcr.io/google_containers/pause*. With _busybox_, the tags _glibc_ and _latest_ will be pulled from DockerHub and pushed to `dockerhub-relativitydev`. Then _busybox_ will be pulled from `dockerhub-relativitydev` to `gcr-helloworld`. The repository name will be `relativitydev/busybox` as defined by `targetRepository`. 

The same thing will happen for *gcr.io/google_containers/pause*, but `targetRepository` isn't defined, so Containership will use the same repository name as `sourceRepository` -- `google_containers/pause`.

Finally, we need to make a Kubernetes secret to securely store registry credentials. In this example, we'll create one secret for multiple credentials. It is base64 encoded.
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: registries-secret
type: kubernetes.io/dockerconfigjson
data:
  # This is a fake secret - not sensitive
  .dockerconfigjson: ewogICJhdXRocyI6IHsKICAgICJkb2NrZXJodWItcmVsYXRpdml0eWRldiI6IHsKICAgICAgImF1dGgiOiAiVkdocGMwbHpUbTkwVW1WaGJEcFRkWEJsY2taaGEyVlRaV055WlhRPSIKICAgIH0KICB9Cn0=
```

Here is what `.dockerconfigjson` looks like decrypted.
```json
{
  "auths": {
    "dockerhub-relativitydev": {
      // This is a fake secret - not sensitive
      "auth": "VGhpc0lzTm90UmVhbDpTdXBlckZha2VTZWNyZXREb2NrZXI="
    },
    "gcr-helloworld": {
      // This is a fake secret - not sensitive
      "auth": "VGhpc0lzTm90UmVhbDpTdXBlckZha2VTZWNyZXRHb29nbGU="
    }
  }
}
```
`dockerhub-relativitydev` and `gcr-helloworld` map to the `secretName` property in the RegistriesConfig. Make sure the names match or the credentials won't be found.

## Deploying Containership

Only one `RegistriesConfig` should be deployed per cluster running the operator. Deploying the `RegistriesConfig`, the `kubernetes.io/dockerconfigjson` secret and the operator together is a simple approach.

Multiple `ContainerManagementObjects` can be deployed a cluster. Two common common setups are
- Deploy one CMO for all images to manage
- Deploy one CMO per repository

CMOs are flexible so you can organize how to manage images for anything.

## Releases

You can find the latest releases [here](https://github.com/relativitydev/containership/releases).

## Contributing

You can find contributing guide [here](./CONTRIBUTING.md).

### Building & deploying locally
Learn how to build & deploy Containership locally [here](./BUILD.md).

### Testing
Learn how to improve testing for Containership [here](./TEST.md).


