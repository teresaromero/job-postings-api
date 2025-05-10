package models

type PostRequestPayload struct {
	Title       string   `json:"title" binding:"required"`
	Company     string   `json:"company" binding:"required"`
	Location    string   `json:"location" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Type        JobType  `json:"type" binding:"required,jobtype"`
	Salary      []int    `json:"salary" binding:"omitempty,dive,gt=0"`
	Perks       []string `json:"perks"`
	Extras      string   `json:"extras"`
}

type ListRequestQueryParams struct {
	Title     string `form:"title"`
	Location  string `form:"location"`
	MaxSalary int    `form:"max_salary"`
	MinSalary int    `form:"min_salary"`
	Company   string `form:"company"`
}

type LoginRequestPayload struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
