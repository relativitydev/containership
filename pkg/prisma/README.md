# Prisma

This package wraps the Prisma API into a convenient Go client. It should have no dependencies from any other internal packages.

## Required Environmental Variables

```
<REGISTRY>_PRISMA_CONSOLE_URL=example.com
<REGISTRY>_PRISMA_USERNAME=containership-user
<REGISTRY>_PRISMA_PASSWORD=<super-secret-password>
```

Replace `<REGISTRY>` with a registry name like `ACRDEV`. You can include these envvars in your `.env` file at the project root.
