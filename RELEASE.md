# Release

VMClarity adopted the [Go MultiMod Releaser](https://github.com/open-telemetry/opentelemetry-go-build-tools/tree/main/multimod) to create a new release for each module.

Find below the steps to create a release.

## 1. Update the new release version

Checkout to a new branch and update the version defined for VMClarity's module-set in `versions.yaml`. Please note that currently the same version is used for all methods. E.g.

```
  vmclarity:
-    version: v0.6.0
+    version: v0.7.0
```

Commit this change and verify the versioning with `make multimod-verify`.

## 2. Bump all dependencies to new release version

To update all `go.mod` files with the new release version, run `make multimod-prerelease` and review the changes in the last commit.

## 3. Tag the new release commit

## 4. Release

