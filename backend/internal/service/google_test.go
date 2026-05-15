package service

import (
	"context"
	"testing"

	"snaptask/backend/internal/model"
)

func TestPushSkipsNotesWithoutCallingGoogle(t *testing.T) {
	service := NewGoogleService()
	results := service.Push(context.Background(), "unused-token", []model.ExtractedItem{{
		Type:     "note",
		Title:    "Keep this for review",
		Detail:   "Original note text",
		Priority: "low",
	}})

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != "skipped" {
		t.Fatalf("expected note to be skipped, got status %q", results[0].Status)
	}
	if results[0].Type != "note" {
		t.Fatalf("expected note type to be preserved, got %q", results[0].Type)
	}
	if results[0].Error == "" {
		t.Fatal("expected skipped note to include a reason")
	}
}
