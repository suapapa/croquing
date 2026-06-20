package pixabay

import "errors"

var (
	// ErrEmptyQuery is returned when a search is attempted without a query.
	ErrEmptyQuery = errors.New("pixabay: query is required")

	// ErrRateLimited indicates the Pixabay API rate limit was exceeded.
	ErrRateLimited = errors.New("pixabay: rate limit exceeded")
)

// APIError represents a non-success HTTP response from Pixabay.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return "pixabay: request failed"
}
