package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/EssaAlshammri/github-activity/adapters"
	"github.com/EssaAlshammri/github-activity/service"
)

const usage = `GitHub User Activity CLI

Usage:
    github-activity [--format=<format>] <username>

Options:
    --format    Output format: 'summary' or 'all' (default: summary)

Example:
    github-activity EssaAlshammri
    github-activity --format=all EssaAlshammari`

func main() {
	format := flag.String("format", "summary", "Output format: 'summary' or 'all'")
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Println(usage)
		os.Exit(1)
	}

	username := args[0]
	githubAdapter := adapters.NewGithubAdapter()
	activityService := service.NewActivityService(githubAdapter)

	switch *format {
	case "summary":
		summaries, err := activityService.FetchActivitySummary(username)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if len(summaries) == 0 {
			fmt.Printf("No recent activity found for user: %s\n", username)
			return
		}

		fmt.Printf("Recent activity summary for %s:\n\n", username)
		for _, summary := range summaries {
			fmt.Printf("- %s\n", summary.Description)
		}

	case "all":
		events, err := activityService.FetchUserActivity(username)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if len(events) == 0 {
			fmt.Printf("No recent activity found for user: %s\n", username)
			return
		}

		fmt.Printf("All recent activities for %s:\n\n", username)
		for _, event := range events {
			fmt.Printf("- [%s] %s: %s\n", event.CreatedAt, event.Type, event.Repo.Name)
		}

	default:
		fmt.Printf("Invalid format: %s. Use 'summary' or 'all'\n", *format)
		os.Exit(1)
	}
}
