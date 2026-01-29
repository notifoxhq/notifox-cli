package app

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/notifoxhq/notifox-go"
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

	// Send alert using SDK
	resp, err := a.client.SendAlertWithOptions(ctx, notifox.AlertRequest{
		Audience: audience,
		Alert:    message,
		Channel:  ch,
	})

	if err != nil {
		return handleError(err)
	}

	// Print success details only if verbose
	if verbose {
		fmt.Printf("Message ID: %s\n", resp.MessageID)
		fmt.Printf("Cost: $%.3f %s\n", resp.Cost, resp.Currency)
		fmt.Printf("Parts: %d\n", resp.Parts)
	}

	return nil
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
