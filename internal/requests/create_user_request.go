package requests

type CreateUserRequest struct {
	UserName string `json:"user_name" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}
