package models

type PostCreateRequest struct {
	Title       string   `json:"title" binding:"required"`
	Company     string   `json:"company" binding:"required"`
	Location    string   `json:"location" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Type        JobType  `json:"type" binding:"required,jobtype"`
	MaxSalary   int      `json:"max_salary" binding:"required,gt=0"`
	MinSalary   int      `json:"min_salary" binding:"required,gt=0"`
	Perks       []string `json:"perks"`
	Extras      string   `json:"extras"`
}

type PostPutRequest struct {
	Title       string   `json:"title"`
	Company     string   `json:"company"`
	Location    string   `json:"location"`
	Description string   `json:"description"`
	Type        JobType  `json:"type" binding:"omitempty,jobtype"`
	MaxSalary   int      `json:"max_salary" binding:"omitempty,gt=0"`
	MinSalary   int      `json:"min_salary" binding:"omitempty,gt=0"`
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

type PostOrderBy struct {
	CreatedAt    string `json:"created_at"`
	Salary       string `json:"salary"`
	CompanyCount int    `json:"company_count"`
}
