// client_test.go
package github

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetUserActivity(t *testing.T) {
	tests := []struct {
		name       string
		response   string
		statusCode int
		wantErr    bool
		wantEvents int
	}{
		{
			name: "success",
			response: `[{
                "id": "1",
                "type": "PushEvent",
                "repo": {
                    "id": 1,
                    "name": "user/repo",
                    "url": "https://api.github.com/repos/user/repo"
                },
                "created_at": "2024-03-20T10:00:00Z"
            }]`,
			statusCode: http.StatusOK,
			wantErr:    false,
			wantEvents: 1,
		},
		{
			name:       "api error",
			response:   `{"message": "Not Found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
			wantEvents: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if !strings.HasSuffix(r.URL.Path, "/events") {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			// Create client with test server URL
			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			events, err := client.GetUserActivity("testuser")

			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserActivity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(events) != tt.wantEvents {
				t.Errorf("GetUserActivity() got %d events, want %d", len(events), tt.wantEvents)
			}
		})
	}
}

func TestGetActivitySummary(t *testing.T) {
	tests := []struct {
		name        string
		response    string
		statusCode  int
		wantErr     bool
		wantSummary int
	}{
		{
			name: "success multiple events",
			response: `[{
                "id": "1",
                "type": "PushEvent",
                "repo": {
                    "id": 1,
                    "name": "user/repo",
                    "url": "https://api.github.com/repos/user/repo"
                },
                "created_at": "2024-03-20T10:00:00Z"
            }]`,
			statusCode:  http.StatusOK,
			wantErr:     false,
			wantSummary: 1,
		},
		{
			name:        "empty response",
			response:    `[]`,
			statusCode:  http.StatusOK,
			wantErr:     false,
			wantSummary: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			summaries, err := client.GetActivitySummary("testuser")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetActivitySummary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(summaries) != tt.wantSummary {
				t.Errorf("GetActivitySummary() got %d summaries, want %d", len(summaries), tt.wantSummary)
			}
		})
	}
}

func TestFormatEventDescription(t *testing.T) {
	tests := []struct {
		eventType string
		count     int
		repoName  string
		want      string
	}{
		{
			eventType: "PushEvent",
			count:     2,
			repoName:  "user/repo",
			want:      "Pushed 2 commit(s) to user/repo",
		},
		{
			eventType: "WatchEvent",
			count:     1,
			repoName:  "user/repo",
			want:      "Starred user/repo",
		},
		{
			eventType: "IssuesEvent",
			count:     3,
			repoName:  "user/repo",
			want:      "Opened 3 issue(s) in user/repo",
		},
		{
			eventType: "PullRequestEvent",
			count:     1,
			repoName:  "user/repo",
			want:      "Opened 1 pull request(s) in user/repo",
		},
		{
			eventType: "UnknownEvent",
			count:     1,
			repoName:  "user/repo",
			want:      "UnknownEvent 1 time(s) in user/repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.eventType, func(t *testing.T) {
			got := formatEventDescription(tt.eventType, tt.count, tt.repoName)
			if got != tt.want {
				t.Errorf("formatEventDescription() = %v, want %v", got, tt.want)
			}
		})
	}
}
