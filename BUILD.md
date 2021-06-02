## Building

### Quick start with [Visual Studio Code Remote - Containers](https://code.visualstudio.com/docs/remote/containers)

This helps you pull and build quickly - dev containers launch the project inside a container with all the tooling
required for a consistent and seamless developer experience.

This means you don't have to install and configure your dev environment as the container handles this for you.

To get started install [VSCode](https://code.visualstudio.com/) and the [Remote Containers extensions](
https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)

Clone the repo and launch code:

```bash
git clone git@github.com:relativitydev/containership.git
cd containership
code .
```

Once VSCode launches run `CTRL+SHIFT+P -> Remote-Containers: Reopen in container` and then use the integrated
terminal to run:

```bash
make build
```

> Note: The first time you run the container it will take some time to build and install the tooling. The image
> will be cached so this is only required the first time.

### Run Locally

#### Makefile & Operator-SDK

This project is using [Operator SDK framework](https://github.com/operator-framework/operator-sdk), make sure you have installed the right version. To check the current version used for Containership check the `OPERATOR_RELEASE_VERSION` in file [.devcontainer/Dockerfile](https://github.com/relativitydev/containership/blob/main/.devcontainer/Dockerfile).

```bash
git clone git@github.com:relativitydev/containership.git
cd containership
make build
```

If the build process fails due to some "checksum mismatch" errors, make sure that `GOPROXY` and `GOSUMDB` environment variables are set properly.
With Go installation on Fedora, for example, it could happen they are wrong.

```bash
go env GOPROXY GOSUMDB
direct
off
```

If not set properly you can just run.

```bash
go env -w GOPROXY=https://proxy.golang.org,direct GOSUMDB=sum.golang.org
```

#### Visual Studio Code

This repo also supports debugging using the Visual Studio Code debugger. This is useful for stepping through code. There is a debugging configuartion already defined in `.vscode/launch.json`. Press `F5` to run the debugger.

## Deploying

### Custom Containership locally outside of a Kubernetes cluster

The Operator SDK framework allows you to run the operator/controller locally outside the cluster without
building an image. This should help during development/debugging of Containership Operator or Gates.
> Note: This approach works only on Linux or macOS.


1. Deploy CRDs
   ```bash
   make install
   ```
2. Run the operator locally with the default Kubernetes config file present at `$HOME/.kube/config`
 and change the operator log level via `--zap-log-level=` if needed
   ```bash
   make run ARGS="--zap-log-level=debug"
   ```
   
## Miscellaneous

### Setting log levels

You can change default log levels for Containership Operator. Containership Operator uses
 [Operator SDK logging](https://sdk.operatorframework.io/docs/building-operators/golang/references/logging/) mechanism.

To change the logging level, find `--zap-log-level=` argument in Operator Deployment section in `config/manager/manager.yaml` file,
 modify its value and redeploy.

Allowed values are `debug`, `info`, `error`, or an integer value greater than `0`, specified as string

Default value: `info`

To change the logging format, find `--zap-encoder=` argument in Operator Deployment section in `config/manager/manager.yaml` file,
 modify its value and redeploy.

Allowed values are `json` and `console`

Default value: `console`