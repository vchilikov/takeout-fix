package preflight

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestHasProcessableTakeout_TrueWhenMediaExists(t *testing.T) {
	restore := stubContentWalker()
	defer restore()

	root := t.TempDir()
	writeFile(t, root, "photo.jpg")

	ok, err := HasProcessableTakeout(root)
	if err != nil {
		t.Fatalf("HasProcessableTakeout error: %v", err)
	}
	if !ok {
		t.Fatalf("expected true when media exists")
	}
}

func TestHasProcessableTakeout_FalseWithoutMedia(t *testing.T) {
	restore := stubContentWalker()
	defer restore()

	root := t.TempDir()
	writeFile(t, root, "orphan.json")
	writeFile(t, root, "notes.txt")

	ok, err := HasProcessableTakeout(root)
	if err != nil {
		t.Fatalf("HasProcessableTakeout error: %v", err)
	}
	if ok {
		t.Fatalf("expected false when no supported media exists")
	}
}

func TestHasProcessableTakeout_PropagatesWalkError(t *testing.T) {
	restore := stubContentWalker()
	defer restore()

	walkDirForContent = func(string, fs.WalkDirFunc) error {
		return errors.New("walk failed")
	}

	ok, err := HasProcessableTakeout(t.TempDir())
	if err == nil {
		t.Fatalf("expected error")
	}
	if ok {
		t.Fatalf("expected false on error")
	}
}

func TestHasProcessableTakeout_StopsOnFirstMedia(t *testing.T) {
	restore := stubContentWalker()
	defer restore()

	walkCalls := 0
	walkDirForContent = func(path string, fn fs.WalkDirFunc) error {
		walkCalls++
		err := fn(filepath.Join(path, "photo.jpg"), fakeDirEntry{name: "photo.jpg"}, nil)
		if !errors.Is(err, errFoundMedia) {
			t.Fatalf("expected early-exit marker, got %v", err)
		}
		return err
	}

	ok, err := HasProcessableTakeout("/tmp/work")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !ok {
		t.Fatalf("expected true when first walked file is media")
	}
	if walkCalls != 1 {
		t.Fatalf("expected walker call once, got %d", walkCalls)
	}
}

func stubContentWalker() func() {
	orig := walkDirForContent
	return func() {
		walkDirForContent = orig
	}
}

func writeFile(t *testing.T, root string, rel string) {
	t.Helper()
	path := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte("x"), 0o600); err != nil {
		t.Fatalf("write %s: %v", rel, err)
	}
}

type fakeDirEntry struct {
	name string
}

func (f fakeDirEntry) Name() string               { return f.name }
func (f fakeDirEntry) IsDir() bool                { return false }
func (f fakeDirEntry) Type() fs.FileMode          { return 0 }
func (f fakeDirEntry) Info() (fs.FileInfo, error) { return fakeFileInfo(f), nil }

type fakeFileInfo struct {
	name string
}

func (f fakeFileInfo) Name() string       { return f.name }
func (f fakeFileInfo) Size() int64        { return 0 }
func (f fakeFileInfo) Mode() fs.FileMode  { return 0 }
func (f fakeFileInfo) ModTime() time.Time { return time.Time{} }
func (f fakeFileInfo) IsDir() bool        { return false }
func (f fakeFileInfo) Sys() any           { return nil }
