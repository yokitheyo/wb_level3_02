package validation

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type URLValidator struct{}

func NewURLValidator() *URLValidator {
	return &URLValidator{}
}

func (v *URLValidator) ValidateURL(url string) error {
	if strings.TrimSpace(url) == "" {
		return errors.New("url cannot be empty")
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return errors.New("url must start with http:// or https://")
	}

	if len(url) > 2048 {
		return errors.New("url is too long (max 2048 characters)")
	}

	return nil
}

type ShortCodeValidator struct{}

func NewShortCodeValidator() *ShortCodeValidator {
	return &ShortCodeValidator{}
}

func (v *ShortCodeValidator) ValidateShortCode(code string) error {
	if code == "" {
		return nil
	}

	if len(code) < 3 || len(code) > 50 {
		return fmt.Errorf("short code must be between 3 and 50 characters, got %d", len(code))
	}

	pattern := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !pattern.MatchString(code) {
		return errors.New("short code can only contain letters, numbers, dashes and underscores")
	}

	return nil
}
