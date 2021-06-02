# Containership

Containership is a Kubernetes operator that automates image management responsibilities. 

Features include:
- pushing and pulling images into multiple container registries
- delete old image tags

Coming soon:
- conditionally promote image based on gates
- scan images for security vulnerabilities

_Note: This is our open source version. We are continuing to bring v2.x to feature parity with our closed source version. Once this is complete, we will switch to only using the open source version._ 

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
  namespace: containership-system
spec:
  registries:
    - name: dockerhub-relativitydev
      hostname: index.docker.io
      secretName: registries-secret
    - name: gcr-helloworld
      hostname: gcr.io
      secretName: registries-secret
```
There are two registries defined, each with a unique name `dockerhub-relativitydev` and `gcr-helloworld`. The `hostname` is where the registry is hosted. `secretName` references the name of Kubernetes secret where the registry's authentication credentials can be found. **The order the registries are listed in the `RegistriesConfig` is the order Containership will promote the images.**

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

After deploying, you should see the supported tags listed in your regsitries. If there were any extra tags in the registries that are not listed in the CMO, they will be deleted.

Finally, we need to make a Kubernetes secret to securely store registry credentials. In this example, we'll create one secret for multiple credentials. It is base64 encoded.
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: registries-secret
  namespace: containership-system
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

### Helm

Looking for contributors!

### Kustomize
This repo has kustomize deployments setup in the `config` directory.

#### Install
- You can deploy using `make`
```
make install
make deploy
```
- Alternatively, you can using `kubectl` directory
```
kubectl apply -k ./config/default
```

#### Uninstall
- You can deploy using `make`
```
make undeploy
make uninstall
```
- Alternatively, you can using `kubectl` directory
```
kubectl delete -k ./config/default
```

### YAML

#### Install
If you want to try Containership on Minikube or a different Kubernetes deployment without using Helm you can still deploy it with kubectl.

- We provide sample YAML declaration which includes our CRDs and all other resources in a file which is available on the GitHub releases page. Run the following command (if needed, replace the version, in this case 2.0.0, with the one you are using):
```
kubectl apply -f https://github.com/relativiydev/containership/releases/download/v2.0.0/containership-2.0.0.yaml
```

- Alternatively you can download the file and deploy it from the local path:
```
kubectl apply -f containership-2.0.0.yaml
```

- You can also find the same YAML declarations in our /config directory on our GitHub repo if you prefer to clone it.
```
git clone https://github.com/relativitydev/containership && cd containership

VERSION=2.0.0 make deploy
```

#### Uninstall
- In case of installing from released YAML file just run the following command (if needed, replace the version, in this case 2.0.0, with the one you are using):
```
kubectl delete -f https://github.com/relativitydev/containership/releases/download/v2.0.0/containership-2.0.0.yaml
```

- If you have downloaded the file locally, you can run:
```
kubectl delete -f containership-2.0.0.yaml
```

- You would need to run these commands from within the directory of the cloned GitHub repo:
```
VERSION=2.0.0 make undeploy
```

### Best Practices

#### Orgainzing Custom Resources
Only one `RegistriesConfig` should be deployed per cluster running the operator. Deploying the `RegistriesConfig`, the `kubernetes.io/dockerconfigjson` secret and the operator together in the same namespace is a good approach.

Multiple `ContainerManagementObjects` can be deployed to a cluster. Two common setups are:
- Deploy one CMO for all images to manage
- Deploy one CMO per repository

CMOs are flexible so you can organize them however you prefer.

#### Key tips
- Don't declare the same image multiple times. Containership does not have any protections for duplicate image references.
- Keep `RegistriesConfig`, it's referenced secret(s), and containership operator in the same namespace.
- The order the registries are listed in `RegistriesConfig` is the same order images will be promoted.

## Releases

You can find the latest releases [here](https://github.com/relativitydev/containership/releases).

## Contributing

You can find contributing guide [here](./CONTRIBUTING.md).

### Building & deploying locally
Learn how to build & deploy Containership locally [here](./BUILD.md).

### Testing
Learn how to improve testing for Containership [here](./TEST.md).
