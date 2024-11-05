package ports

import "github.com/EssaAlshammri/github-activity/domain"

// Secondary Port (driven)
type GithubActivityPort interface {
	GetUserActivity(username string) (domain.Events, error)
	GetActivitySummary(username string) (domain.ActivitySummaries, error)
}

// Primary Port (driving)
type ActivityService interface {
	FetchUserActivity(username string) (domain.Events, error)
	FetchActivitySummary(username string) (domain.ActivitySummaries, error)
}
