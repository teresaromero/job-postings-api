package sorting

import (
	"job-postings-api/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSortPosts(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name       string
		posts      []*models.Post
		companyMap map[string]int
		want       []*models.Post
	}{
		{
			name: "sorts by company posts count",
			posts: []*models.Post{
				{Company: "A", Salary: []int{50000, 60000}, CreatedAt: now},
				{Company: "B", Salary: []int{50000, 60000}, CreatedAt: now},
				{Company: "C", Salary: []int{50000, 60000}, CreatedAt: now},
			},
			companyMap: map[string]int{
				"A": 1,
				"B": 3,
				"C": 2,
			},
			want: []*models.Post{
				{Company: "B", Salary: []int{50000, 60000}, CreatedAt: now},
				{Company: "C", Salary: []int{50000, 60000}, CreatedAt: now},
				{Company: "A", Salary: []int{50000, 60000}, CreatedAt: now},
			},
		},
		{
			name: "sorts by recent first posts count",
			posts: []*models.Post{
				{Company: "A", Salary: []int{50000, 60000}, CreatedAt: now.Add(-8 * 24 * time.Hour)},
				{Company: "A", Salary: []int{50000, 60000}, CreatedAt: now.Add(-2 * time.Hour)},
				{Company: "A", Salary: []int{50000, 60000}, CreatedAt: now},
			},
			companyMap: map[string]int{
				"A": 1,
				"B": 3,
				"C": 2,
			},
			want: []*models.Post{
				{Company: "A", Salary: []int{50000, 60000}, CreatedAt: now},
				{Company: "A", Salary: []int{50000, 60000}, CreatedAt: now.Add(-2 * time.Hour)},
				{Company: "A", Salary: []int{50000, 60000}, CreatedAt: now.Add(-8 * 24 * time.Hour)}},
		},
		{
			name: "sorts by higher salary first",
			posts: []*models.Post{
				{Company: "A", Salary: []int{50000, 60000}, CreatedAt: now},
				{Company: "A", Salary: []int{70000, 80000}, CreatedAt: now},
				{Company: "A", Salary: []int{60000, 70000}, CreatedAt: now},
			},
			companyMap: map[string]int{
				"A": 1,
				"B": 3,
				"C": 2,
			},
			want: []*models.Post{
				{Company: "A", Salary: []int{70000, 80000}, CreatedAt: now},
				{Company: "A", Salary: []int{60000, 70000}, CreatedAt: now},
				{Company: "A", Salary: []int{50000, 60000}, CreatedAt: now},
			},
		},
		{
			name: "sorts by higher salary first with different companies",
			posts: []*models.Post{
				{Company: "A", Salary: []int{50000, 60000}, CreatedAt: now},
				{Company: "B", Salary: []int{60000, 70000}, CreatedAt: now},
				{Company: "C", Salary: []int{70000, 80000}, CreatedAt: now},
			},
			companyMap: map[string]int{
				"A": 1,
				"B": 3,
				"C": 2,
			},
			want: []*models.Post{
				{Company: "C", Salary: []int{70000, 80000}, CreatedAt: now},
				{Company: "B", Salary: []int{60000, 70000}, CreatedAt: now},
				{Company: "A", Salary: []int{50000, 60000}, CreatedAt: now},
			},
		},
		{
			name: "sorts by higher salary first with different companies and different created at",
			posts: []*models.Post{
				{Company: "A", Salary: []int{50000, 60000}, CreatedAt: now.Add(-8 * 24 * time.Hour)},
				{Company: "B", Salary: []int{60000, 70000}, CreatedAt: now.Add(-2 * time.Hour)},
				{Company: "C", Salary: []int{70000, 80000}, CreatedAt: now},
			},
			companyMap: map[string]int{
				"A": 1,
				"B": 3,
				"C": 2,
			},
			want: []*models.Post{
				{Company: "C", Salary: []int{70000, 80000}, CreatedAt: now},
				{Company: "B", Salary: []int{60000, 70000}, CreatedAt: now.Add(-2 * time.Hour)},
				{Company: "A", Salary: []int{50000, 60000}, CreatedAt: now.Add(-8 * 24 * time.Hour)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SortPosts(tt.posts, tt.companyMap)

			assert.Equal(t, tt.want, tt.posts, "posts should be sorted correctly")
		})
	}
}
