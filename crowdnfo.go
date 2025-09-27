package crowdnfo

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"

	"github.com/crowdnfo/crowdnfo-go/internal"
	"github.com/crowdnfo/crowdnfo-go/internal/api"
	"github.com/crowdnfo/crowdnfo-go/internal/files"
	"github.com/crowdnfo/crowdnfo-go/internal/mediainfo"
)

// Options holds all parameters for processing a release.
type Options struct {
	ReleasePath     string
	MediaInfoPath   string // optional, defaults to "mediainfo" in PATH
	MediaInfoJSON   []byte // optional, if provided, skip mediainfo execution
	Category        string
	NFOFilePath     string // optional
	APIKey          string
	ArchiveDir      string
	MaxHashFileSize int64
}

// Valid CrowdNFO categories
var validCategories = []string{"Movies", "TV", "Games", "Software", "Music", "Audiobooks", "Books", "Other"}

// Built-in regex patterns for category detection
var categoryRegexPatterns = []struct {
	Pattern  string
	Category string
}{
	{`(?i)\b(audiobook|abook|abookde|hörbuch|hoerbuch|horbuch|m4b)\b`, "Audiobooks"},
	{`(?i)\b(ebook|epaper|pdf|epub|mobi)\b`, "Books"},
	{`(?i)\b((s\d{1,4}e\d{1,4})|(s\d{1,4})|(e\d{1,4})|season|staffel|episode|folge|(\d{4}-\d{2}-\d{2}))\b`, "TV"},
	{`(?i)\b(elamigos|gog|xbox|xbox360|x360|ps\d|nintendo|nsw|amiga|atari|wii[u]?)\b`, "Games"},
	{`(?i)\b(patch|crack|cracked|keygen|keymaker|keyfilemaker|x64|dvt|btcr|macos)\b`, "Software"},
	{`(?i)\b((\d{3,4}[pi])|bluray|dvdrip|webrip|hdtv|bdrip|dvd|remux|mpeg[-]?2|vc[-]?1|avc|hevc|([xh][. ]?26[456]))\b`, "Movies"},
	{`(?i)\b(mp3|flac|webflac|aac|wav|album|artist|discography|single|vinyl|cd|\d+bit|\d+khz)\b`, "Music"},
}

// ProcessRelease is the main entrypoint for uploading release info to CrowdNFO.
func ProcessRelease(opts Options, logger *log.Logger) error {

	mediaInfoPath := checkMediaInfoAvailable(opts.MediaInfoPath)
	mediaInfoJSON := opts.MediaInfoJSON

	releaseName := files.GetBaseOrName(opts.ReleasePath)
	if releaseName == "" {
		return fmt.Errorf("Could not determine release name from path: %s", opts.ReleasePath)
	}

	if opts.APIKey == "" {
		return fmt.Errorf("API key is required")
	}

	category := getCategory(opts.Category, releaseName)
	if category == "" {
		return fmt.Errorf("Invalid category: %s", category)
	}

	// Check if this is a season pack
	if internal.IsSeasonPack(releaseName) || internal.IsSeasonPackFallback(opts.ReleasePath) {
		logger.Printf("Detected Season Pack")
		if err := processSeasonPack(opts.APIKey, opts.ReleasePath, releaseName, category, opts.ArchiveDir, mediaInfoPath, opts.MaxHashFileSize, logger); err != nil {
			return err
		}
		return nil
	}

	mediaFile, err := files.FindBiggestFile(opts.ReleasePath)
	if err != nil || mediaFile == "" {
		mediaFile, err = files.FindFirstAudioFile(opts.ReleasePath)
		if err != nil || mediaFile == "" {
			return fmt.Errorf("No media file found in: %s", opts.ReleasePath)
		}
	}

	// Generate MediaInfo and hash if media file found
	if mediaFile != "" && mediaInfoPath != "" && len(mediaInfoJSON) == 0 {
		// Generate MediaInfo JSON only for non-hash-only files
		if !files.IsHashOnlyFile(mediaFile) {
			mediaInfoJSON, err = mediainfo.GenerateMediaInfoJSON(mediaFile, mediaInfoPath)
			if err != nil {
				logger.Printf("%s - Failed to generate MediaInfo: %s", releaseName, err)
			}
		}
	}

	var hash string
	// Calculate hash for any file found (media or ISO/IMG)
	if mediaFile != "" {
		shouldHash, err := shouldCalculateHash(mediaFile, opts.MaxHashFileSize)
		if err != nil {
			return err
		} else if shouldHash {
			hash, err = calculateSHA256(mediaFile)
			if err != nil {
				return err
			}
		} else {
			logger.Printf("Skip Hashing: File exceeds max_hash_file_size limit")
		}
	}

	nfoFile, err := files.FindNFOFile(opts.ReleasePath)
	if err != nil {
		logger.Printf("No NFO File found")
		nfoFile = "" // Set empty string for upload function
	}

	err = api.UploadToCrowdNFO(opts.APIKey, releaseName, category, hash, opts.ReleasePath, mediaInfoJSON, nfoFile, opts.ArchiveDir, logger)

	if err != nil {
		return fmt.Errorf("upload to CrowdNFO failed: %w", err)
	}

	return nil
}

