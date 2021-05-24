# Release Process

The release process of a new version of Containership involves the following:

## 0. Prerequisites

Look at the last released version in the releases page: https://github.com/relativitydev/containership/releases
For example: currently it is 2.0.0
The next version will thus be 2.1.0

## 1. Changelog

Provide a new section in `CHANGELOG.md` for the new version that is being released along with the new features, patches and deprecations it introduces.

It should not include every single change but solely what matters to our customers, for example issue template that has changed is not important.

## 2. Create Containership release on GitHub

Creating a new release in the releases page (https://github.com/relativitydev/containership/releases) will trigger a GitHub workflow which will create a new image with the latest code and tagged with the next version (in this example 2.1.0).

Containership Deployment YAML file (eg. containership-2.1.0.yaml) is also automatically created and attached to the Release as part of the workflow.

> Note: The Docker Hub repo with all the different images can be seen here: https://hub.docker.com/r/relativitydev/containership/tags

### Release template

Every release should use the following template to create the GitHub release.

> 💡 Don't forget to update the version in the template

Here's the template:

```markdown
We are happy to release Containership <INSERT-CORRECT-VERSION> 🎉

Here are some highlights:

- <list highlights>

Learn how to deploy Containership by reading [our documentation](https://github.com/relativitydev/containership/docs/INSERT-CORRECT-VERSION/deploy/).

### New

- <list items>

### Improvements

- <list items>

### Breaking Changes

- <list items>

### Other

- <list items>
```

## 3. Update Helm Charts

Update the `version` and `appVersion` in our [chart definition](https://github.com/relativitydev/containership/blob/master/helm/Chart.yaml).