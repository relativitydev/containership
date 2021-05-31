# Testing

Tests should be written and used whenever possible. Tests are automatically run when a PR is created amd must pass to merge.

To run tests locally
```
make test
```

## Unit Testing

This repo uses [testify](https://github.com/stretchr/testify) for assertions and [mockery](https://github.com/vektra/mockery) to mock third party dependencies through dependency injection.

## Integration Testing

TODO

## Functional Testing
You can easily spin up a local [kind](https://kind.sigs.k8s.io/) cluster for functional testing.

1. Create Kind Cluster
   ```bash
   make kind-start
   ```
2. Deploy CRDs
   ```bash
   make install
   ```
3. Create sample CRs
   ```bash
   make deploy-samples
   ```
4. Launch app with `make run` or VSCode debugger