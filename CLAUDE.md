# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`go-org-readwise` is a CLI tool that syncs highlights from the Readwise API to org-mode files using the denote file naming convention. It implements only the export endpoint of the Readwise v2 API, fetching all highlights and creating/updating org files in a target directory.

## Building and Running

This project uses Nix flakes for development and builds:

```bash
# Build the project
nix build

# Run the project
nix run . -- -targetFolder /path/to/org/files -apiKeyFile /path/to/key

# Or with environment variable
READWISE_KEY=your_key nix run . -- -targetFolder /path/to/org/files

# Enter development shell
nix develop

# Inside dev shell, use standard Go commands
go build
go test -v ./...
```

Standard Go commands (when not using Nix):
```bash
# Build
go build -v ./...

# Run tests
go test -v ./...

# Run a single test
go test -v ./internal/org -run TestFunctionName
```

## Architecture

The codebase follows a clean 3-layer architecture:

### 1. Main Entry Point (`main.go`)
- Handles CLI flags:
  - `-targetFolder` (required): Destination folder for org files
  - `-apiKeyFile` (optional, falls back to `READWISE_KEY` env var): API key file
  - `-archiveURLs` (optional, default false): Archive document URLs using monolith
- Manages state file (`.readwise-sync.state`) in the target folder to track incremental syncs
- Coordinates the sync flow: fetch from Readwise API → write to org files

### 2. Readwise Package (`internal/readwise/`)
- **readwise.go**: Implements `FetchFromAPI()` which handles paginated API calls to the Readwise export endpoint
- **types.go**: Defines API response types (`Export`, `Result`, `Highlight`, `Tag`)
- Key constant: `FormatUpdatedAfter = "2006-01-02T15:04:05"` for incremental sync timestamps

### 3. Org Package (`internal/org/`)
- **sync.go**: Core sync logic that creates new org files or appends to existing ones
  - Implements denote file naming convention: `DATE==SIGNATURE--TITLE__KEYWORDS.EXTENSION`
  - Example: `20240511T100401==readwise=books--economics-in-the-euro-area__finance.org`
  - Signature format: `==readwise={category}` where category comes from Readwise
  - URL archiving: When `-archiveURLs` is enabled, archives URLs using monolith to `.archive` folder
- **org.go**: Type definitions for `Document` (full new file) and `PartialDocument` (appending highlights)
  - `Document.ArchivePath`: Relative path to archived HTML file (when archiving is enabled)
- **write.go**: Template rendering using Go's `text/template` to generate org-mode formatted content
  - Includes `#+property: ARCHIVE:` header when ArchivePath is set

### Key Design Decisions

1. **Denote Integration**: Files use denote naming convention with `==readwise={category}` as the signature field to mark synced content
2. **Incremental Sync**: Uses `.readwise-sync.state` file to store last sync timestamp and fetch only updated highlights
3. **Dual Document Types**:
   - `Document` for creating new org files with full metadata
   - `PartialDocument` for appending new highlights to existing files
4. **Tag Transformation**: Both book-level tags (in filename/filetags) and highlight-level tags are sluggified and converted to org-mode tag format (`:tag1:tag2:`)
5. **URL Archiving**: Optional feature using monolith to archive source URLs
   - Archives stored in `.archive/` folder with same naming as org files (`.html` extension)
   - Fails gracefully with warning if monolith not available or URL unreachable
   - Archive path added as `#+property: ARCHIVE:` header in org file

## Date/Time Formats

- Denote ID format: `20060102T150405` (compact date-time)
- Org date format: `2006-01-02 Mon 15:04` (human-readable with day name)
- Readwise API format: `2006-01-02T15:04:05` (ISO-8601 style)

## CI/CD

GitHub Actions workflows:
- **ci.yml**: Runs tests (`go test -v ./...`) and builds (`go build -v ./...`), plus golangci-lint
- **nix.yml**: Nix-based builds
- **release.yml**: Handles releases

All GitHub Actions are pinned to commit SHAs for security.
