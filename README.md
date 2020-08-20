# Action to Publish and Verify Release Assets into Asset Transparency Log

This action is designed to trigger on GitHub Release events and adds the release assets (including GitHub generated source tarballs) to the [Asset Transparency Log](https://www.transparencylog.com).

## Inputs

**None**

## Outputs

### `verified`

The list of verified URLs

### `failed`

The list of URLs that failed to match the asset logs digest

## Example Workflow

[See example workflow](https://github.com/transparencylog/github-releases-asset-transparency-verify-action/blob/main/.github/workflows/asset-transparency.yaml)

### API Docs Used

- https://docs.github.com/en/actions/reference/context-and-expression-syntax-for-github-actions#github-context
- https://docs.github.com/en/actions/configuring-and-managing-workflows/using-environment-variables
- https://docs.github.com/en/actions/reference/events-that-trigger-workflows#release
- https://pkg.go.dev/github.com/google/go-github/v32/github?tab=doc#ReleaseEvent
