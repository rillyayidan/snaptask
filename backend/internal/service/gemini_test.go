package service

import "testing"

func TestNormalizeImageMimeTypeKeepsSupportedUploadType(t *testing.T) {
	pngHeader := []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}
	got := normalizeImageMimeType("image/png; charset=binary", pngHeader)
	if got != "image/png" {
		t.Fatalf("expected image/png, got %q", got)
	}
}

func TestNormalizeImageMimeTypeFallsBackToDetectedType(t *testing.T) {
	jpegHeader := []byte{0xff, 0xd8, 0xff, 0xdb}
	got := normalizeImageMimeType("application/octet-stream", jpegHeader)
	if got != "image/jpeg" {
		t.Fatalf("expected detected image/jpeg, got %q", got)
	}
}

func TestSupportedImageMimeTypeRejectsPDF(t *testing.T) {
	if isSupportedImageMimeType("application/pdf") {
		t.Fatal("expected PDF MIME type to be rejected")
	}
}
