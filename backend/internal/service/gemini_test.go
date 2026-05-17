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

func TestParseExtractedItemsAcceptsMarkdownFencedArray(t *testing.T) {
	items, err := parseExtractedItems("```json\n[{\"type\":\"task\",\"title\":\"Send deck\",\"priority\":\"high\"}]\n```")
	if err != nil {
		t.Fatalf("expected fenced JSON to parse: %v", err)
	}
	if len(items) != 1 || items[0].Title != "Send deck" {
		t.Fatalf("expected parsed task item, got %#v", items)
	}
}

func TestParseExtractedItemsAcceptsWrappedItemsArray(t *testing.T) {
	items, err := parseExtractedItems("{\"items\":[{\"type\":\"event\",\"title\":\"Friday sync\",\"priority\":\"medium\"}]}")
	if err != nil {
		t.Fatalf("expected wrapped items JSON to parse: %v", err)
	}
	if len(items) != 1 || items[0].Type != "event" {
		t.Fatalf("expected wrapped event item, got %#v", items)
	}
}
