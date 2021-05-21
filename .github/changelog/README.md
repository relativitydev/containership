# Generate changelog with git-chglog

Github Action for creating a CHANGELOG.md file based on semver and conventional commits.

## Usage
### Pre-requisites
Create a workflow .yml file in your repositories .github/workflows directory. An example workflow is available below. For more information, reference the GitHub Help Documentation for Creating a workflow file.

Further more you need to have [git-chlog](https://github.com/git-chglog/git-chglog) configured and have the configuration added to your git repository.

### Inputs
 - `config_dir`: git-chglog configuration directory. Default: `.ghglog`
 - `filename`: Filename to write the changelog to. Default: `CHANGELOG_TMP.md`
 - `tag`: Optional, Generate changelog only for this tag.

### Outputs
 - `changelog`: Changelog content
 - `filepath`: A filepath to the changelog if you cannot take the content output

### Example workflow
When a git tag is added to a repo, this workflow will execute to generate a changelog and output it for use in other workflow steps.

```yaml
name: Go Release

on:
  push:
    tags:
      - "*"

jobs:
  goreleaser:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v2.3.4
        with:
          fetch-depth: 0
      - name: Generate Release Notes
        id: changelog
        uses: ./changelog
      - name: Set up Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.14
      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v2.5.0
        with:
          version: latest
          args: release --rm-dist --release-notes=${{ steps.changelog.outputs.filepath }}
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/nuuday/github-changelog-action/tags).
