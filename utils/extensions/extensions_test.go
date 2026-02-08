package extensions

import (
	"errors"
	"slices"
	"strings"
	"testing"
)

func TestAreExtensionsCompatible(t *testing.T) {
	tests := []struct {
		ext1 string
		ext2 string
		want bool
	}{
		{".jpg", ".jpeg", true},
		{".mov", ".mp4", true},
		{".PNG", ".png", true},
		{".jpg", ".png", false},
	}

	for _, tt := range tests {
		got := areExtensionsCompatible(tt.ext1, tt.ext2)
		if got != tt.want {
			t.Fatalf("compatibility mismatch for %q/%q: want %v, got %v", tt.ext1, tt.ext2, tt.want, got)
		}
	}
}

func TestGenerateRandomSuffix(t *testing.T) {
	got, err := generateRandomSuffix()
	if err != nil {
		t.Fatalf("generateRandomSuffix error: %v", err)
	}
	if len(got) != 5 {
		t.Fatalf("expected suffix length 5, got %d", len(got))
	}
	for _, r := range got {
		if !strings.ContainsRune("abcdefghijklmnopqrstuvwxyz0123456789", r) {
			t.Fatalf("unexpected rune in suffix: %q", r)
		}
	}
}

func TestParseFileTypeExtension(t *testing.T) {
	tests := []struct {
		name   string
		output string
		want   string
	}{
		{
			name:   "plain extension",
			output: ".jpg\n",
			want:   ".jpg",
		},
		{
			name:   "without dot",
			output: "jpg\n",
			want:   ".jpg",
		},
		{
			name:   "ignores warnings",
			output: "Warning: duplicate tags\njpg\n",
			want:   ".jpg",
		},
		{
			name:   "only errors",
			output: "Error: bad file\n",
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseFileTypeExtension(tt.output); got != tt.want {
				t.Fatalf("parseFileTypeExtension mismatch: want %q, got %q", tt.want, got)
			}
		})
	}
}

func TestGetNewExtension_DoesNotUseDoubleDashSeparator(t *testing.T) {
	var gotArgs []string
	run := func(args []string) (string, error) {
		gotArgs = append([]string(nil), args...)
		return ".jpg\n", nil
	}

	ext, err := getNewExtension("photo.jpg", run)
	if err != nil {
		t.Fatalf("getNewExtension returned error: %v", err)
	}
	if ext != ".jpg" {
		t.Fatalf("extension mismatch: want .jpg, got %s", ext)
	}
	if slices.Contains(gotArgs, "--") {
		t.Fatalf("did not expect -- separator in exiftool args: %v", gotArgs)
	}
}

func TestGetNewExtension_RequiresRunner(t *testing.T) {
	if _, err := getNewExtension("photo.jpg", nil); err == nil {
		t.Fatalf("expected error for nil runner")
	}
}

func TestGetNewExtension_PropagatesRunnerError(t *testing.T) {
	run := func([]string) (string, error) {
		return "", errors.New("boom")
	}
	if _, err := getNewExtension("photo.jpg", run); err == nil {
		t.Fatalf("expected runner error")
	}
}
