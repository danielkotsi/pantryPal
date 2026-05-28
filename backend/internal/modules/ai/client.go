package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"pantrypal/backend/internal/config"
)

var (
	ErrMissingAPIKey      = errors.New("gemini api key is not configured")
	ErrEmptyPrompt        = errors.New("gemini prompt must not be empty")
	ErrUnexpectedResponse = errors.New("gemini returned no text response")
)

type Client struct {
	apiKey         string
	model          string
	baseURL        string
	timeout        time.Duration
	retryMax       int
	retryBackoff   time.Duration
	responseFormat string
	httpClient     *http.Client
	Endpoint       string
}

type GenerateRequest struct {
	Prompt           string
	ResponseMIMEType string
	Temperature      *float64
}

type GenerateResponse struct {
	Text       string
	Model      string
	Latency    time.Duration
	StatusCode int
	RequestID  string
}

type geminiGenerateContentRequest struct {
	Contents          []geminiContent          `json:"contents"`
	GenerationConfig  geminiGenerationConfig   `json:"generationConfig,omitempty"`
	SystemInstruction *geminiSystemInstruction `json:"systemInstruction,omitempty"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiSystemInstruction struct {
	Parts []geminiPart `json:"parts"`
}

type geminiGenerationConfig struct {
	CandidateCount   int      `json:"candidateCount,omitempty"`
	Temperature      *float64 `json:"temperature,omitempty"`
	ResponseMIMEType string   `json:"responseMimeType,omitempty"`
}

type geminiGenerateContentResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
	UsageMetadata struct {
		PromptTokenCount     int `json:"promptTokenCount"`
		CandidatesTokenCount int `json:"candidatesTokenCount"`
		TotalTokenCount      int `json:"totalTokenCount"`
	} `json:"usageMetadata"`
}

type Config struct {
	APIKey         string
	Model          string
	BaseURL        string
	Timeout        time.Duration
	RetryMax       int
	RetryBackoff   time.Duration
	ResponseFormat string
}

func ConfigFromApp(cfg config.Config) Config {
	return Config{
		APIKey:         cfg.GeminiAPIKey,
		Model:          cfg.GeminiModel,
		BaseURL:        cfg.GeminiBaseURL,
		Timeout:        cfg.GeminiTimeout,
		RetryMax:       cfg.GeminiRetryMax,
		RetryBackoff:   cfg.GeminiRetryBackoff,
		ResponseFormat: cfg.GeminiResponseFormat,
	}
}

func NewClient(cfg Config) (*Client, error) {
	if strings.TrimSpace(cfg.APIKey) == "" {
		return nil, ErrMissingAPIKey
	}
	if strings.TrimSpace(cfg.Model) == "" {
		cfg.Model = "gemini-2.5-flash"
	}
	if strings.TrimSpace(cfg.BaseURL) == "" {
		cfg.BaseURL = "https://generativelanguage.googleapis.com"
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 40 * time.Second
	}
	if cfg.RetryMax < 0 {
		cfg.RetryMax = 0
	}
	if cfg.RetryBackoff <= 0 {
		cfg.RetryBackoff = 500 * time.Millisecond
	}
	if strings.TrimSpace(cfg.ResponseFormat) == "" {
		cfg.ResponseFormat = "application/json"
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.ResponseHeaderTimeout = cfg.Timeout
	transport.DialContext = (&net.Dialer{Timeout: cfg.Timeout}).DialContext
	endpoint := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=%s", cfg.APIKey)

	return &Client{
		apiKey:         cfg.APIKey,
		model:          cfg.Model,
		baseURL:        strings.TrimRight(cfg.BaseURL, "/"),
		timeout:        cfg.Timeout,
		retryMax:       cfg.RetryMax,
		retryBackoff:   cfg.RetryBackoff,
		responseFormat: cfg.ResponseFormat,
		httpClient: &http.Client{
			Timeout:   cfg.Timeout,
			Transport: transport,
		},
		Endpoint: endpoint,
	}, nil
}

func (c *Client) Generate(ctx context.Context, req GenerateRequest) (GenerateResponse, error) {
	prompt := strings.TrimSpace(req.Prompt)
	if prompt == "" {
		return GenerateResponse{}, ErrEmptyPrompt
	}

	requestBody := geminiGenerateContentRequest{
		Contents: []geminiContent{{
			Parts: []geminiPart{{Text: prompt}},
		}},
		SystemInstruction: &geminiSystemInstruction{
			Parts: []geminiPart{{Text: "Return only the requested content. Do not wrap JSON in markdown fences."}},
		},
		GenerationConfig: geminiGenerationConfig{
			CandidateCount:   1,
			Temperature:      req.Temperature,
			ResponseMIMEType: c.responseMIMEType(req),
		},
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return GenerateResponse{}, err
	}

	var lastErr error
	for attempt := 0; attempt <= c.retryMax; attempt++ {
		start := time.Now()
		response, err := c.doRequest(ctx, c.Endpoint, body)
		if err == nil {
			response.Latency = time.Since(start)
			return response, nil
		}
		lastErr = err
		if attempt == c.retryMax || !isRetryable(err) {
			break
		}

		select {
		case <-ctx.Done():
			return GenerateResponse{}, ctx.Err()
		case <-time.After(c.retryBackoff * time.Duration(attempt+1)):
		}
	}

	return GenerateResponse{}, lastErr
}

func (c *Client) doRequest(ctx context.Context, endpoint string, body []byte) (GenerateResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return GenerateResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return GenerateResponse{}, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return GenerateResponse{}, err
	}

	requestID := strings.TrimSpace(resp.Header.Get("x-request-id"))
	if requestID == "" {
		requestID = strings.TrimSpace(resp.Header.Get("x-goog-request-id"))
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return GenerateResponse{}, httpStatusError{
			StatusCode: resp.StatusCode,
			Body:       string(responseBody),
			RequestID:  requestID,
		}
	}

	var parsed geminiGenerateContentResponse
	if err := json.Unmarshal(responseBody, &parsed); err != nil {
		return GenerateResponse{}, err
	}

	text := extractResponseText(parsed)
	if text == "" {
		return GenerateResponse{}, ErrUnexpectedResponse
	}

	return GenerateResponse{
		Text:       text,
		Model:      c.model,
		StatusCode: resp.StatusCode,
		RequestID:  requestID,
	}, nil
}

func (c *Client) responseMIMEType(req GenerateRequest) string {
	if strings.TrimSpace(req.ResponseMIMEType) != "" {
		return req.ResponseMIMEType
	}
	return c.responseFormat
}

func extractResponseText(resp geminiGenerateContentResponse) string {
	parts := make([]string, 0)
	for _, candidate := range resp.Candidates {
		for _, part := range candidate.Content.Parts {
			text := strings.TrimSpace(part.Text)
			if text != "" {
				parts = append(parts, text)
			}
		}
	}
	return strings.TrimSpace(strings.Join(parts, "\n"))
}

func isRetryable(err error) bool {
	var statusErr httpStatusError
	if errors.As(err, &statusErr) {
		return statusErr.StatusCode == http.StatusTooManyRequests || statusErr.StatusCode >= http.StatusInternalServerError
	}

	var netErr net.Error
	return errors.As(err, &netErr)
}

type httpStatusError struct {
	StatusCode int
	Body       string
	RequestID  string
}

func (e httpStatusError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("gemini request failed with status %d (request_id=%s): %s", e.StatusCode, e.RequestID, strings.TrimSpace(e.Body))
	}
	return fmt.Sprintf("gemini request failed with status %d: %s", e.StatusCode, strings.TrimSpace(e.Body))
}
