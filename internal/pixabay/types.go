package pixabay

const (
	defaultBaseURL    = "https://pixabay.com/api/"
	defaultImageType  = "photo"
	defaultOrder      = "popular"
	defaultPage       = 1
	defaultPerPage    = 20
	minPerPage        = 3
	maxPerPage        = 200
)

// SearchParams holds query parameters for a Pixabay image search.
type SearchParams struct {
	Query      string
	Order      string
	Page       int
	PerPage    int
	ImageType  string
	SafeSearch bool
}

// SearchResult is the normalized search response.
type SearchResult struct {
	Total     int
	TotalHits int
	Hits      []Image
	RateLimit RateLimit
}

// Image contains the fields used by Croquing from a Pixabay hit.
type Image struct {
	ID            int
	PageURL       string
	PreviewURL    string
	WebformatURL  string
	LargeImageURL string
	Width         int
	Height        int
	Views         int
	Downloads     int
	Likes         int
}

// RateLimit captures Pixabay rate-limit response headers.
type RateLimit struct {
	Limit     int
	Remaining int
	Reset     int
}

type apiSearchResponse struct {
	Total     int         `json:"total"`
	TotalHits int         `json:"totalHits"`
	Hits      []apiHit    `json:"hits"`
}

type apiHit struct {
	ID            int    `json:"id"`
	PageURL       string `json:"pageURL"`
	PreviewURL    string `json:"previewURL"`
	WebformatURL  string `json:"webformatURL"`
	LargeImageURL string `json:"largeImageURL"`
	ImageWidth    int    `json:"imageWidth"`
	ImageHeight   int    `json:"imageHeight"`
	Views         int    `json:"views"`
	Downloads     int    `json:"downloads"`
	Likes         int    `json:"likes"`
}

func normalizeSearchParams(params SearchParams) SearchParams {
	if params.Order == "" {
		params.Order = defaultOrder
	}
	if params.Page <= 0 {
		params.Page = defaultPage
	}
	if params.PerPage <= 0 {
		params.PerPage = defaultPerPage
	}
	if params.PerPage < minPerPage {
		params.PerPage = minPerPage
	}
	if params.PerPage > maxPerPage {
		params.PerPage = maxPerPage
	}
	if params.ImageType == "" {
		params.ImageType = defaultImageType
	}
	return params
}

func toSearchResult(resp apiSearchResponse, rateLimit RateLimit) SearchResult {
	hits := make([]Image, 0, len(resp.Hits))
	for _, hit := range resp.Hits {
		hits = append(hits, Image{
			ID:            hit.ID,
			PageURL:       hit.PageURL,
			PreviewURL:    hit.PreviewURL,
			WebformatURL:  hit.WebformatURL,
			LargeImageURL: hit.LargeImageURL,
			Width:         hit.ImageWidth,
			Height:        hit.ImageHeight,
			Views:         hit.Views,
			Downloads:     hit.Downloads,
			Likes:         hit.Likes,
		})
	}

	return SearchResult{
		Total:     resp.Total,
		TotalHits: resp.TotalHits,
		Hits:      hits,
		RateLimit: rateLimit,
	}
}
