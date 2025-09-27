package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Lokilicious/crowdnfo-go/internal/files"
)

var BASE_URL = "https://crowdnfo.net/api/releases"

const (
	MediaInfoType = "MediaInfo"
	NFOType       = "NFO"
	FileListType  = "FileList"
)

// UploadToCrowdNFO uploads release data to CrowdNFO.
// On failure, returns an error. If multiple errors occurred, returns an *UploadError
// which contains all error messages and the count of successful uploads.
func UploadToCrowdNFO(apiKey string, releaseName, category, hash, releasePath string, mediaInfoJSON []byte, nfoFile, archiveDir string, logger *log.Logger) error {
	fileListEntries, err := files.CreateFileList(releasePath, releaseName)
	if err != nil {
		logger.Printf("%s - Failed to create File List: %v", releaseName, err)
		fileListEntries = nil
	}
	uploadAssets(apiKey, releaseName, category, hash, archiveDir, mediaInfoJSON, nfoFile, fileListEntries, logger)

	return nil
}

// UploadEpisodeToCrowdNFO uploads release data to CrowdNFO.
// On failure, returns an error. If multiple errors occurred, returns an *UploadError
// which contains all error messages and the count of successful uploads.
func UploadEpisodeToCrowdNFO(apiKey string, episodeInfo files.EpisodeInfo, category, hash string, mediaInfoJSON []byte, archiveDir string, logger *log.Logger) {
	fileListEntries, err := files.CreateEpisodeFileList(episodeInfo)
	if err != nil {
		logger.Printf("%s - Failed to create File List: %v", episodeInfo, err)
		fileListEntries = nil
	}
	uploadAssets(apiKey, episodeInfo.ReleaseName, category, hash, archiveDir, mediaInfoJSON, episodeInfo.NFOFile, fileListEntries, logger)
}

func uploadAssets(apiKey, releaseName, category, hash, archiveDir string, mediaInfoJSON []byte, nfoFile string, fileListEntries []files.FileListEntry, logger *log.Logger) {
	// MediaInfo
	if len(mediaInfoJSON) > 0 {
		if err := uploadFile(apiKey, releaseName, MediaInfoType, "", mediaInfoJSON, hash, category, archiveDir); err != nil {
			logger.Printf("%s - %s: %v", releaseName, MediaInfoType, err)
		}
	}
	// NFO
	if nfoFile != "" {
		nfoData, err := os.ReadFile(nfoFile)
		if err != nil {
			logger.Printf("%s - %s: %v", releaseName, NFOType, err)
		} else {
			nfoFileName := filepath.Base(nfoFile)
			if err := uploadFile(apiKey, releaseName, NFOType, nfoFileName, nfoData, hash, category, archiveDir); err != nil {
				logger.Printf("%s - %s: %v", releaseName, NFOType, err)
			}
		}
	}
	// FileList
	if len(fileListEntries) > 0 {
		fileListRequest := files.FileListRequest{
			ReleaseName: releaseName,
			Category:    category,
			Entries:     fileListEntries,
		}
		if err := uploadFileList(apiKey, fileListRequest); err != nil {
			logger.Printf("%s - %s: %v", releaseName, FileListType, err)
		}
	}
}

func uploadFile(apiKey string, releaseName, fileType, originalFileName string, fileData []byte, hash, category, archiveDir string) error {
	url := fmt.Sprintf("%s/%s/files", BASE_URL, releaseName)

	// Create multipart form
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Add form fields
	writer.WriteField("FileType", fileType)
	if originalFileName != "" {
		writer.WriteField("OriginalFileName", originalFileName)
	}
	if category != "" {
		writer.WriteField("Category", category)
	}
	if hash != "" {
		writer.WriteField("FileHash", hash)
	}

	// Add file
	part, err := writer.CreateFormFile("File", getFileName(fileType, releaseName, originalFileName))
	if err != nil {
		return fmt.Errorf("create form file: %w", err)
	}
	part.Write(fileData)

	writer.Close()

	// Create request
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return fmt.Errorf("create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-Api-Key", apiKey)
	req.Header.Set("User-Agent", getUserAgent())

	// Send request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check for update headers
	//checkUpdateHeaders(resp.Header)

	// Read response body for error details
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Archive the uploaded file
	archiveFile := filepath.Join(archiveDir, getFileName(fileType, releaseName, originalFileName))
	if err := os.WriteFile(archiveFile, fileData, 0644); err != nil {
		return fmt.Errorf("failed to archive uploaded %s file: %w", fileType, err)
	}

	return nil
}

func getFileName(fileType, releaseName, originalFileName string) string {
	if fileType == "NFO" && originalFileName != "" {
		return originalFileName
	}
	return fmt.Sprintf("%s.json", releaseName)
}

// uploadFileList uploads a file list to CrowdNFO
func uploadFileList(apiKey string, fileListRequest files.FileListRequest) error {
	url := fmt.Sprintf("%s/%s/filelists", BASE_URL, fileListRequest.ReleaseName)

	// Convert to JSON
	jsonData, err := json.Marshal(fileListRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal file list: %w", err)
	}

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", apiKey)
	req.Header.Set("User-Agent", getUserAgent())

	// Send request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == http.StatusUnauthorized {
			return fmt.Errorf("unauthorized: please check your API key in config.json")
		}
		if resp.StatusCode == http.StatusBadRequest {
			return fmt.Errorf("%s", string(body))
		}
		return fmt.Errorf("file list upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// getUserAgent returns the User-Agent string for API requests
func getUserAgent() string {
	return fmt.Sprintf("crowdnfo-go/%s")
}
