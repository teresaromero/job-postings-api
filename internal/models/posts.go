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
