package httpserver

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/suapapa/croquis-king/internal/lobby"
	"github.com/suapapa/croquis-king/internal/pixabay"
)

type pixabayHandler struct {
	store  lobby.Store
	client *pixabay.Client
}

type searchImagesResponse struct {
	Total     int                 `json:"total"`
	TotalHits int                 `json:"total_hits"`
	Hits      []searchImageItem   `json:"hits"`
	RateLimit rateLimitResponse   `json:"rate_limit"`
}

type searchImageItem struct {
	PixabayID     int    `json:"pixabay_id"`
	PageURL       string `json:"page_url"`
	PreviewURL    string `json:"preview_url"`
	WebformatURL  string `json:"webformat_url"`
	LargeImageURL string `json:"large_image_url"`
	Width         int    `json:"width"`
	Height        int    `json:"height"`
	Views         int    `json:"views"`
	Downloads     int    `json:"downloads"`
	Likes         int    `json:"likes"`
}

type rateLimitResponse struct {
	Limit     int `json:"limit"`
	Remaining int `json:"remaining"`
	Reset     int `json:"reset"`
}

func newPixabayHandler(store lobby.Store, client *pixabay.Client) *pixabayHandler {
	return &pixabayHandler{
		store:  store,
		client: client,
	}
}

func (h *pixabayHandler) search(c *gin.Context) {
	query := strings.TrimSpace(c.Query("q"))
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "q is required"})
		return
	}

	order := c.DefaultQuery("order", "popular")
	if order != "popular" && order != "latest" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order must be popular or latest"})
		return
	}

	page, err := parsePositiveInt(c.DefaultQuery("page", "1"), 1)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page must be a positive integer"})
		return
	}

	perPage, err := parsePositiveInt(c.DefaultQuery("per_page", "20"), 20)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "per_page must be a positive integer"})
		return
	}

	result, err := h.client.Search(c.Request.Context(), pixabay.SearchParams{
		Query:      query,
		Order:      order,
		Page:       page,
		PerPage:    perPage,
		SafeSearch: true,
	})
	if err != nil {
		writePixabayError(c, err)
		return
	}

	c.JSON(http.StatusOK, toSearchImagesResponse(result))
}

func toSearchImagesResponse(result pixabay.SearchResult) searchImagesResponse {
	hits := make([]searchImageItem, 0, len(result.Hits))
	for _, hit := range result.Hits {
		hits = append(hits, searchImageItem{
			PixabayID:     hit.ID,
			PageURL:       hit.PageURL,
			PreviewURL:    hit.PreviewURL,
			WebformatURL:  hit.WebformatURL,
			LargeImageURL: hit.LargeImageURL,
			Width:         hit.Width,
			Height:        hit.Height,
			Views:         hit.Views,
			Downloads:     hit.Downloads,
			Likes:         hit.Likes,
		})
	}

	return searchImagesResponse{
		Total:     result.Total,
		TotalHits: result.TotalHits,
		Hits:      hits,
		RateLimit: rateLimitResponse{
			Limit:     result.RateLimit.Limit,
			Remaining: result.RateLimit.Remaining,
			Reset:     result.RateLimit.Reset,
		},
	}
}

func writePixabayError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, pixabay.ErrEmptyQuery):
		c.JSON(http.StatusBadRequest, gin.H{"error": "q is required"})
	case errors.Is(err, pixabay.ErrRateLimited):
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "pixabay rate limit exceeded"})
	default:
		var apiErr *pixabay.APIError
		if errors.As(err, &apiErr) {
			status := apiErr.StatusCode
			if status < http.StatusBadRequest {
				status = http.StatusBadGateway
			}
			c.JSON(status, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": "pixabay search failed"})
	}
}

func parsePositiveInt(raw string, defaultValue int) (int, error) {
	if strings.TrimSpace(raw) == "" {
		return defaultValue, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return 0, err
	}
	return value, nil
}

func registerPixabayRoutes(r *gin.RouterGroup, store lobby.Store, client *pixabay.Client) {
	handler := newPixabayHandler(store, client)
	admin := r.Group("")
	admin.Use(requireAdmin(store))
	admin.GET("/pixabay/search", handler.search)
}
