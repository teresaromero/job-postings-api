package models

const (
	EmployerUserType UserType = "employer"
	EmployeeUserType UserType = "employee"
)

type UserType string

func (ut UserType) String() string {
	return string(ut)
}

type User struct {
	Type UserType `json:"type"`
}
