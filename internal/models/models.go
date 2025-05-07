package models

type UserType string

const (
	EmployerUserType = "employer"
	EmployeeUserType = "employee"
)

type JobType string

const (
	FullTime   JobType = "full-time"
	PartTime   JobType = "part-time"
	Contract   JobType = "contract"
	Internship JobType = "internship"
	Freelance  JobType = "freelance"
)

type Post struct {
	Title       string   `json:"title"`
	Company     string   `json:"company"`
	Location    string   `json:"location"`
	Description string   `json:"description"`
	Type        JobType  `json:"type"`
	Salary      []int    `json:"salary"`
	Perks       []string `json:"perks"`
	Extras      string   `json:"extras"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

type PostOrderBy struct {
	CreatedAt    string `json:"created_at"`
	Salary       string `json:"salary"`
	CompanyCount int    `json:"company_count"`
}

type PostFilter struct {
	Title    string `json:"title"`
	Location string `json:"location"`
	Salary   []int  `json:"salary_range"`
}

type User struct {
	Type UserType `json:"type"`
}
