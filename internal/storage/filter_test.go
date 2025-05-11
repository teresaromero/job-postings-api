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
			name: "match company lowecase",
			job:  testJob,
			filter: models.ListRequestQueryParams{
				Company: "test company",
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
			name: "match title partial",
			job:  testJob,
			filter: models.ListRequestQueryParams{
				Title: "Software",
			},
			want: true,
		},
		{
			name: "match title",
			job:  testJob,
			filter: models.ListRequestQueryParams{
				Title: "Software Engineer",
			},
			want: true,
		},
		{
			name: "match title partial lowercase",
			job:  testJob,
			filter: models.ListRequestQueryParams{
				Title: "software",
			},
			want: true,
		},
		{
			name: "match title lowercase",
			job:  testJob,
			filter: models.ListRequestQueryParams{
				Title: "software engineer",
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
			name: "match location lowercase",
			job:  testJob,
			filter: models.ListRequestQueryParams{
				Location: "remote",
			},
			want: true,
		},
		{
			name: "no match max salary with salary range",
			job: models.Post{
				Salary: []int{50000, 70000},
			},
			filter: models.ListRequestQueryParams{
				MaxSalary: 7000,
			},
			want: false,
		},
		{
			name: "match max salary with salary range",
			job: models.Post{
				Salary: []int{50000, 70000},
			},
			filter: models.ListRequestQueryParams{
				MaxSalary: 80000,
			},
			want: true,
		},
		{
			name: "no match min salary with salary range",
			job: models.Post{
				Salary: []int{50000, 70000},
			},
			filter: models.ListRequestQueryParams{
				MinSalary: 80000,
			},
			want: false,
		},
		{
			name: "match min salary with salary range",
			job: models.Post{
				Salary: []int{50000, 70000},
			},
			filter: models.ListRequestQueryParams{
				MinSalary: 60000,
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
