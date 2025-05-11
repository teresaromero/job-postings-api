package storage

import (
	"job-postings-api/internal/models"
	"strings"
)

// isAllowedByFilter checks if the job does not match the filter criteria.
func isAllowedByFilter(job models.Post, filter models.ListRequestQueryParams) bool {

	if filter.Company != "" && !strings.EqualFold(filter.Company, job.Company) {
		return false
	}
	if filter.Title != "" && (!strings.EqualFold(job.Title, filter.Title) &&
		!strings.Contains(strings.ToLower(job.Title), strings.ToLower(filter.Title))) {
		return false
	}
	if filter.Location != "" && !strings.EqualFold(job.Location, filter.Location) {
		return false
	}
	if filter.MaxSalary != 0 &&
		job.Salary[1] > filter.MaxSalary {
		return false
	}
	if filter.MinSalary != 0 &&
		job.Salary[0] < filter.MinSalary &&
		job.Salary[1] < filter.MinSalary {
		return false
	}
	return true
}
