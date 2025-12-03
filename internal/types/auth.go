package types

type RegisterUserInput struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type RegisterUserOutput struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}



type LoginOutput struct {
	AccessToken string `json:"access_token"`
}

type VerifyUserInput struct {
	Email string `json:"email"`
	Token string `json:"token"`
}


type MessageOutput struct {
	Message string `json:"message"`
}

type ResendVerificationInput struct {
	Email string `json:"email"`

}

type ForgotPasswordInput struct {
	Email string `json:"email"`
}



type ResetPasswordInput struct {
    Token string `json:"token"`
	Password string `json:"password"`
}