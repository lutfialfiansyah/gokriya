package usecase

import (
	"context"
	"encoding/json"
	"golang.org/x/sync/errgroup"
	"time"

	"github.com/bxcodec/go-clean-arch/domain"
)

type usersUsecase struct {
	usersRepo    domain.UsersRepository
	//roleRepo     domain.RolesRepository
	contextTimeout time.Duration
}

// NewUsersUsecase will create new an articleUsecase object representation of domain.ArticleUsecase interface
func NewUsersUsecase(a domain.UsersRepository, timeout time.Duration) domain.UsersUsecase {
	return &usersUsecase{
		usersRepo:    a,
		//roleRepo:     ar,
		contextTimeout: timeout,
	}
}
func (u usersUsecase) Fetch(ctx context.Context, page int,size int) ([]*domain.UsersList, error) {
	_, c := errgroup.WithContext(ctx)
	getUsers,err := u.usersRepo.Fetch(c,page,size)
	listUsers := make([]*domain.UsersList,0)
	if err != nil {
		return nil,err
	}
	for _,element := range getUsers{
		var data domain.UsersObj
		if errUnmarshal := json.Unmarshal([]byte(element.Data), &data); errUnmarshal != nil {
			return nil, errUnmarshal
		}
		usr := domain.UsersList{
			Id:       element.ID,
			Username: data.Username,
			Email:    data.Email,
			Status:   data.Status,
		}
		listUsers = append(listUsers,&usr)
	}
	return listUsers,nil
}

func (u usersUsecase) GetByID(ctx context.Context, id string) (*domain.UsersListDetail, error) {
	_, c := errgroup.WithContext(ctx)
	getUsersId,err := u.usersRepo.GetByID(c,id)
	if err != nil {
		return nil,err
	}
	var data domain.UsersObj
	if errUnmarshal := json.Unmarshal([]byte(getUsersId.DataUser), &data); errUnmarshal != nil {
		return nil, errUnmarshal
	}
	var dataRole domain.RolesData
	if errUnmarshal := json.Unmarshal([]byte(getUsersId.DataRole), &dataRole); errUnmarshal != nil {
		return nil, errUnmarshal
	}
	result := domain.UsersListDetail{
		Id:       getUsersId.Id,
		Username: data.Username,
		Email:    data.Email,
		Role:     dataRole.RoleName,
	}

	return &result,nil
}

func (u usersUsecase) Update(ctx context.Context, ar *domain.Users) error {
	panic("implement me")
}

func (u usersUsecase) GetByData(ctx context.Context, title string) (domain.Users, error) {
	panic("implement me")
}

func (u usersUsecase) Store(ctx context.Context, users *domain.Users) error {
	panic("implement me")
}

func (u usersUsecase) Delete(ctx context.Context, id string) error {
	_, c := errgroup.WithContext(ctx)
	err := u.usersRepo.Delete(c,id)
	if err != nil {
		return err
	}
	return nil
}


/*
* In this function below, I'm using errgroup with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */
