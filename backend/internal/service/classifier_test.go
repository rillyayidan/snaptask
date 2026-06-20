package service

import (
	"strings"
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
	if items[0].Priority != "high" {
		t.Fatalf("expected urgent priority alias to normalize to high, got %q", items[0].Priority)
	}
	if items[0].DueDate == nil || *items[0].DueDate != "2026-05-15T15:00:00+07:00" {
		t.Fatalf("expected due date to be trimmed, got %#v", items[0].DueDate)
	}
}

func TestNormalizeItemsDropsBlankTitlesAndEmptyDueDates(t *testing.T) {
	emptyDueDate := "   "
	items := NormalizeItems([]model.ExtractedItem{
		{Title: "   ", Detail: "   ", DueDate: &emptyDueDate},
		{Title: "Call client", DueDate: &emptyDueDate},
	})

	if len(items) != 1 {
		t.Fatalf("expected only titled item to remain, got %#v", items)
	}
	if items[0].DueDate != nil {
		t.Fatalf("expected empty due date to normalize to nil, got %#v", items[0].DueDate)
	}
}

func TestNormalizeItemsBuildsTitleFromDetail(t *testing.T) {
	items := NormalizeItems([]model.ExtractedItem{{
		Type:     "task",
		Title:    "   ",
		Detail:   "  Please send the signed vendor agreement before lunch tomorrow.  ",
		Priority: "normal",
	}})

	if len(items) != 1 {
		t.Fatalf("expected detail-only item to remain, got %#v", items)
	}
	if items[0].Title != "Please send the signed vendor agreement before lunch tomorrow." {
		t.Fatalf("expected fallback title from detail, got %q", items[0].Title)
	}
	if items[0].Detail != "Please send the signed vendor agreement before lunch tomorrow." {
		t.Fatalf("expected detail to stay trimmed, got %q", items[0].Detail)
	}
}

func TestNormalizeItemsTruncatesFallbackTitleFromDetail(t *testing.T) {
	items := NormalizeItems([]model.ExtractedItem{{
		Detail: "Confirm the pitch deck, pricing appendix, legal notes, onboarding dates, success metrics, and stakeholder list before the 3 PM prep call.",
	}})

	if len(items) != 1 {
		t.Fatalf("expected detail-only item to remain, got %#v", items)
	}
	if len(items[0].Title) > 83 {
		t.Fatalf("expected fallback title to be truncated, got %q", items[0].Title)
	}
	if !strings.HasSuffix(items[0].Title, "...") {
		t.Fatalf("expected truncated fallback title to end with ellipsis, got %q", items[0].Title)
	}
}

func TestNormalizeItemsDropsNullishDueDateStrings(t *testing.T) {
	for _, raw := range []string{"null", "N/A", "-", "unknown", `"null"`} {
		raw := raw
		items := NormalizeItems([]model.ExtractedItem{{
			Title:   "Send recap",
			DueDate: &raw,
		}})

		if len(items) != 1 {
			t.Fatalf("expected one item for due date %q, got %#v", raw, items)
		}
		if items[0].DueDate != nil {
			t.Fatalf("expected due date %q to normalize to nil, got %#v", raw, items[0].DueDate)
		}
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

func TestNormalizeItemsAcceptsCommonGeminiAliases(t *testing.T) {
	items := NormalizeItems([]model.ExtractedItem{
		{Type: "meeting", Title: "Client sync", Priority: "urgent"},
		{Type: "memo", Title: "Budget note", Priority: "minor"},
	})

	if len(items) != 2 {
		t.Fatalf("expected 2 normalized items, got %#v", items)
	}
	if items[0].Type != "event" || items[0].Priority != "high" {
		t.Fatalf("expected meeting/urgent aliases to normalize to event/high, got %#v", items[0])
	}
	if items[1].Type != "note" || items[1].Priority != "low" {
		t.Fatalf("expected memo/minor aliases to normalize to note/low, got %#v", items[1])
	}
}

func TestNormalizeItemsInfersPriorityFromVerboseLabels(t *testing.T) {
	items := NormalizeItems([]model.ExtractedItem{
		{Title: "Pay vendor", Priority: "High priority"},
		{Title: "Read thread", Priority: "low confidence"},
	})

	if items[0].Priority != "high" {
		t.Fatalf("expected verbose high priority label to normalize to high, got %#v", items[0])
	}
	if items[1].Priority != "low" {
		t.Fatalf("expected verbose low priority label to normalize to low, got %#v", items[1])
	}
}
