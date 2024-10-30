package models

type Person struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Type      string `json:"type" validate:"required,ValidateType"`
	Age       int    `json:"age" validate:"required,gt=0"`
	Courses   []int  `json:"courses" validate:"required"`
}
