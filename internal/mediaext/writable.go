package mediaext

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/vchilikov/takeout-fix/internal/exifcmd"
)

var (
	writableOnce       sync.Once
	writableExtSet     map[string]struct{}
	writableExtLoadErr error
)

var runListWritableTypes = func() (string, error) {
	bin, err := exifcmd.Resolve()
	if err != nil {
		return "", err
	}

	cmd := exec.Command(bin, "-listwf")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("run exiftool -listwf: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	return string(out), nil
}

// IsWritableExtension reports whether exiftool can write metadata into files
// with the provided extension (for example ".jpg" or ".avi").
func IsWritableExtension(ext string) (bool, error) {
	normalized := strings.ToLower(strings.TrimSpace(ext))
	if normalized == "" {
		return false, nil
	}
	if !strings.HasPrefix(normalized, ".") {
		normalized = "." + normalized
	}

	writableOnce.Do(func() {
		output, err := runListWritableTypes()
		if err != nil {
			writableExtLoadErr = err
			return
		}
		writableExtSet = parseWritableExtensionSet(output)
	})
	if writableExtLoadErr != nil {
		return false, writableExtLoadErr
	}

	_, ok := writableExtSet[normalized]
	return ok, nil
}

func parseWritableExtensionSet(output string) map[string]struct{} {
	set := make(map[string]struct{})
	for _, token := range strings.Fields(output) {
		cleaned := strings.Trim(token, " \t\r\n,;:()[]{}")
		if cleaned == "" {
			continue
		}
		if !isWritableToken(cleaned) {
			continue
		}
		set["."+strings.ToLower(cleaned)] = struct{}{}
	}
	return set
}

func isWritableToken(token string) bool {
	if token == "" || token != strings.ToUpper(token) {
		return false
	}

	hasLetter := false
	for _, r := range token {
		switch {
		case r >= 'A' && r <= 'Z':
			hasLetter = true
		case r >= '0' && r <= '9':
			// allowed
		default:
			return false
		}
	}

	return hasLetter
}
