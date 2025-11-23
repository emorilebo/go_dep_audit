# Publishing Guide

## Prerequisites

- Access to the GitHub repository.
- Git installed.

## Steps

1. **Update Version**: Ensure the code is ready and version numbers are updated.

2. **Tag the Release**:
   ```bash
   git tag v0.1.0
   git push origin v0.1.0
   ```

3. **Publish to Go Proxy**:
   The Go proxy will automatically pick up the new tag. You can force it by visiting:
   `https://proxy.golang.org/github.com/emorilebo/go_dep_audit/@v/v0.1.0.info`

4. **Create GitHub Release**:
   Go to the GitHub repository -> Releases -> Draft a new release.
   Select the tag `v0.1.0`.
   Generate release notes.

## CI/CD

The `.github/workflows/ci.yml` file handles testing and linting on every push.
