package main

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

const inspirobotApiUrl = "https://inspirobot.me/api?generate=true"

// generateInspirobotImage asks InspiroBot to generate an image.
// This will return a URL to the image, if successful.
func generateInspirobotImage() (*url.URL, error) {
	// Create an HTTP request.
	req, err := http.NewRequest("GET", inspirobotApiUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "discord-inspirobot/0.1.0 (Golang; https://github.com/eth-p/discord-inspirobot)")

	// Perform the HTTP request.
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Read the response URL.
	var sb strings.Builder
	_, err = io.Copy(&sb, res.Body)
	if err != nil {
		return nil, err
	}

	// Parse the response URL.
	responseUrl, err := url.Parse(sb.String())
	if err != nil {
		return nil, err
	}

	return responseUrl, nil
}
