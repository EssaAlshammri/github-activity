package service

import (
	"github.com/EssaAlshammri/github-activity/domain"
	"github.com/EssaAlshammri/github-activity/ports"
)

type ActivityService struct {
	githubPort ports.GithubActivityPort
}

func NewActivityService(githubPort ports.GithubActivityPort) *ActivityService {
	return &ActivityService{
		githubPort: githubPort,
	}
}

func (s *ActivityService) FetchUserActivity(username string) (domain.Events, error) {
	return s.githubPort.GetUserActivity(username)
}

func (s *ActivityService) FetchActivitySummary(username string) (domain.ActivitySummaries, error) {
	return s.githubPort.GetActivitySummary(username)
}
