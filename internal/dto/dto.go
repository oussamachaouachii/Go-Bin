package dto

type SnippetCreateDTO struct {
	Title   string `validate:"required,min=4,max=20"`
	Content string `validate:"required,min=4,max=200"`
	Expires int    `validate:"required,oneof=7 30 365"`
}
