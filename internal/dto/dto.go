package dto

type SnippetCreateDTO struct {
	Title   string `validate:"required,min=4,max=20"`
	Content string `validate:"required,min=4,max=200"`
	Expires int    `validate:"required,oneof=7 30 365"`
}

type UserCreateDTO struct {
	Name     string `validate:"required,min=4,max=20"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=30"`
}

type UserLoginDTO struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=30"`
}
