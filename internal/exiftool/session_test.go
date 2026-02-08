package exiftool

import "testing"

func TestHasErrorLine(t *testing.T) {
	if !hasErrorLine("Warning: x\nError: boom\n") {
		t.Fatalf("expected error line detection")
	}
	if hasErrorLine("1 image files updated\n") {
		t.Fatalf("did not expect error line")
	}
}

func TestFirstErrorLine(t *testing.T) {
	got := firstErrorLine("Warning: x\nError: boom\n")
	if got != "Error: boom" {
		t.Fatalf("firstErrorLine mismatch: got %q", got)
	}
}
