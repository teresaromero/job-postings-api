package models

import "github.com/go-playground/validator/v10"

const (
	FullTime   JobType = "full-time"
	PartTime   JobType = "part-time"
	Contract   JobType = "contract"
	Internship JobType = "internship"
	Freelance  JobType = "freelance"
)

type JobType string

func (jt JobType) IsValid() bool {
	switch jt {
	case FullTime, PartTime, Contract, Internship, Freelance:
		return true
	}
	return false
}

func JobTypeValidator(fl validator.FieldLevel) bool {
	jobType, ok := fl.Field().Interface().(JobType)
	if !ok {
		return false
	}
	return jobType.IsValid()
}
