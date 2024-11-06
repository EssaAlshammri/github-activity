package github

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type Event struct {
	ID        string     `json:"id"`
	Type      string     `json:"type"`
	Repo      Repository `json:"repo"`
	CreatedAt string     `json:"created_at"`
}

type Repository struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type ActivitySummary struct {
	Description string
	Count       int
	RepoName    string
}

func NewClient() *Client {
	return &Client{
		baseURL:    "https://api.github.com",
		httpClient: &http.Client{},
	}
}

func (c *Client) GetUserActivity(username string) ([]Event, error) {
	url := fmt.Sprintf("%s/users/%s/events?per_page=100", c.baseURL, username)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch events: %s", resp.Status)
	}

	var events []Event
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, err
	}

	return events, nil
}

func (c *Client) GetActivitySummary(username string) ([]ActivitySummary, error) {
	events, err := c.GetUserActivity(username)
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

	var summaries []ActivitySummary
	for eventType, repos := range summaryMap {
		for repoName, count := range repos {
			summaries = append(summaries, ActivitySummary{
				Description: formatEventDescription(eventType, count, repoName),
				Count:       count,
				RepoName:    repoName,
			})
		}
	}

	return summaries, nil
}

func formatEventDescription(eventType string, count int, repoName string) string {
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
