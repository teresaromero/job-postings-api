package sorting

import (
	"job-postings-api/internal/models"
	"slices"
)

type Sorter struct {
	order []string
	rules map[string]func(input *models.SortInput)
}

func NewSorter(order []string) *Sorter {
	defaultOrder := []string{"0", "1", "2"}
	if len(order) == 0 {
		order = defaultOrder
	}

	return &Sorter{
		order: order,
		// define sorting rules with initial order of priority
		// 2. Posts created in the last 7 days rank first above older posts
		// 1. Posts with higher salaries rank first above posts with lower salaries
		// 0. Posts made by companies with more open job posts rank first than those for companies with less posts
		rules: map[string]func(input *models.SortInput){
			"0": sortByCompanyPostsCount,
			"1": sortByHighSalary,
			"2": sortByLastSevenDays,
		},
	}
}

func (s *Sorter) Sort(input *models.SortInput) {
	for _, idx := range s.order {
		s.rules[idx](input)
	}
}

func sortByCompanyPostsCount(input *models.SortInput) {
	// sort by companies with more open job posts first
	slices.SortStableFunc(input.Posts, func(a, b *models.Post) int {
		if input.CompanyMap[a.Company] != input.CompanyMap[b.Company] {
			return input.CompanyMap[b.Company] - input.CompanyMap[a.Company]
		}
		return 0
	})
}

func sortByHighSalary(input *models.SortInput) {
	// sort by the salary first, highter salary first
	slices.SortStableFunc(input.Posts, func(a, b *models.Post) int {
		return b.Salary[1] - a.Salary[1]
	})
}

func sortByLastSevenDays(input *models.SortInput) {
	// sort by the created at date, newer posts first
	slices.SortStableFunc(input.Posts, func(a, b *models.Post) int {
		// if is not created in the last 7 days, keep the order
		if !a.CreatedAtLastSevenDays() && !b.CreatedAtLastSevenDays() {
			return 0
		}

		// if both are created out of the  days, sort by created at date
		if a.CreatedAt != b.CreatedAt {
			return int(b.CreatedAt.Sub(a.CreatedAt).Seconds())
		}
		return 0
	})
}
