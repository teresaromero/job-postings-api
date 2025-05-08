package models

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestJobType_IsValid(t *testing.T) {
	tests := []struct {
		name    string
		jobType JobType
		want    bool
	}{
		{
			name:    "valid full-time job type",
			jobType: FullTime,
			want:    true,
		},
		{
			name:    "valid part-time job type",
			jobType: PartTime,
			want:    true,
		},
		{
			name:    "valid contract job type",
			jobType: Contract,
			want:    true,
		},
		{
			name:    "valid internship job type",
			jobType: Internship,
			want:    true,
		},
		{
			name:    "valid freelance job type",
			jobType: Freelance,
			want:    true,
		},
		{
			name:    "invalid job type",
			jobType: "invalid",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.jobType.IsValid(); got != tt.want {
				t.Errorf("JobType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJobTypeValidator(t *testing.T) {
	validate := validator.New()
	_ = validate.RegisterValidation("jobtype", JobTypeValidator)

	type TestStruct struct {
		Job JobType `validate:"jobtype"`
	}

	tests := []struct {
		name    string
		input   TestStruct
		wantErr bool
	}{
		{
			name:    "valid full-time job type",
			input:   TestStruct{Job: FullTime},
			wantErr: false,
		},
		{
			name:    "valid part-time job type",
			input:   TestStruct{Job: PartTime},
			wantErr: false,
		},
		{
			name:    "valid contract job type",
			input:   TestStruct{Job: Contract},
			wantErr: false,
		},
		{
			name:    "valid internship job type",
			input:   TestStruct{Job: Internship},
			wantErr: false,
		},
		{
			name:    "valid freelance job type",
			input:   TestStruct{Job: Freelance},
			wantErr: false,
		},
		{
			name:    "invalid job type",
			input:   TestStruct{Job: "invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("JobTypeValidator() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
