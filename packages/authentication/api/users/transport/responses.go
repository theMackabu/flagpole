package transport

type LoginResponse struct{
	Token string `json:"token"`
}

type SignupResponse struct{
	Created bool `json:"created"`
	Message string `json:"message"`
}
