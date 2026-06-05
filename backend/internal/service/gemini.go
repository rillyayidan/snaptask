package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	mimepkg "mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"snaptask/backend/internal/model"
)

const maxGeminiImageBytes = 10 * 1024 * 1024

const defaultExtractionTimezone = "Asia/Jakarta"

const extractorPromptBase = `You are a task extractor. Analyze this screenshot and extract ALL action items, commitments, and time-sensitive information.

Respond ONLY with a JSON array. No explanation, no markdown, no preamble.

Format:
[
  {
    "type": "task" | "event" | "note",
    "title": "short action title",
    "detail": "original text from screenshot",
    "due_date": "ISO 8601 or null",
    "priority": "high" | "medium" | "low"
  }
]

Due date rules:
- Use YYYY-MM-DD when only the date is known.
- Use RFC3339 with timezone offset when both date and time are known.
- Use null when the screenshot does not provide enough information.

Extract implicit commitments too. "Eh nanti kirim filenya ya" = task: "Kirim file ke [sender]".`

type GeminiService struct {
	apiKey   string
	model    string
	timezone string
	location *time.Location
	client   *http.Client
}

func NewGeminiService() *GeminiService {
	modelName := os.Getenv("GEMINI_MODEL")
	if modelName == "" {
		modelName = "gemini-2.0-flash"
	}
	timezone, location := configuredExtractionLocation()
	return &GeminiService{
		apiKey:   os.Getenv("GEMINI_API_KEY"),
		model:    modelName,
		timezone: timezone,
		location: location,
		client:   &http.Client{Timeout: 25 * time.Second},
	}
}

func (s *GeminiService) Extract(ctx context.Context, file multipart.File, header *multipart.FileHeader) ([]model.ExtractedItem, error) {
	if s.apiKey == "" {
		return nil, errors.New("GEMINI_API_KEY is not configured")
	}

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read image: %w", err)
	}
	if len(imageBytes) == 0 {
		return nil, errors.New("uploaded image is empty")
	}
	if len(imageBytes) > maxGeminiImageBytes {
		return nil, fmt.Errorf("uploaded image is too large; max size is %d MB", maxGeminiImageBytes/1024/1024)
	}

	mimeType := normalizeImageMimeType(header.Header.Get("Content-Type"), header.Filename, imageBytes)
	if !isSupportedImageMimeType(mimeType) {
		return nil, fmt.Errorf("unsupported image type %q; upload a PNG, JPEG, WebP, GIF, HEIC, or HEIF screenshot", mimeType)
	}

	payload := geminiRequest{
		Contents: []geminiContent{{
			Parts: []geminiPart{
				{Text: buildExtractorPrompt(time.Now().In(s.location), s.timezone)},
				{InlineData: &geminiInlineData{
					MimeType: mimeType,
					Data:     base64.StdEncoding.EncodeToString(imageBytes),
				}},
			},
		}},
		GenerationConfig: geminiGenerationConfig{
			Temperature:      0.1,
			ResponseMimeType: "application/json",
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal gemini request: %w", err)
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", s.model, s.apiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create gemini request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call gemini: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read gemini response: %w", err)
	}
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("gemini returned %d: %s", resp.StatusCode, string(respBody))
	}

	var geminiResp geminiResponse
	if err := json.Unmarshal(respBody, &geminiResp); err != nil {
		return nil, fmt.Errorf("decode gemini response: %w", err)
	}
	extractionText, err := extractionTextFromResponse(geminiResp)
	if err != nil {
		return nil, err
	}

	items, err := parseExtractedItems(extractionText)
	if err != nil {
		return nil, err
	}
	return NormalizeItems(items), nil
}

func buildExtractorPrompt(now time.Time, timezone string) string {
	timezone = strings.TrimSpace(timezone)
	if timezone == "" {
		timezone = "UTC"
	}
	return fmt.Sprintf(`%s

Reference date:
- Today is %s.
- Current local time is %s.
- Local timezone is %s.

Resolve relative dates and times against this reference, including Indonesian phrases like "besok", "lusa", "Jumat ini", and "Jumat depan".`, extractorPromptBase, now.Format("2006-01-02"), now.Format(time.RFC3339), timezone)
}

