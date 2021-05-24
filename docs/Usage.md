# Usage
---

### Image Promotion

As defined in the `ContainerManagementObject` (CMO), the destination registries will be synced to only include the `supportedTags`. If a tag is not found, it will be imported from `source`.

### Internal Images

These are images internal to Relativity. They are defined in a CMO like this:
```yaml
destinations:
  azurecontainerregistries:
    - name: acrreg.azurecr.io
sourceImageName: acrdev.azurecr.io/busybox
supportedTags:
  - 1.32.0
  - latest
```

### External Images

These are images pulled into private container registries from an external source like _docker.io_, _quay.io_, etc. External image sources and destination are defined in a CMO like this:

```yaml
destinations:
  azurecontainerregistries:
    - name: acrdev.azurecr.io
sourceImageName: docker.io/alpine
supportedTags:
  - edge
  - 3.12.0
```

### Promotion Rings

Containership will promote images in rings if defined in the CMO. If no ring is defined, each destination is treated as ring zero.

```yaml
destinations:
  azurecontainerregistries:
    - name: acrdev.azurecr.io
      ring: 0
    - name: acrreg.azurecr.io
      ring: 1
sourceImageName: library/busybox
supportedTags:
  - 1.32.0
  - latest
```

### Promotion Gates

There are a number of gates that can be evaluated during the image promotion lifecycle. Currently, there are two types:
1. PreImport - before image import
2. PreDelete - before image deletion

A directory of available gates can be found [here](../pkg/gates/README.md).
