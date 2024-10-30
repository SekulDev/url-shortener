package pkg

import "net/url"

func IsValidUrl(urlParam string) bool {
	parsedURL, err := url.Parse(urlParam)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return false
	}

	return parsedURL.Scheme == "http" || parsedURL.Scheme == "https"
}
