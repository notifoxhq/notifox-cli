package app

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/notifoxhq/notifox-go"
)

const (
	maxRetries   = 3
	initialDelay = 1 * time.Second
)

// App represents the application layer
type App struct {
	client *notifox.Client
}

// New creates a new App instance
func New(version string) (*App, error) {
	// Get API key from environment
	apiKey := os.Getenv("NOTIFOX_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("NOTIFOX_API_KEY environment variable is required")
	}

	// Build client options
	var opts []notifox.ClientOption

	opts = append(opts, notifox.WithUserAgent(fmt.Sprintf("notifox-cli/%s", version)))

	// Optional base URL override
	if baseURL := os.Getenv("NOTIFOX_BASE_URL"); baseURL != "" {
		opts = append(opts, notifox.WithBaseURL(baseURL))
	}

	// Create client
	client, err := notifox.NewClient(apiKey, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create notifox client: %w", err)
	}

	return &App{
		client: client,
	}, nil
}

// SendAlert sends an alert to the specified audience via the specified channel
func (a *App) SendAlert(audience, channel, message string, verbose bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Convert channel string to notifox.Channel type
	var ch notifox.Channel
	switch channel {
	case "sms":
		ch = notifox.SMS
	case "email":
		ch = notifox.Email
	default:
		return fmt.Errorf("invalid channel: %s (must be 'sms' or 'email')", channel)
	}

	// Send alert with retries for transient errors only
	req := notifox.AlertRequest{
		Audience: audience,
		Alert:    message,
		Channel:  ch,
	}

	var resp *notifox.AlertResponse
	var lastErr error
	delay := initialDelay

	for attempt := 0; attempt <= maxRetries; attempt++ {
		resp, lastErr = a.client.SendAlertWithOptions(ctx, req)
		if lastErr == nil {
			break
		}
		if !isRetryable(lastErr) || attempt == maxRetries {
			return handleError(lastErr)
		}
		time.Sleep(delay)
		delay *= 2
	}

	if lastErr != nil {
		return handleError(lastErr)
	}

	// Print success details only if verbose
	if verbose {
		fmt.Printf("Message ID: %s\n", resp.MessageID)
		fmt.Printf("Cost: $%.3f %s\n", resp.Cost, resp.Currency)
		fmt.Printf("Parts: %d\n", resp.Parts)
	}

	return nil
}

// isRetryable returns true only for errors that might succeed on retry.
// We do not retry: bad request (4xx), auth, insufficient balance, rate limit.
// We do retry: connection errors, internal server error (5xx).
func isRetryable(err error) bool {
	switch e := err.(type) {
	case *notifox.NotifoxConnectionError:
		return true
	case *notifox.NotifoxAPIError:
		return e.StatusCode >= 500
	default:
		return false
	}
}

// handleError converts SDK errors to user-friendly messages
func handleError(err error) error {
	switch e := err.(type) {
	case *notifox.NotifoxAuthenticationError:
		return fmt.Errorf("authentication failed: %s (status: %d)", e.ResponseText, e.StatusCode)
	case *notifox.NotifoxInsufficientBalanceError:
		return fmt.Errorf("insufficient balance: %s", e.ResponseText)
	case *notifox.NotifoxRateLimitError:
		return fmt.Errorf("rate limit exceeded: %s", e.ResponseText)
	case *notifox.NotifoxAPIError:
		return fmt.Errorf("API error: %s (status: %d)", e.ResponseText, e.StatusCode)
	case *notifox.NotifoxConnectionError:
		return fmt.Errorf("connection error: %v", e.Err)
	default:
		return fmt.Errorf("error: %v", err)
	}
}
