package scripts

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	toolkit = map[string]any{
		"sha256":             Sha256,
		"sha512":             Sha512,
		"uuid":               UUIDv4,
		"parseTime":          ParseTime,
		"durationFromMillis": DurationFromMillis,
		"now":                Now,
	}
)

func Sha256(value string) string {
	hash := sha256.New()
	_, _ = hash.Write([]byte(strings.ToLower(value)))
	bytes := hash.Sum(nil)
	return fmt.Sprintf("%x", bytes)
}

func Sha512(value string) string {
	hash := sha512.New()
	_, _ = hash.Write([]byte(strings.ToLower(value)))
	bytes := hash.Sum(nil)
	return fmt.Sprintf("%x", bytes)
}

func UUIDv4() string {
	return uuid.NewString()
}

func DurationFromMillis(ms int) time.Duration {
	return time.Duration(ms) * time.Millisecond
}

func Now() time.Time {
	return time.Now()
}

var (
	formats = []string{
		"2006-01-02 15:04:05",
		"02.01.2006 15:04:05",
		"2006-01-02T15:04:05.999Z07:00",
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05.000Z",
		"2006-01-02",
		"01-02-2006",
		"02.01.2006",
	}
)

func ParseTime(value string, extraFormats ...string) *time.Time {
	allFormats := formats
	if len(extraFormats) > 0 {
		allFormats = make([]string, len(formats))
		copy(allFormats, formats)
		allFormats = append(allFormats, extraFormats...)
	}

	for _, format := range allFormats {
		t, err := time.Parse(format, value)
		if err == nil {
			return &t
		}
	}

	return nil
}
