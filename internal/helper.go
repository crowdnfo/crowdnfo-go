package internal

import (
	"regexp"

	"github.com/crowdnfo/crowdnfo-go/internal/files"
	"github.com/crowdnfo/crowdnfo-go/typing"
)

// isSeasonPack determines if the given job name corresponds to a season pack
func IsSeasonPack(jobName string) bool {
	// First check for traditional season patterns (S\d{2,4}) but NOT episodes (SxxExx)
	seasonPattern := regexp.MustCompile(`(?i)\bS\d{2,4}\b`)
	episodePattern := regexp.MustCompile(`(?i)\bS\d{2,4}E\d{2,4}\b`)

	if seasonPattern.MatchString(jobName) && !episodePattern.MatchString(jobName) {
		return true
	}

	// Check for ISO date format years (S2024, S2023, etc.)
	isoYearPattern := regexp.MustCompile(`(?i)\bS(20\d{2})\b`)
	if isoYearPattern.MatchString(jobName) {
		return true
	}

	return false
}

// isSeasonPackFallback checks if a directory should be treated as season pack based on video file count
func IsSeasonPackFallback(finalDir string) bool {
	videoFiles, err := files.FindAllVideoFiles(finalDir)
	if err != nil {
		return false
	}

	// If we find 3 or more video files, treat as season pack
	return len(videoFiles) >= 3
}

func MergeProcessResults(a, b *typing.ProcessResult) *typing.ProcessResult {
	if a == nil && b == nil {
		return &typing.ProcessResult{}
	}
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	return &typing.ProcessResult{
		Warnings: append(a.Warnings, b.Warnings...),
	}
}
