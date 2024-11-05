package domain

type ActivitySummary struct {
	Description string
	Count      int
	RepoName   string
}

type ActivitySummaries []ActivitySummary