// processSeasonPack handles the processing of season packs
func processSeasonPack(apiKey string, releasePath string, releaseName string, category string, archiveDir string, mediaInfoPath string, maxHashFileSize int64, logger *log.Logger) error {

	// Find all video files in the season pack
	videoFiles, err := files.FindAllVideoFiles(releasePath)
	if err != nil {
		return err
	}

	if len(videoFiles) == 0 {
		logger.Printf("No Video Files found in %s", releasePath)
		return nil
	}

	// Extract episode information for each video file
	episodes := make([]files.EpisodeInfo, 0)
	generalNFO := files.FindGeneralNFO(releasePath)

	for _, videoFile := range videoFiles {
		episodeInfo := files.ExtractEpisodeInfo(videoFile, releaseName, generalNFO)
		if episodeInfo.ReleaseName != "" { // Only process valid episodes
			episodes = append(episodes, episodeInfo)
		}
	}

	if len(episodes) == 0 {
		logger.Printf("No Episodes found in %s", releasePath)
		return nil
	}

	for _, episode := range episodes {
		// Calculate SHA256 for this episode (check file size limit first)
		var hash string
		shouldHash, err := shouldCalculateHash(releasePath, maxHashFileSize)
		if err != nil {
			logger.Printf("%s - %s", episode.ReleaseName, err)
		} else if shouldHash {
			hash, err = calculateSHA256(episode.VideoFile.Path)
			if err != nil {
				logger.Printf("%s - Failed to generate Hash: %s", episode.ReleaseName, err)
				continue
			}
		}

		// Generate MediaInfo JSON for this episode
		var mediaInfoJSON []byte
		if mediaInfoPath != "" {
			mediaInfoJSON, err = mediainfo.GenerateMediaInfoJSON(episode.VideoFile.Path, mediaInfoPath)
			if err != nil {
				logger.Printf("%s - Failed to generate MediaInfo: %s", episode.ReleaseName, err)
			}
		}
		// Upload this episode to CrowdNFO API with file list
		api.UploadEpisodeToCrowdNFO(apiKey, episode, category, hash, mediaInfoJSON, archiveDir, logger)

	}

	return nil
}

// matchCategoryByRegex tries to determine category from release name using built-in regex patterns
func matchCategoryByRegex(releaseName string) string {
	for _, regexRule := range categoryRegexPatterns {
		regex, err := regexp.Compile(regexRule.Pattern)
		if err != nil {
			continue
		}
		if regex.MatchString(releaseName) {
			return regexRule.Category
		}
	}
	return ""
}

// IsValidCategory checks if the provided category is in validCategories
func isValidCategory(category string) bool {
	for _, valid := range validCategories {
		if category == valid {
			return true
		}
	}
	return false
}

// shouldCalculateHash checks if hash should be calculated based on file size and config
func shouldCalculateHash(filePath string, maxHashFileSize int64) (bool, error) {
	// If max_hash_file_size is "0", always calculate hash
	if maxHashFileSize == 0 {
		return true, nil
	}

	// If max_hash_file_size is "<0", never calculate hash
	if maxHashFileSize < 0 {
		return false, nil
	}

	// Get file size
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false, err
	}

	if fileInfo.Size() > maxHashFileSize {
		return false, nil
	}

	return true, nil
}

func calculateSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func checkMediaInfoAvailable(path string) string {
	mediaInfoPath := path
	if mediaInfoPath == "" {
		mediaInfoPath = "mediainfo" // fallback to PATH
	}
	// Check if binary exists before running
	if _, err := exec.LookPath(mediaInfoPath); err != nil {
		mediaInfoPath = ""
	}

	return mediaInfoPath
}

func getCategory(category string, releaseName string) string {
	if category == "" {
		category = matchCategoryByRegex(releaseName)
	} else if !isValidCategory(category) {
		return ""
	}

	return category
}
