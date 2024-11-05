package adapters

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/EssaAlshammri/github-activity/domain"
)

type GithubAdapter struct {
	client *http.Client
}

func NewGithubAdapter() *GithubAdapter {
	return &GithubAdapter{
		client: &http.Client{},
	}
}

func (ga *GithubAdapter) GetUserActivity(username string) (domain.Events, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/events?per_page=100", username)
	resp, err := ga.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch events: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var events domain.Events
	if err := json.Unmarshal(body, &events); err != nil {
		return nil, err
	}

	return events, nil
}

func (ga *GithubAdapter) GetActivitySummary(username string) (domain.ActivitySummaries, error) {
	events, err := ga.GetUserActivity(username)
	if err != nil {
		return nil, err
	}

	summaryMap := make(map[string]map[string]int)
	for _, event := range events {
		if _, exists := summaryMap[event.Type]; !exists {
			summaryMap[event.Type] = make(map[string]int)
		}
		summaryMap[event.Type][event.Repo.Name]++
	}

	var summaries domain.ActivitySummaries
	for eventType, repos := range summaryMap {
		for repoName, count := range repos {
			description := formatDescription(eventType, count, repoName)
			summaries = append(summaries, domain.ActivitySummary{
				Description: description,
				Count:       count,
				RepoName:    repoName,
			})
		}
	}

	return summaries, nil
}

func formatDescription(eventType string, count int, repoName string) string {
	switch eventType {
	case "PushEvent":
		return fmt.Sprintf("Pushed %d commit(s) to %s", count, repoName)
	case "IssuesEvent":
		return fmt.Sprintf("Opened %d issue(s) in %s", count, repoName)
	case "WatchEvent":
		return fmt.Sprintf("Starred %s", repoName)
	case "CreateEvent":
		return fmt.Sprintf("Created %d branch(es) in %s", count, repoName)
	case "PullRequestEvent":
		return fmt.Sprintf("Opened %d pull request(s) in %s", count, repoName)
	case "PullRequestReviewEvent":
		return fmt.Sprintf("Reviewed %d pull request(s) in %s", count, repoName)
	case "IssueCommentEvent":
		return fmt.Sprintf("Commented %d time(s) on issues in %s", count, repoName)
	default:
		return fmt.Sprintf("%s %d time(s) in %s", eventType, count, repoName)
	}
}
