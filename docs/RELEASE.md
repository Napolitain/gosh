# Release Process

This document describes the release process for gosh.

## Calendar Versioning

gosh uses calendar versioning in the format: `v.YYYY.MM.DD.BUILD`

Examples:
- `v2025.10.26.0` - First release on October 26, 2025
- `v2025.10.26.1` - Second release on the same day
- `v2025.11.01.0` - First release on November 1, 2025

## Automated Releases

Releases are automatically created when code is pushed to the `main` branch (excluding documentation changes).

### Automatic Release Workflow

1. Push changes to `main` branch
2. The `auto-release` workflow runs:
   - Runs `go mod tidy`
   - Builds the project
   - Runs all tests
   - Generates a calendar version tag
   - Builds binaries for:
     - Linux (amd64)
     - macOS (amd64, arm64)
     - Windows (amd64)
   - Creates a GitHub release
   - Attaches binaries to the release

### Manual Release

You can also trigger a release manually:

1. Go to Actions tab in GitHub
2. Select "Auto Release" workflow
3. Click "Run workflow"
4. Select the `main` branch
5. Click "Run workflow"

## Manual Tag-Based Release

If you need to create a release from a specific tag:

1. Create and push a tag:
   ```bash
   git tag v2025.10.26.0
   git push origin v2025.10.26.0
   ```

2. The `release` workflow will automatically:
   - Build binaries for all platforms
   - Create a GitHub release
   - Attach binaries to the release

## CI Pipeline

Every push and pull request to `main` runs:
- `go mod tidy` (with validation)
- `go build -v ./...`
- `go test -v ./...`

## Platform Support

Releases include binaries for:
- **Linux**: `gosh-linux-amd64`
- **macOS Intel**: `gosh-darwin-amd64`
- **macOS Apple Silicon**: `gosh-darwin-arm64`
- **Windows**: `gosh-windows-amd64.exe`

## Version Information

Binaries include version information embedded during build. To check the version:

```bash
./gosh --version
```

## Download Latest Release

Visit the [Releases](https://github.com/Napolitain/gosh/releases) page to download the latest binaries.
