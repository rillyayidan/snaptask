package service

import (
	"strings"

	"snaptask/backend/internal/model"
)

func NormalizeItems(items []model.ExtractedItem) []model.ExtractedItem {
	normalized := make([]model.ExtractedItem, 0, len(items))
	for _, item := range items {
		item.Type = normalizeType(item.Type)
		item.Priority = normalizePriority(item.Priority)
		item.Title = strings.TrimSpace(item.Title)
		item.Detail = strings.TrimSpace(item.Detail)
		item.DueDate = normalizeDueDateValue(item.DueDate)
		if item.Title == "" {
			continue
		}
		normalized = append(normalized, item)
	}
	return normalized
}

func normalizeType(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "event", "calendar":
		return "event"
	case "note":
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
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func normalizePriority(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "high", "medium", "low":
		return strings.ToLower(strings.TrimSpace(value))
	default:
		return "medium"
	}
}
