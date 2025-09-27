package typing

import (
	"fmt"
	"strings"
)

// UploadError represents multiple errors that occurred during upload.
// SuccessCount is the number of successful uploads.
// Errors contains all error messages.
type UploadError struct {
	SuccessCount int
	Errors       []string
}

// Implement the error interface
func (e *UploadError) Error() string {
	summary := ""
	if e.SuccessCount > 0 {
		summary = fmt.Sprintf("partial_failure: %d successful, %d failed\n", e.SuccessCount, len(e.Errors))
	} else {
		summary = fmt.Sprintf("total_failure: %d upload(s) failed\n", len(e.Errors))
	}
	details := strings.Join(e.Errors, "\n")
	return fmt.Sprintf("%sDetails:\n%s", summary, details)
}
