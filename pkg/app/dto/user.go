package dto

type UserDTO struct {
	Id    int64      `json:"id"`
	Cards []*CardDTO `json:"cards"`
}

type UserRegAndLognDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserRegistrationSuccessDTO struct {
	Id int64 `json:"id"`
}

type UserLoginDTO struct {
	Login string `json:"login"`
}
