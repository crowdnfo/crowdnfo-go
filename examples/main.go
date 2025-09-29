package main

import (
	"fmt"
	"log"

	"github.com/crowdnfo/crowdnfo-go" // import your library
)

func main() {
	opts := crowdnfo.Options{
		ReleasePath:     "", // path to the release directory
		MediaInfoPath:   "", // path to mediainfo binary (optional, defaults to "mediainfo" in PATH)
		Category:        "", // e.g., "TV", "Movies" (optional, auto-detected if empty)
		NFOFilePath:     "", // path to the NFO file (optional, auto-detected if empty)
		APIKey:          "", // your CrowdNFO API key
		MaxHashFileSize: 0,  // max file size for hashing in bytes (0 for no limit, -1 for do not hash)
		ArchiveDir:      "", // directory to archive uploaded metadata, empty for no archiving
		ProgressCB: func(stage, releaseName, detail string) {
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
