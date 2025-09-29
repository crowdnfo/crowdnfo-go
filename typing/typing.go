package typing

// ProcessResult holds the result of processing a release, including any non-fatal warnings.
type ProcessResult struct {
	Warnings []error
}

type ProgressCB func(stage string, releasename string, detail string)
