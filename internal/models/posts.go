package models

import "time"

type PostList struct {
	Posts []*Post `json:"posts"`
	Count int     `json:"count"`
}

type Post struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Company     string    `json:"company"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
	Type        JobType   `json:"type"`
	Salary      []int     `json:"salary"`
	Perks       []string  `json:"perks"`
	Extras      string    `json:"extras"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreatedAtLastSevenDays checks if the post was created in the last 7 days
// and returns true if it was, false otherwise.
func (p *Post) CreatedAtLastSevenDays() bool {
	return time.Since(p.CreatedAt).Hours() < 7*24
}
