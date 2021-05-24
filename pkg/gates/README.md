# Gates

## Purpose

Gates run evaluation logic to determine if a image can continue. Gates currently can be set to run before image import or image deletion. They return `true` if the gate is _open_ and `false` if the gate is _closed_.

## Gate Directory

### Prisma Gate

Integrates with Prisma to determine if an image has vulnerabilities. If vulnerabilities exist, the gate is _closed_.

<details>
<summary>Usage Details</summary>

The metadata object must contain the following
```
"registry":   "acrdev.azurecr.io",
"repository": "busybox",
"tag":        "latest",
```

Required environmental variables
Name|Description|Possible Values|Default
---|---|---|---
`<REGISTRY>_PRISMA_URL`|Prisma console hostname||example.com
`<REGISTRY>_PRISMA_USERNAME`|Prisma user||containership-user
`<REGISTRY>_PRISMA_PASSWORD`|Prisma user secret||
`PRISMA_VULNERABILITY_LEVEL`|The vulnerability level that will close the gate|`critical`,`high`,`medium`,`low`|`low`

Replace `<REGISTRY>` with a registry name like `ACRDEV`.

</details>

### Deployed Image Gate

Determines if an image is running on any host across all clusters in an environment (dev, test, reg, prod). If it is running on at least one host, the gate is _closed_.

<details>
<summary>Usage Details</summary>

The metadata object must contain the following
```
"registry":   "acrdev.azurecr.io",
"repository": "busybox",
"tag":        "latest",
```

Required environmental variables
Name|Description|Possible Values|Default
---|---|---|---
`<REGISTRY>_PRISMA_URL`|Prisma console hostname||example.com
`<REGISTRY>_PRISMA_USERNAME`|Prisma user||containership-user
`<REGISTRY>_PRISMA_PASSWORD`|Prisma user secret||

Replace `<REGISTRY>` with a registry name like `ACRDEV`.

</details>