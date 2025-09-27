package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Lokilicious/crowdnfo-go" // import your library
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
		ArchiveDir:      "",  // directory to archive uploaded metadata (optional)
	}

	err := crowdnfo.ProcessRelease(opts, log.New(os.Stdout, "crowdnfo: ", log.LstdFlags))

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
