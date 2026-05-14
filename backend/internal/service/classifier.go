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

func normalizePriority(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "high", "medium", "low":
		return strings.ToLower(strings.TrimSpace(value))
	default:
		return "medium"
	}
}
