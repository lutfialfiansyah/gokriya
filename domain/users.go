package domain

import "context"

// Users ...
type Users struct {
	ID     string  `json:"id"`
	Data   string `json:"data" validate:"required"`
	RoleId string  `json:"role_id"`
}

type UsersJoinRole struct {
	Id 		string	`json:"id"`
	DataUser string `json:"data_user"`
	DataRole string	`json:"data_role"`
}



type UsersList struct {
	Id 			string `json:"id"`
	Username   	string `json:"username"`
	Email   	string `json:"email"`
	Status 		Status  `json:"status"`
}

type UsersObj struct {
	Email   	string `json:"email"`
	Status 		Status  `json:"status"`
	Password	string	`json:"password"`
	Username   	string `json:"username"`
}

type Status struct {
	IsActive   	bool `json:"is_active"`
}

type UsersListDetail struct {
	Id 			string 	`json:"id"`
	Username   	string 	`json:"username"`
	Email   	string 	`json:"email"`
	Role 		string 	`json:"role_name"`
}

type UsersUsecase interface {
	Fetch(ctx context.Context, page int,size int) ([]*UsersList, error)
	GetByID(ctx context.Context, id string) (*UsersListDetail, error)
	Update(ctx context.Context, ar *Users) error
	GetByData(ctx context.Context, title string) (Users, error)
	Store(context.Context, *Users) error
	Delete(ctx context.Context, id string) error
}

// UsersRepository represent the article's repository contract
type UsersRepository interface {
	Fetch(ctx context.Context, page int,size int) (res []Users, err error)
	GetByID(ctx context.Context, id string) (UsersJoinRole, error)
	Update(ctx context.Context, ar *Users) error
	GetByData(ctx context.Context, title string) (Users, error)
	Store(context.Context, *Users) error
	Delete(ctx context.Context, id string) error
}