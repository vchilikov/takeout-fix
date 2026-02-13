package mediaext

import (
	"strings"
	"sync"
	"testing"
)

func TestParseWritableExtensionSet(t *testing.T) {
	output := "Writable file types:\n  3GP AVI HEIC JPG MP4 MOV WEBP\n"
	got := parseWritableExtensionSet(output)

	for _, ext := range []string{".3gp", ".avi", ".heic", ".jpg", ".mp4", ".mov", ".webp"} {
		if _, ok := got[ext]; !ok {
			t.Fatalf("expected %s in writable set, got %v", ext, got)
		}
	}
	if _, ok := got[".writable"]; ok {
		t.Fatalf("unexpected non-extension token in set: %v", got)
	}
}

func TestIsWritableExtension_UsesCachedList(t *testing.T) {
	origRun := runListWritableTypes
	origSet := writableExtSet
	origErr := writableExtLoadErr
	defer func() {
		runListWritableTypes = origRun
		writableOnce = sync.Once{}
		writableExtSet = origSet
		writableExtLoadErr = origErr
	}()

	writableOnce = sync.Once{}
	writableExtSet = nil
	writableExtLoadErr = nil

	calls := 0
	runListWritableTypes = func() (string, error) {
		calls++
		return "Writable file types:\n  JPG MP4 HEIC\n", nil
	}

	ok, err := IsWritableExtension(".jpg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatalf("expected .jpg to be writable")
	}

	ok, err = IsWritableExtension("AVI")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatalf("expected AVI to be non-writable")
	}

	if calls != 1 {
		t.Fatalf("expected cached list to be loaded once, got %d", calls)
	}
}

func TestIsWritableToken(t *testing.T) {
	tests := []struct {
		token string
		want  bool
	}{
		{token: "JPG", want: true},
		{token: "3GP", want: true},
		{token: "Writable", want: false},
		{token: "file", want: false},
		{token: "XMP:", want: false},
		{token: strings.ToLower("MP4"), want: false},
	}

	for _, tt := range tests {
		if got := isWritableToken(tt.token); got != tt.want {
			t.Fatalf("isWritableToken(%q) = %v, want %v", tt.token, got, tt.want)
		}
	}
}
