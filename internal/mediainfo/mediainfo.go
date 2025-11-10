package mediainfo

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const MinMediaInfoVersion = 2300

func GenerateMediaInfoJSON(filePath, mediaInfoPath string) ([]byte, error) {
	if mediaInfoPath == "" {
		return nil, fmt.Errorf("MediaInfo is not available")
	}

	cmd := exec.Command(mediaInfoPath, "--Output=JSON", filePath)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run MediaInfo: %v", err)
	}

	return output, nil
}

// CheckMediaInfoVersion checks if MediaInfo version is >= MinMediaInfoVersion
func CheckMediaInfoVersion(mediaInfoPath string) error {
	if mediaInfoPath == "" {
		return fmt.Errorf("MediaInfo path is empty")
	}

	cmd := exec.Command(mediaInfoPath, "--Version")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get MediaInfo version: %v", err)
	}

	version, err := parseMediaInfoVersion(string(output))
	if err != nil {
		return fmt.Errorf("failed to parse MediaInfo version: %v", err)
	}

	if version < MinMediaInfoVersion {
		return fmt.Errorf("MediaInfo version %d is too old (minimum required: %d)", version, MinMediaInfoVersion)
	}

	return nil
}

// parseMediaInfoVersion extracts version number from MediaInfo --Version output
// Expected format:
// "MediaInfo Command line,
//
//	MediaInfoLib - v25.07"
func parseMediaInfoVersion(versionOutput string) (int, error) {
	// Remove any leading/trailing whitespace
	versionOutput = strings.TrimSpace(versionOutput)

	// Pattern to match version like "MediaInfoLib - v24.06", "v23.11", etc.
	// This handles both the actual format and potential variations
	versionPattern := regexp.MustCompile(`(?:MediaInfoLib\s*-\s*)?v?(\d+)\.(\d+)`)
	matches := versionPattern.FindStringSubmatch(versionOutput)

	if len(matches) < 3 {
		return 0, fmt.Errorf("could not parse version from: %s", versionOutput)
	}

	major, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, fmt.Errorf("invalid major version: %s", matches[1])
	}

	minor, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0, fmt.Errorf("invalid minor version: %s", matches[2])
	}

	// Convert to format like 2507, 2300 for easy comparison
	version := major*100 + minor

	return version, nil
}
