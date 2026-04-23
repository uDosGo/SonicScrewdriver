package secrets

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// ProxyServer handles API proxy requests
type ProxyServer struct {
	secretStore  *SecretStore
	rateLimits   map[string]int
	callCounts   map[string]int
	lastReset    time.Time
	mu           sync.RWMutex
}

// ProxyRequest represents a proxy request
type ProxyRequest struct {
	Provider string                 `json:"provider"`
	Method   string                 `json:"method"`
	Path     string                 `json:"path"`
	Headers  map[string]string      `json:"headers"`
	Body     map[string]interface{} `json:"body"`
}

// ProxyResponse represents a proxy response
type ProxyResponse struct {
	Status  int                    `json:"status"`
	Body    map[string]interface{} `json:"body"`
	Headers map[string]string      `json:"headers"`
}

// NewProxyServer creates a new proxy server
func NewProxyServer(secretStore *SecretStore) *ProxyServer {
	return &ProxyServer{
		secretStore: secretStore,
		rateLimits: map[string]int{
			"openrouter": 60,
			"deepseek":   100,
			"gemini":     60,
			"github":    5000,
		},
		callCounts: make(map[string]int),
		lastReset:  time.Now(),
	}
}

// HandleRequest handles a proxy request
func (p *ProxyServer) HandleRequest(req ProxyRequest) (ProxyResponse, error) {
	// Check rate limits
	if err := p.checkRateLimit(req.Provider); err != nil {
		return ProxyResponse{Status: 429}, err
	}

	// Get API key from secret store
	apiKey, err := p.secretStore.GetSecret(req.Provider + "_api_key")
	if err != nil {
		return ProxyResponse{Status: 401}, fmt.Errorf("failed to get API key: %v", err)
	}

	// Construct the actual API request
	url := p.getProviderURL(req.Provider) + req.Path
	
	// Convert body to JSON
	bodyJSON, err := json.Marshal(req.Body)
	if err != nil {
		return ProxyResponse{Status: 400}, fmt.Errorf("failed to marshal body: %v", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest(req.Method, url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		return ProxyResponse{Status: 500}, fmt.Errorf("failed to create request: %v", err)
	}

	// Add headers
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// Add authorization header
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return ProxyResponse{Status: 502}, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ProxyResponse{Status: 500}, fmt.Errorf("failed to read response: %v", err)
	}

	// Parse response body
	var responseBody map[string]interface{}
	if err := json.Unmarshal(body, &responseBody); err != nil {
		responseBody = map[string]interface{}{
			"raw": string(body),
		}
	}

	// Increment call count
	p.incrementCallCount(req.Provider)

	return ProxyResponse{
		Status: resp.StatusCode,
		Body:   responseBody,
		Headers: map[string]string{
			"Content-Type": resp.Header.Get("Content-Type"),
		},
	}, nil
}

// checkRateLimit checks if rate limit is exceeded
func (p *ProxyServer) checkRateLimit(provider string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Reset counters if it's been more than 1 minute
	if time.Since(p.lastReset) > time.Minute {
		p.callCounts = make(map[string]int)
		p.lastReset = time.Now()
	}

	count := p.callCounts[provider]
	limit := p.rateLimits[provider]

	if count >= limit {
		return fmt.Errorf("rate limit exceeded for %s: %d/%d", provider, count, limit)
	}

	return nil
}

// incrementCallCount increments the call count for a provider
func (p *ProxyServer) incrementCallCount(provider string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.callCounts[provider]++
}

// getProviderURL returns the base URL for a provider
func (p *ProxyServer) getProviderURL(provider string) string {
	switch provider {
	case "openrouter":
		return "https://openrouter.ai/api/v1"
	case "deepseek":
		return "https://api.deepseek.com/v1"
	case "gemini":
		return "https://generativelanguage.googleapis.com/v1"
	case "github":
		return "https://api.github.com"
	default:
		return ""
	}
}

// GetStatus returns the current status of all proxy endpoints
func (p *ProxyServer) GetStatus() map[string]map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	status := make(map[string]map[string]interface{})

	for provider, limit := range p.rateLimits {
		count := p.callCounts[provider]
		status[provider] = map[string]interface{}{
			"calls":       count,
			"errors":      0,
			"rate_limit":  fmt.Sprintf("%d/min", limit),
			"status":      "healthy",
		}
	}

	return status
}

// GetHealth checks the health of all proxy endpoints
func (p *ProxyServer) GetHealth() map[string]string {
	health := make(map[string]string)

	providers := []string{"openrouter", "deepseek", "gemini", "github"}
	
	for _, provider := range providers {
		url := p.getProviderURL(provider)
		if url == "" {
			health[provider] = "unknown"
			continue
		}

		// Simple health check - just try to connect
		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Get(url)
		if err != nil {
			health[provider] = fmt.Sprintf("unhealthy: %v", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			health[provider] = "healthy"
		} else {
			health[provider] = fmt.Sprintf("unhealthy: status %d", resp.StatusCode)
		}
	}

	return health
}

// TestAPIKey tests if an API key is valid
func (p *ProxyServer) TestAPIKey(provider, apiKey string) error {
	url := p.getProviderURL(provider)
	if url == "" {
		return errors.New("unknown provider")
	}

	// Create a simple test request
	req, err := http.NewRequest("GET", url+"/models", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("API key test failed with status %d", resp.StatusCode)
}