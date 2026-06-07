package service

import (
	"strings"
	"testing"
	"time"
)

func TestBuildExtractorPromptIncludesReferenceDateAndTimezone(t *testing.T) {
	now := time.Date(2026, 5, 20, 9, 30, 0, 0, time.FixedZone("WIB", 7*60*60))
	prompt := buildExtractorPrompt(now, "Asia/Jakarta")

	for _, expected := range []string{
		"Today is 2026-05-20.",
		"Current local time is 2026-05-20T09:30:00+07:00.",
		"Local timezone is Asia/Jakarta.",
		`"besok"`,
	} {
		if !strings.Contains(prompt, expected) {
			t.Fatalf("expected prompt to contain %q, got:\n%s", expected, prompt)
		}
	}
}

func TestNormalizeImageMimeTypeKeepsSupportedUploadType(t *testing.T) {
	pngHeader := []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}
	got := normalizeImageMimeType("image/png; charset=binary", "chat.png", pngHeader)
	if got != "image/png" {
		t.Fatalf("expected image/png, got %q", got)
	}
}

func TestNormalizeImageMimeTypeAcceptsJPEGAlias(t *testing.T) {
	got := normalizeImageMimeType("image/jpg", "chat.jpg", []byte("not inspected"))
	if got != "image/jpeg" {
		t.Fatalf("expected image/jpeg alias, got %q", got)
	}
}

func TestNormalizeImageMimeTypeFallsBackToDetectedType(t *testing.T) {
	jpegHeader := []byte{0xff, 0xd8, 0xff, 0xdb}
	got := normalizeImageMimeType("application/octet-stream", "chat.bin", jpegHeader)
	if got != "image/jpeg" {
		t.Fatalf("expected detected image/jpeg, got %q", got)
	}
}

func TestNormalizeImageMimeTypeFallsBackToHEICFilename(t *testing.T) {
	got := normalizeImageMimeType("application/octet-stream", "shared-screenshot.HEIC", []byte("heic bytes"))
	if got != "image/heic" {
		t.Fatalf("expected HEIC filename fallback, got %q", got)
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

func TestParseExtractedItemsAcceptsSingleItemObject(t *testing.T) {
	items, err := parseExtractedItems("{\"type\":\"task\",\"title\":\"Send invoice\",\"priority\":\"medium\"}")
	if err != nil {
		t.Fatalf("expected single item JSON to parse: %v", err)
	}
	if len(items) != 1 || items[0].Title != "Send invoice" {
		t.Fatalf("expected parsed single task item, got %#v", items)
	}
}

func TestParseExtractedItemsAcceptsPrefacedJSON(t *testing.T) {
	items, err := parseExtractedItems("Here are the extracted tasks:\n[{\"type\":\"task\",\"title\":\"Send invoice\",\"priority\":\"medium\"}]")
	if err != nil {
		t.Fatalf("expected prefaced JSON to parse: %v", err)
	}
	if len(items) != 1 || items[0].Title != "Send invoice" {
		t.Fatalf("expected parsed prefaced task, got %#v", items)
	}
}

func TestParseExtractedItemsAcceptsUTF8BOMPayload(t *testing.T) {
	items, err := parseExtractedItems("\ufeff[{\"type\":\"note\",\"title\":\"Read later\",\"priority\":\"low\"}]")
	if err != nil {
		t.Fatalf("expected BOM-prefixed JSON to parse: %v", err)
	}
	if len(items) != 1 || items[0].Title != "Read later" {
		t.Fatalf("expected parsed BOM-prefixed item, got %#v", items)
	}
}

func TestParseExtractedItemsAcceptsTrailingText(t *testing.T) {
	items, err := parseExtractedItems("[{\"type\":\"note\",\"title\":\"Remember budget\",\"priority\":\"low\"}]\nDone.")
	if err != nil {
		t.Fatalf("expected JSON with trailing text to parse: %v", err)
	}
	if len(items) != 1 || items[0].Type != "note" {
		t.Fatalf("expected parsed note item, got %#v", items)
	}
}

func TestExtractionTextFromResponseUsesLaterTextParts(t *testing.T) {
	response := geminiResponse{}
	response.Candidates = append(response.Candidates, struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	}{})
	response.Candidates[0].Content.Parts = append(response.Candidates[0].Content.Parts,
		struct {
			Text string `json:"text"`
		}{Text: "   "},
		struct {
			Text string `json:"text"`
		}{Text: "[{\"type\":\"task\",\"title\":\"Send brief\"}]"},
	)

	text, err := extractionTextFromResponse(response)
	if err != nil {
		t.Fatalf("expected extraction text: %v", err)
	}
	if text != "[{\"type\":\"task\",\"title\":\"Send brief\"}]" {
		t.Fatalf("expected later text part, got %q", text)
	}
}

func TestExtractionTextFromResponseJoinsTextParts(t *testing.T) {
	response := geminiResponse{}
	response.Candidates = append(response.Candidates, struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	}{})
	response.Candidates[0].Content.Parts = append(response.Candidates[0].Content.Parts,
		struct {
			Text string `json:"text"`
		}{Text: "["},
		struct {
			Text string `json:"text"`
		}{Text: "{\"type\":\"task\",\"title\":\"Pay invoice\"}]"},
	)

	text, err := extractionTextFromResponse(response)
	if err != nil {
		t.Fatalf("expected extraction text: %v", err)
	}
	if text != "[\n{\"type\":\"task\",\"title\":\"Pay invoice\"}]" {
		t.Fatalf("expected joined text parts, got %q", text)
	}
}
