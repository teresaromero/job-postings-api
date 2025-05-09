package storage

import (
	"job-postings-api/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isAllowedByFilter(t *testing.T) {
	testJob := models.Post{
		Title:    "Software Engineer",
		Company:  "Test Company",
		Location: "Remote",
		Salary:   []int{50000, 80000},
	}

	tests := []struct {
		name   string
		job    models.Post
		filter models.ListRequestQueryParams
		want   bool
	}{
		{
			name: "match company",
			job:  testJob,
			filter: models.ListRequestQueryParams{
				Company: "Test Company",
			},
			want: true,
		},
		{
			name: "no match company",
			job:  testJob,
			filter: models.ListRequestQueryParams{
				Company: "Other Company",
			},
			want: false,
		},
		{
			name: "match title",
			job:  testJob,
			filter: models.ListRequestQueryParams{
				Title: "Software",
			},
			want: true,
		},
		{
			name: "no match title",
			job:  testJob,
			filter: models.ListRequestQueryParams{
				Title: "Manager",
			},
			want: false,
		},
		{
			name: "match location",
			job:  testJob,
			filter: models.ListRequestQueryParams{
				Location: "Remote",
			},
			want: true,
		},
		{
			name: "match max salary",
			job:  testJob,
			filter: models.ListRequestQueryParams{
				MaxSalary: 60000,
			},
			want: true,
		},
		{
			name: "match min salary",
			job:  testJob,
			filter: models.ListRequestQueryParams{
				MinSalary: 70000,
			},
			want: true,
		},
		{
			name: "no matches",
			job:  testJob,
			filter: models.ListRequestQueryParams{
				Company:   "Other",
				Title:     "Manager",
				Location:  "Office",
				MinSalary: 90000,
				MaxSalary: 40000,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isAllowedByFilter(tt.job, tt.filter)
			assert.Equal(t, tt.want, got)
		})
	}
}
