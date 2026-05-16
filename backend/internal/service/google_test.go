package service

import (
	"context"
	"testing"
	"time"

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

func TestNormalizeGoogleDueDateAcceptsRFC3339WithOffset(t *testing.T) {
	raw := "2026-05-15T15:00:00+07:00"
	got, ok := normalizeGoogleDueDate(&raw)
	if !ok {
		t.Fatal("expected RFC3339 due date to be accepted")
	}
	if got != "2026-05-15T08:00:00Z" {
		t.Fatalf("expected UTC due date, got %q", got)
	}
}

func TestNormalizeGoogleDueDateAcceptsDateOnly(t *testing.T) {
	raw := "2026-05-15"
	got, ok := normalizeGoogleDueDate(&raw)
	if !ok {
		t.Fatal("expected date-only due date to be accepted")
	}
	if got != time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC).Format(time.RFC3339) {
		t.Fatalf("expected midnight UTC due date, got %q", got)
	}
}

func TestNormalizeGoogleDueDateRejectsInvalidDate(t *testing.T) {
	raw := "next Friday"
	if got, ok := normalizeGoogleDueDate(&raw); ok {
		t.Fatalf("expected invalid due date to be rejected, got %q", got)
	}
}
