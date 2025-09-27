package files

// Structures for season pack processing
type VideoFile struct {
	Path string
	Dir  string
	Name string
}

// FileListEntry represents a single file in a file list
type FileListEntry struct {
	FilePath      string `json:"filePath"`
	FileSizeBytes int64  `json:"fileSizeBytes"`
}

// FileListRequest represents the JSON structure for file list API requests
type FileListRequest struct {
	ReleaseName string          `json:"releaseName"`
	Category    string          `json:"category"`
	Entries     []FileListEntry `json:"entries"`
}

type EpisodeInfo struct {
	VideoFile   VideoFile
	EpisodeNum  string // "E01", "E19" etc.
	ReleaseName string
	NFOFile     string
}
