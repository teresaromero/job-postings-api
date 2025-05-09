package sorting

import (
	"job-postings-api/internal/models"
	"slices"
)

// SortPosts sorts the posts by the following rules:
// 1. Posts created in the last 7 days rank first above older posts
// 2. Posts with higher salaries rank first above posts with lower salaries
// 3. Posts made by companies with more open job posts rank first than those for companies with less posts
func SortPosts(posts []*models.Post, companyMap map[string]int) {
	// sort by companies with more open job posts first
	slices.SortFunc(posts, func(a, b *models.Post) int {
		if companyMap[a.Company] != companyMap[b.Company] {
			return companyMap[b.Company] - companyMap[a.Company]
		}
		return 0
	})

	// sort by the salary first, highter salary first
	slices.SortFunc(posts, func(a, b *models.Post) int {
		if a.Salary[1] != b.Salary[1] {
			return b.Salary[1] - a.Salary[1]
		}
		return b.Salary[0] - a.Salary[0]
	})

	// sort by the created at date, newer posts first
	slices.SortFunc(posts, func(a, b *models.Post) int {
		if a.CreatedAt != b.CreatedAt {
			return int(b.CreatedAt.Sub(a.CreatedAt).Seconds())
		}
		return 0
	})
}
