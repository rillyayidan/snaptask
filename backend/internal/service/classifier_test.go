package service

import (
	"testing"

	"snaptask/backend/internal/model"
)

func TestNormalizeItemsTrimsMetadata(t *testing.T) {
	dueDate := " 2026-05-15T15:00:00+07:00 "
	items := NormalizeItems([]model.ExtractedItem{{
		Type:     " calendar ",
		Title:    " Send deck ",
		Detail:   " Original message ",
		DueDate:  &dueDate,
		Priority: " urgent ",
	}})

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].Type != "event" {
		t.Fatalf("expected calendar alias to normalize to event, got %q", items[0].Type)
	}
	if items[0].Title != "Send deck" || items[0].Detail != "Original message" {
		t.Fatalf("expected title and detail to be trimmed, got %#v", items[0])
	}
	if items[0].Priority != "medium" {
		t.Fatalf("expected invalid priority to normalize to medium, got %q", items[0].Priority)
	}
	if items[0].DueDate == nil || *items[0].DueDate != "2026-05-15T15:00:00+07:00" {
		t.Fatalf("expected due date to be trimmed, got %#v", items[0].DueDate)
	}
}

func TestNormalizeItemsDropsBlankTitlesAndEmptyDueDates(t *testing.T) {
	emptyDueDate := "   "
	items := NormalizeItems([]model.ExtractedItem{
		{Title: "   ", DueDate: &emptyDueDate},
		{Title: "Call client", DueDate: &emptyDueDate},
	})

	if len(items) != 1 {
		t.Fatalf("expected only titled item to remain, got %#v", items)
	}
	if items[0].DueDate != nil {
		t.Fatalf("expected empty due date to normalize to nil, got %#v", items[0].DueDate)
	}
}

func TestNormalizeItemsDeduplicatesSameTitleTypeAndDueDate(t *testing.T) {
	dueDate := "2026-05-23T15:00:00+07:00"
	items := NormalizeItems([]model.ExtractedItem{
		{
			Type:     "task",
			Title:    "Send invoice",
			Detail:   "short",
			DueDate:  &dueDate,
			Priority: "medium",
		},
		{
			Type:     "task",
			Title:    "  send   invoice ",
			Detail:   "Send the revised invoice to Nadia before the client call.",
			DueDate:  &dueDate,
			Priority: "high",
		},
	})

	if len(items) != 1 {
		t.Fatalf("expected duplicate items to collapse, got %#v", items)
	}
	if items[0].Priority != "high" {
		t.Fatalf("expected merged item to keep highest priority, got %#v", items[0])
	}
	if items[0].Detail != "Send the revised invoice to Nadia before the client call." {
		t.Fatalf("expected merged item to keep richer detail, got %#v", items[0])
	}
}

func TestNormalizeItemsKeepsSameTitleWhenDueDateDiffers(t *testing.T) {
	friday := "2026-05-23"
	monday := "2026-05-26"
	items := NormalizeItems([]model.ExtractedItem{
		{Type: "task", Title: "Follow up", DueDate: &friday},
		{Type: "task", Title: "Follow up", DueDate: &monday},
	})

	if len(items) != 2 {
		t.Fatalf("expected items with different due dates to remain separate, got %#v", items)
	}
}
