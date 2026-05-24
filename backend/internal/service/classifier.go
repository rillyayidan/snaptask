package service

import (
	"strings"

	"snaptask/backend/internal/model"
)

func NormalizeItems(items []model.ExtractedItem) []model.ExtractedItem {
	normalized := make([]model.ExtractedItem, 0, len(items))
	seen := make(map[string]int, len(items))
	for _, item := range items {
		item.Type = normalizeType(item.Type)
		item.Priority = normalizePriority(item.Priority)
		item.Title = strings.TrimSpace(item.Title)
		item.Detail = strings.TrimSpace(item.Detail)
		item.DueDate = normalizeDueDateValue(item.DueDate)
		if item.Title == "" {
			continue
		}
		key := dedupeKey(item)
		if existingIndex, ok := seen[key]; ok {
			normalized[existingIndex] = mergeDuplicateItem(normalized[existingIndex], item)
			continue
		}
		seen[key] = len(normalized)
		normalized = append(normalized, item)
	}
	return normalized
}

func normalizeType(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "event", "calendar", "meeting", "appointment", "schedule":
		return "event"
	case "note", "memo", "info", "reference":
		return "note"
	default:
		return "task"
	}
}

func normalizeDueDateValue(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	normalized := strings.ToLower(strings.Trim(trimmed, `"'`))
	switch normalized {
	case "", "null", "nil", "none", "n/a", "na", "-", "unknown", "not specified":
		return nil
	}
	return &trimmed
}

func normalizePriority(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "high", "urgent", "important":
		return "high"
	case "medium", "normal":
		return "medium"
	case "low", "minor":
		return "low"
	case "":
		return "medium"
	default:
		normalized := strings.ToLower(strings.TrimSpace(value))
		switch {
		case strings.Contains(normalized, "urgent") || strings.Contains(normalized, "high"):
			return "high"
		case strings.Contains(normalized, "low"):
			return "low"
		case strings.Contains(normalized, "medium") || strings.Contains(normalized, "normal"):
			return "medium"
		}
		return "medium"
	}
}

func dedupeKey(item model.ExtractedItem) string {
	title := strings.ToLower(strings.Join(strings.Fields(item.Title), " "))
	dueDate := ""
	if item.DueDate != nil {
		dueDate = strings.ToLower(strings.TrimSpace(*item.DueDate))
	}
	return item.Type + "|" + title + "|" + dueDate
}

func mergeDuplicateItem(existing model.ExtractedItem, incoming model.ExtractedItem) model.ExtractedItem {
	if len(strings.TrimSpace(incoming.Detail)) > len(strings.TrimSpace(existing.Detail)) {
		existing.Detail = incoming.Detail
	}
	if existing.DueDate == nil && incoming.DueDate != nil {
		existing.DueDate = incoming.DueDate
	}
	if priorityRank(incoming.Priority) > priorityRank(existing.Priority) {
		existing.Priority = incoming.Priority
	}
	return existing
}

func priorityRank(value string) int {
	switch value {
	case "high":
		return 3
	case "medium":
		return 2
	case "low":
		return 1
	default:
		return 0
	}
}
