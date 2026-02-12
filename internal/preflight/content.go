package preflight

import (
	"errors"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/vchilikov/takeout-fix/internal/mediaext"
)

var (
	walkDirForContent = filepath.WalkDir
	errFoundMedia     = errors.New("processable media found")
)

// HasProcessableTakeout returns true when a folder looks like extracted
// Takeout content by finding at least one media evidence item.
func HasProcessableTakeout(path string) (bool, error) {
	err := walkDirForContent(path, func(_ string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}

		ext := filepath.Ext(d.Name())
		for _, supportedExt := range mediaext.Supported {
			if strings.EqualFold(ext, supportedExt) {
				return errFoundMedia
			}
		}
		return nil
	})
	if err == nil {
		return false, nil
	}
	if errors.Is(err, errFoundMedia) {
		return true, nil
	}
	return false, err
}
