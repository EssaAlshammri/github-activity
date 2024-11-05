package domain

type Events []Event

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
