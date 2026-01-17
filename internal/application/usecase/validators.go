package usecase

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

func GenerateShortCode() (string, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("failed to generate UUID: %w", err)
	}

	hash := md5.Sum([]byte(u.String()))
	encoded := base64.URLEncoding.EncodeToString(hash[:])
	return strings.TrimRight(encoded, "=")[:6], nil
}

type ShortCodeValidator struct {
	customCodeRegex *regexp.Regexp
}

func NewShortCodeValidator() *ShortCodeValidator {
	return &ShortCodeValidator{
		customCodeRegex: regexp.MustCompile(`^[a-zA-Z0-9_-]+$`),
	}
}

func (v *ShortCodeValidator) ValidateCustomShort(code string) error {
	if len(code) < 3 {
		return fmt.Errorf("short code must be at least 3 characters")
	}
	if len(code) > 50 {
		return fmt.Errorf("short code must be at most 50 characters")
	}
	if !v.customCodeRegex.MatchString(code) {
		return fmt.Errorf("short code can only contain alphanumeric characters, underscore, and hyphen")
	}
	return nil
}
