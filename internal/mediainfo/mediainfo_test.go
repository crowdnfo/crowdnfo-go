package mediainfo

import (
	"testing"
)

func TestParseMediaInfoVersion(t *testing.T) {
	tests := []struct {
		name            string
		versionOutput   string
		expectedVersion int
		expectError     bool
	}{
		{
			name:            "Actual MediaInfo format v25.07",
			versionOutput:   "MediaInfo Command line,\nMediaInfoLib - v25.07",
			expectedVersion: 2507,
			expectError:     false,
		},
		{
			name:            "MediaInfo format v23.00",
			versionOutput:   "MediaInfo Command line,\nMediaInfoLib - v23.00",
			expectedVersion: 2300,
			expectError:     false,
		},
		{
			name:            "MediaInfo format v24.06",
			versionOutput:   "MediaInfo Command line,\nMediaInfoLib - v24.06",
			expectedVersion: 2406,
			expectError:     false,
		},
		{
			name:            "Legacy format v23.11",
			versionOutput:   "MediaInfo v23.11",
			expectedVersion: 2311,
			expectError:     false,
		},
		{
			name:            "Version without v prefix",
			versionOutput:   "MediaInfoLib - 24.06",
			expectedVersion: 2406,
			expectError:     false,
		},
		{
			name:            "Version 20.09",
			versionOutput:   "MediaInfoLib - v20.09",
			expectedVersion: 2009,
			expectError:     false,
		},
		{
			name:            "Invalid format - no version",
			versionOutput:   "MediaInfo Command line",
			expectedVersion: 0,
			expectError:     true,
		},
		{
			name:            "Empty string",
			versionOutput:   "",
			expectedVersion: 0,
			expectError:     true,
		},
		{
			name:            "Version with extra whitespace",
			versionOutput:   "MediaInfo Command line,\n\nMediaInfoLib - v24.06",
			expectedVersion: 2406,
			expectError:     false,
		},
		{
			name:            "Minimum version 23.0",
			versionOutput:   "MediaInfoLib - v23.0",
			expectedVersion: 2300,
			expectError:     false,
		},
		{
			name:            "Single digit minor version",
			versionOutput:   "MediaInfoLib - v23.1",
			expectedVersion: 2301,
			expectError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, err := parseMediaInfoVersion(tt.versionOutput)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if version != tt.expectedVersion {
				t.Errorf("Expected version %d, got %d", tt.expectedVersion, version)
			}
		})
	}
}

func TestMinMediaInfoVersion(t *testing.T) {
	// Test that the minimum version constant is set correctly
	if MinMediaInfoVersion != 2300 {
		t.Errorf("Expected MinMediaInfoVersion to be 2300 (23.0), got %d", MinMediaInfoVersion)
	}
}
