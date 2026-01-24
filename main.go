package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mathis/notifox-cli/internal/app"
)

const (
	exitSuccess = 0
	exitFailure = 1
)

func main() {
	// Parse command
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(exitFailure)
	}

	command := os.Args[1]
	if command != "send" {
		fmt.Fprintf(os.Stderr, "Error: unknown command '%s'\n", command)
		fmt.Fprintf(os.Stderr, "Available commands: send\n")
		os.Exit(exitFailure)
	}

	// Parse flags for send command
	fs := flag.NewFlagSet("send", flag.ExitOnError)
	audience := fs.String("a", "", "audience to send the alert to")
	audienceLong := fs.String("audience", "", "audience to send the alert to")
	channel := fs.String("c", "", "channel to send through (sms|email)")
	channelLong := fs.String("channel", "", "channel to send through (sms|email)")
	message := fs.String("m", "", "message to send")
	messageLong := fs.String("message", "", "message to send")

	fs.Parse(os.Args[2:])

	// Merge short and long flags
	if *audienceLong != "" {
		*audience = *audienceLong
	}
	if *channelLong != "" {
		*channel = *channelLong
	}
	if *messageLong != "" {
		*message = *messageLong
	}

	// Read message from stdin or flag
	msg, err := readMessage(*message)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading message: %v\n", err)
		os.Exit(exitFailure)
	}

	// Validate inputs
	if *audience == "" {
		fmt.Fprintf(os.Stderr, "Error: -a/--audience is required\n")
		os.Exit(exitFailure)
	}
	if *channel == "" {
		fmt.Fprintf(os.Stderr, "Error: -c/--channel is required\n")
		os.Exit(exitFailure)
	}
	if msg == "" {
		fmt.Fprintf(os.Stderr, "Error: message is required (provide via -m/--message or stdin)\n")
		os.Exit(exitFailure)
	}

	// Validate channel
	if *channel != "sms" && *channel != "email" {
		fmt.Fprintf(os.Stderr, "Error: channel must be 'sms' or 'email', got '%s'\n", *channel)
		os.Exit(exitFailure)
	}

	// Create app and send alert
	notifoxApp, err := app.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing app: %v\n", err)
		os.Exit(exitFailure)
	}

	err = notifoxApp.SendAlert(*audience, *channel, msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending alert: %v\n", err)
		os.Exit(exitFailure)
	}

	fmt.Println("Alert sent successfully!")
	os.Exit(exitSuccess)
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

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: notifox <command> [flags]\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Commands:\n")
	fmt.Fprintf(os.Stderr, "  send    Send an alert\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Use 'notifox send -h' for help on the send command\n")
}
