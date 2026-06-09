# gen-apidocs

`gen-apidocs` generates the Kubernetes API reference documentation from an OpenAPI `swagger.json`. It produces two outputs:

- **HTML backend** — the single-page API reference published at `https://kubernetes.io/docs/reference/generated/kubernetes-api/v<X.Y>/`.
- **Markdown backend** — Hugo-native pages under `content/en/docs/reference/kubernetes-api/` in `kubernetes/website`.

Both backends are supported. Pick the one that matches the destination in `kubernetes/website`.

For the canonical end-to-end release walkthrough, see [Generating Reference Documentation for the Kubernetes API](https://kubernetes.io/docs/contribute/generate-ref-docs/kubernetes-api/).

## Inputs

| Input | Purpose |
|---|---|
| `K8S_RELEASE` | Kubernetes release version, for example `1.36.0`. The generator uses the `X.Y` prefix to locate per-release config. |
| `K8S_ROOT` | Path to a local `kubernetes/kubernetes` checkout at the matching release tag. Required for the swagger refresh step. |
| `K8S_WEBROOT` | Path to a local `kubernetes/website` checkout. Required for the copy targets. |
| `gen-apidocs/config/v<X_Y>/swagger.json` | OpenAPI spec for the release. |
| `gen-apidocs/config/v<X_Y>/toc.yaml` | Section and category layout used by the renderers. |
| `gen-apidocs/config/v<X_Y>/config.yaml` | Resource grouping rules and supplementary metadata. |

## Prepare the swagger input

The `swagger.json` checked into `kubernetes/kubernetes` at `api/openapi-spec/swagger.json` is missing many enum fields that the API reference needs. Regenerate the file with `OpenAPIEnums=true` before running the generator.

### Option A: from a temporary source checkout (no local k/k needed)

From the reference-docs repository root, with only `K8S_RELEASE` set:

```shell
make updateapispec-enums-from-source
```

This shallow-clones the release tag into a temporary directory, enables
`OpenAPIEnums=true` in that copy only, runs k/k's `hack/update-openapi-spec.sh`,
copies the result into `gen-apidocs/config/v<X_Y>/swagger.json`, verifies the
enum metadata, and deletes the temporary checkout. It needs network access and
k/k's host OpenAPI prerequisites, and leaves k/k's checked-in swagger unchanged.
Set `KEEP_TMP=1` to preserve the checkout and generation log for debugging.

### Option B: from a prepared `K8S_ROOT` checkout

From your `K8S_ROOT` checkout, at the release tag:

1. Edit `hack/update-openapi-spec.sh` and set `OpenAPIEnums=true`.
2. Run `hack/update-openapi-spec.sh` to regenerate `api/openapi-spec/swagger.json`.

Then, from the reference-docs repository root:

```shell
make updateapispec
```

This copies the regenerated swagger into `gen-apidocs/config/v<X_Y>/swagger.json`.

## Build

Run from the repository root:

```shell
make api      # HTML backend     -> gen-apidocs/build/html/
make apimd    # Markdown backend -> gen-apidocs/build/markdown/
```

Copy generated output into a `kubernetes/website` checkout:

```shell
make copyapi      # publish HTML reference
make copyapimd    # publish Markdown reference
```

## Output

- `gen-apidocs/build/html/index.html` and `navData.js` — single-page HTML reference.
- `gen-apidocs/build/markdown/` — Hugo-compatible Markdown tree organized by API group.
