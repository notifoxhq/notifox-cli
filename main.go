package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mathis/notifox-cli/internal/app"
	"github.com/spf13/cobra"
)

var version = "v0.0.0"

func main() {
	rootCmd := &cobra.Command{
		Use:   "notifox",
		Short: "A CLI for sending internal alerts through Notifox",
		Long:  "A small, Unix-friendly CLI for sending internal alerts through Notifox (SMS + Email).",
		Run: func(cmd *cobra.Command, args []string) {
			// Show help if no subcommand provided
			cmd.Help()
		},
	}

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}

	sendCmd := &cobra.Command{
		Use:   "send",
		Short: "Send an alert",
		Long:  "Send an alert to a specified audience via SMS or email.",
		RunE:  runSend,
	}

	sendCmd.Flags().StringP("audience", "a", "", "audience to send the alert to")
	sendCmd.Flags().StringP("channel", "c", "", "channel to send through (sms|email)")
	sendCmd.Flags().StringP("message", "m", "", "message to send")
	sendCmd.Flags().BoolP("verbose", "v", false, "show detailed output (message ID, cost, parts)")

	rootCmd.AddCommand(sendCmd)
	rootCmd.AddCommand(versionCmd)

	// Enable --version flag (cobra handles this automatically when Version is set)
	rootCmd.Version = version
	rootCmd.SetVersionTemplate(version + "\n")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runSend(cmd *cobra.Command, args []string) error {
	// Get flags
	audience, _ := cmd.Flags().GetString("audience")
	channel, _ := cmd.Flags().GetString("channel")
	message, _ := cmd.Flags().GetString("message")
	verbose, _ := cmd.Flags().GetBool("verbose")

	// Get values from flags or environment variables (flags override env vars, like AWS CLI)
	finalAudience := audience
	if finalAudience == "" {
		finalAudience = os.Getenv("NOTIFOX_AUDIENCE")
	}

	finalChannel := channel
	if finalChannel == "" {
		finalChannel = os.Getenv("NOTIFOX_CHANNEL")
	}

	// Read message from stdin or flag
	msg, err := readMessage(message)
	if err != nil {
		return fmt.Errorf("error reading message: %w", err)
	}

	// Validate inputs
	if finalAudience == "" {
		return fmt.Errorf("audience is required (use -a/--audience or set NOTIFOX_AUDIENCE)")
	}
	if finalChannel == "" {
		return fmt.Errorf("channel is required (use -c/--channel or set NOTIFOX_CHANNEL)")
	}
	if msg == "" {
		return fmt.Errorf("message is required (provide via -m/--message or stdin)")
	}

	// Validate channel
	if finalChannel != "sms" && finalChannel != "email" {
		return fmt.Errorf("channel must be 'sms' or 'email', got '%s'", finalChannel)
	}

	// Create app and send alert
	notifoxApp, err := app.New()
	if err != nil {
		return fmt.Errorf("error initializing app: %w", err)
	}

	err = notifoxApp.SendAlert(finalAudience, finalChannel, msg, verbose)
	if err != nil {
		return fmt.Errorf("error sending alert: %w", err)
	}

	return nil
}

func readMessage(messageFlag string) (string, error) {
	// Precedence: --message flag wins, otherwise stdin
	if messageFlag != "" {
		return messageFlag, nil
	}

	// Check if stdin is a pipe (not a terminal)
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}

	// If stdin is a pipe (not a character device), read it
	// This handles: echo "text" | notifox send ...
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(data)), nil
	}

	// Stdin is a terminal, so no input from pipe
	return "", nil
}
