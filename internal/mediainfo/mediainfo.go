package mediainfo

import (
	"fmt"
	"os/exec"
)

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