func configuredExtractionLocation() (string, *time.Location) {
	timezone := strings.TrimSpace(os.Getenv("SNAPTASK_TIMEZONE"))
	if timezone == "" {
		timezone = defaultExtractionTimezone
	}
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return "UTC", time.UTC
	}
	return timezone, location
}

func extractionTextFromResponse(response geminiResponse) (string, error) {
	for _, candidate := range response.Candidates {
		var parts []string
		for _, part := range candidate.Content.Parts {
			if text := strings.TrimSpace(part.Text); text != "" {
				parts = append(parts, text)
			}
		}
		if len(parts) > 0 {
			return strings.Join(parts, "\n"), nil
		}
	}
	return "", errors.New("gemini returned no extraction text")
}

func parseExtractedItems(text string) ([]model.ExtractedItem, error) {
	text = cleanGeminiJSONText(text)

	var items []model.ExtractedItem
	if err := json.Unmarshal([]byte(text), &items); err == nil {
		return items, nil
	}

	var wrapped struct {
		Items []model.ExtractedItem `json:"items"`
	}
	if err := json.Unmarshal([]byte(text), &wrapped); err != nil {
		return nil, fmt.Errorf("parse extracted JSON: %w", err)
	}
	if wrapped.Items == nil {
		return nil, errors.New("parse extracted JSON: missing items array")
	}
	return wrapped.Items, nil
}

func cleanGeminiJSONText(text string) string {
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)
	if json.Valid([]byte(text)) {
		return text
	}
	if payload, ok := extractJSONPayload(text); ok {
		return payload
	}
	return text
}

func extractJSONPayload(text string) (string, bool) {
	for index, char := range text {
		if char != '[' && char != '{' {
			continue
		}

		decoder := json.NewDecoder(strings.NewReader(text[index:]))
		var payload json.RawMessage
		if err := decoder.Decode(&payload); err == nil && len(payload) > 0 {
			return string(payload), true
		}
	}
	return "", false
}

func normalizeImageMimeType(uploaded string, filename string, imageBytes []byte) string {
	if uploaded != "" {
		if mediaType, _, err := mimepkg.ParseMediaType(uploaded); err == nil {
			uploaded = mediaType
		}
		uploaded = strings.ToLower(strings.TrimSpace(uploaded))
	}
	uploaded = normalizeImageMimeAlias(uploaded)
	if isSupportedImageMimeType(uploaded) {
		return uploaded
	}

	detected := strings.ToLower(http.DetectContentType(imageBytes))
	if isSupportedImageMimeType(detected) {
		return detected
	}

	if heicType := heicMimeTypeFromFilename(filename); heicType != "" {
		return heicType
	}
	return detected
}

func normalizeImageMimeAlias(value string) string {
	switch value {
	case "image/jpg", "image/pjpeg":
		return "image/jpeg"
	default:
		return value
	}
}

func heicMimeTypeFromFilename(filename string) string {
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".heic":
		return "image/heic"
	case ".heif":
		return "image/heif"
	default:
		return ""
	}
}

func isSupportedImageMimeType(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "image/png", "image/jpeg", "image/webp", "image/gif", "image/heic", "image/heif":
		return true
	default:
		return false
	}
}

type geminiRequest struct {
	Contents         []geminiContent        `json:"contents"`
	GenerationConfig geminiGenerationConfig `json:"generationConfig"`
}

type geminiGenerationConfig struct {
	Temperature      float64 `json:"temperature"`
	ResponseMimeType string  `json:"responseMimeType"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text       string            `json:"text,omitempty"`
	InlineData *geminiInlineData `json:"inline_data,omitempty"`
}

type geminiInlineData struct {
	MimeType string `json:"mime_type"`
	Data     string `json:"data"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}
