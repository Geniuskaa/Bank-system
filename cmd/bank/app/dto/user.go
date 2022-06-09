package dto

type UserDTO struct {
	Id    int64      `json:"id"`
	Cards []*CardDTO `json:"cards"`
}

type NewUserDTO struct {
	Id int64 `json:"id"`
}
