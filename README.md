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
	"fmt"
	"log"

	"github.com/crowdnfo/crowdnfo-go" // import your library
)

func main() {
	opts := crowdnfo.Options{
		ReleasePath:     "",  // path to the release directory
		MediaInfoPath:   "",  // path to mediainfo binary (optional, defaults to "mediainfo" in PATH)
		MediaInfoJSON:   nil, // actual MediaInfo JSON data as byte[] (optional)
		Category:        "",  // e.g., "TV", "Movies" (optional, auto-detected if empty)
		NFOFilePath:     "",  // path to the NFO file (optional, auto-detected if empty)
		APIKey:          "",  // your CrowdNFO API key
		MaxHashFileSize: 0,   // max file size for hashing in bytes (0 for no limit, -1 for do not hash)
		ArchiveDir:      "",  // directory to archive uploaded metadata, empty for no archiving
		ProgressCB: func(stage, detail string) {
			fmt.Printf("[%s]\t%s\n", stage, detail)
		},
	}

	result, err := crowdnfo.ProcessRelease(opts)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	for _, warn := range result.Warnings {
		log.Printf("Warning: %v", warn)
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
