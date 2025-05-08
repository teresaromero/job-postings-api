package models

const (
	EmployerUserType = "employer"
	EmployeeUserType = "employee"
)

type UserType string

type User struct {
	Type UserType `json:"type"`
}
