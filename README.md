# crowdnfo-go

A Go library for uploading release information to [CrowdNFO](https://crowdnfo.net/).
This package aims to provide convenient and idiomatic Go APIs for automating CrowdNFO submissions, including support for video, audio, and season pack releases.

---

## Features

- Upload release metadata and NFO files to CrowdNFO
- Automatic category detection (Movies, TV, Games, etc.)
- MediaInfo integration for video/audio files
- Season pack and multi-episode support
- File hashing and validation
- Simple, idiomatic Go API

---

## Installation

```sh
go get github.com/crowdnfo/crowdnfo-go
```

---

## Usage

```go
package main

import (
    "log"
    "github.com/crowdnfo/crowdnfo-go"
)

func main() {
    opts := crowdnfo.Options{
        ReleasePath:     "/path/to/release",
        MediaInfoPath:   "", // optional, defaults to "mediainfo" in PATH
        Category:        "", // optional, auto-detect if empty
        APIKey:          "your-crowdnfo-api-key",
        ArchiveDir:      "", // optional
        MaxHashFileSize: 0,  // 0 = always hash, <0 = never hash
    }

    logger := log.New(os.Stdout, "[crowdnfo] ", log.LstdFlags)

    if err := crowdnfo.ProcessRelease(opts, logger); err != nil {
        logger.Fatalf("Failed to process release: %v", err)
    }
}
```

---

## Requirements

- Go 1.21 or newer (for `slices.Contains`)
- [mediainfo](https://mediaarea.net/en/MediaInfo) binary in your PATH (optional, for media info extraction)
- A valid CrowdNFO API key

---

## Acknowledgments

Special thanks to [pixelhunterX](https://github.com/pixelhunterX/) for laying the base of this project and inspiring its development.

---

## License

MIT License. See [LICENSE](./LICENSE) for details.
