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

To update all `go.mod` files with the new release version, run `make multimod-prerelease` and review the changes in the last commit. Then, create a pull request.

## 3. Create and push tags

Once the previous changes have been approved and merged, pull the latest changes in the `main` branch. Then, create the tags for the last commit and push them with `make multimod-push-tags`.

## 4. Release

Finally, create a Release on GitHub.
