package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"snaptask/backend/internal/model"
)

type GoogleService struct {
	client *http.Client
}

func NewGoogleService() *GoogleService {
	return &GoogleService{client: &http.Client{Timeout: 15 * time.Second}}
}

func (s *GoogleService) Push(ctx context.Context, accessToken string, items []model.ExtractedItem) []model.PushResult {
	items = NormalizeItems(items)
	results := make([]model.PushResult, len(items))
	var wg sync.WaitGroup

	for idx, item := range items {
		wg.Add(1)
		go func(i int, current model.ExtractedItem) {
			defer wg.Done()
			switch current.Type {
			case "event":
				results[i] = s.createEvent(ctx, accessToken, current)
			case "note":
				results[i] = skippedResult(current, "notes are review-only and are not pushed to Google")
			default:
				results[i] = s.createTask(ctx, accessToken, current)
			}
		}(idx, item)
	}

	wg.Wait()
	return results
}

func (s *GoogleService) createTask(ctx context.Context, token string, item model.ExtractedItem) model.PushResult {
	payload := map[string]any{
		"title": item.Title,
		"notes": item.Detail,
	}
	if due, ok := normalizeGoogleDueDate(item.DueDate); ok {
		payload["due"] = due
	}

	id, err := s.postGoogle(ctx, token, "https://tasks.googleapis.com/tasks/v1/lists/@default/tasks", payload)
	return pushResult(item, id, err)
}

func (s *GoogleService) createEvent(ctx context.Context, token string, item model.ExtractedItem) model.PushResult {
	start, end := calendarEventTimes(item.DueDate, time.Now().Add(24*time.Hour))
	payload := map[string]any{
		"summary":     item.Title,
		"description": item.Detail,
		"start":       start,
		"end":         end,
	}

	id, err := s.postGoogle(ctx, token, "https://www.googleapis.com/calendar/v3/calendars/primary/events", payload)
	return pushResult(item, id, err)
}

func (s *GoogleService) postGoogle(ctx context.Context, token, url string, payload map[string]any) (string, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("google api returned %d: %s", resp.StatusCode, string(respBody))
	}

	var decoded struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(respBody, &decoded); err != nil {
		return "", err
	}
	return decoded.ID, nil
}

func normalizeGoogleDueDate(value *string) (string, bool) {
	parsed, ok := parseDueDate(value)
	if !ok {
		return "", false
	}
	return parsed.Format(time.RFC3339), true
}

func calendarEventTimes(value *string, fallback time.Time) (map[string]string, map[string]string) {
	if value != nil {
		raw := strings.TrimSpace(*value)
		if parsed, err := time.Parse("2006-01-02", raw); err == nil {
			return map[string]string{"date": parsed.Format("2006-01-02")},
				map[string]string{"date": parsed.AddDate(0, 0, 1).Format("2006-01-02")}
		}
		if parsed, err := time.Parse(time.RFC3339, raw); err == nil {
			start := parsed.UTC()
			return map[string]string{"dateTime": start.Format(time.RFC3339)},
				map[string]string{"dateTime": start.Add(time.Hour).Format(time.RFC3339)}
		}
	}

	start := fallback.UTC()
	return map[string]string{"dateTime": start.Format(time.RFC3339)},
		map[string]string{"dateTime": start.Add(time.Hour).Format(time.RFC3339)}
}

func parseDueDate(value *string) (time.Time, bool) {
	if value == nil {
		return time.Time{}, false
	}
	raw := strings.TrimSpace(*value)
	if raw == "" {
		return time.Time{}, false
	}
	if parsed, err := time.Parse(time.RFC3339, raw); err == nil {
		return parsed.UTC(), true
	}
	if parsed, err := time.Parse("2006-01-02", raw); err == nil {
		return parsed.UTC(), true
	}
	return time.Time{}, false
}

func pushResult(item model.ExtractedItem, id string, err error) model.PushResult {
	result := model.PushResult{
		Title:  item.Title,
		Type:   item.Type,
		Status: "ok",
		ID:     id,
	}
	if err != nil {
		result.Status = "error"
		result.Error = err.Error()
	}
	return result
}

func skippedResult(item model.ExtractedItem, reason string) model.PushResult {
	return model.PushResult{
		Title:  item.Title,
		Type:   item.Type,
		Status: "skipped",
		Error:  reason,
	}
}
