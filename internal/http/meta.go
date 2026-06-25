package httpserver

import (
	"html"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	defaultAppLogo       = "/example_logo.png"
	defaultAppTitle      = "Croquing"
	defaultOGDescription = "Real-time croquis meetups with synchronized photos and timers."
)

type socialMeta struct {
	Title       string
	Description string
	Image       string
	URL         string
}

func resolveAppTitle(appName string) string {
	if trimmed := strings.TrimSpace(appName); trimmed != "" {
		return trimmed
	}
	return defaultAppTitle
}

func resolveAppLogo(appLogo string) string {
	if trimmed := strings.TrimSpace(appLogo); trimmed != "" {
		return trimmed
	}
	return defaultAppLogo
}

func requestAbsoluteURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	} else if proto := c.GetHeader("X-Forwarded-Proto"); proto != "" {
		scheme = strings.TrimSpace(strings.Split(proto, ",")[0])
	} else if c.Request.URL != nil && c.Request.URL.Scheme != "" {
		scheme = c.Request.URL.Scheme
	}
	host := c.Request.Host
	if host == "" {
		host = "localhost"
	}

	pageURL := url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   c.Request.URL.Path,
	}
	return pageURL.String()
}

func absoluteAssetURL(c *gin.Context, assetPath string) string {
	if parsed, err := url.Parse(assetPath); err == nil && parsed.IsAbs() {
		return assetPath
	}

	base, err := url.Parse(requestAbsoluteURL(c))
	if err != nil {
		return assetPath
	}

	path := assetPath
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	base.Path = path
	base.RawQuery = ""
	base.Fragment = ""
	return base.String()
}

func buildSocialMetaTags(c *gin.Context, appName, appLogo string) string {
	meta := socialMeta{
		Title:       resolveAppTitle(appName),
		Description: defaultOGDescription,
		Image:       absoluteAssetURL(c, resolveAppLogo(appLogo)),
		URL:         requestAbsoluteURL(c),
	}
	return renderSocialMetaTags(meta)
}

func renderSocialMetaTags(meta socialMeta) string {
	title := html.EscapeString(meta.Title)
	description := html.EscapeString(meta.Description)
	image := html.EscapeString(meta.Image)
	pageURL := html.EscapeString(meta.URL)

	return "" +
		`<meta property="og:type" content="website" />` + "\n    " +
		`<meta property="og:title" content="` + title + `" />` + "\n    " +
		`<meta property="og:description" content="` + description + `" />` + "\n    " +
		`<meta property="og:image" content="` + image + `" />` + "\n    " +
		`<meta property="og:url" content="` + pageURL + `" />` + "\n    " +
		`<meta name="twitter:card" content="summary_large_image" />` + "\n    " +
		`<meta name="twitter:title" content="` + title + `" />` + "\n    " +
		`<meta name="twitter:description" content="` + description + `" />` + "\n    " +
		`<meta name="twitter:image" content="` + image + `" />` + "\n    "
}

func injectBeforeHeadClose(htmlContent, injection string) string {
	const needle = "</head>"
	if idx := strings.Index(htmlContent, needle); idx >= 0 {
		return htmlContent[:idx] + injection + htmlContent[idx:]
	}
	return htmlContent
}

func withAppTitle(htmlContent, appName string) string {
	title := html.EscapeString(resolveAppTitle(appName))

	const openTag = "<title>"
	const closeTag = "</title>"

	start := strings.Index(htmlContent, openTag)
	if start < 0 {
		return htmlContent
	}
	contentStart := start + len(openTag)
	end := strings.Index(htmlContent[contentStart:], closeTag)
	if end < 0 {
		return htmlContent
	}

	return htmlContent[:contentStart] + title + htmlContent[contentStart+end:]
}
