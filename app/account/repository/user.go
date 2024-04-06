package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/NotFound1911/filestore/app/account/domain"
	"github.com/NotFound1911/filestore/app/account/repository/dao"
)

var (
	ErrDuplicateUser = dao.ErrDuplicateEmailOrPhone
	ErrUserNotFound  = dao.ErrRecordNotFound
)

//go:generate mockgen -source=./user.go -package=repomocks -destination=./mocks/user.mock.go UserRepository
type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	UpdateNonZeroFields(ctx context.Context, user domain.User) error
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindById(ctx context.Context, uid int64) (domain.User, error)
}
type CachedUserRepository struct {
	dao dao.UserDAO
	//cache   cache.UserCache
	logFunc func(format string, a ...any)
}

func (repo *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, repo.toEntity(u))
}

func (repo *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), err
}

func (repo *CachedUserRepository) UpdateNonZeroFields(ctx context.Context, user domain.User) error {
	err := repo.dao.UpdateById(ctx, repo.toEntity(user))
	if err != nil {
		return err
	}
	// todo delete cache
	return nil
}

func (repo *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := repo.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *CachedUserRepository) FindById(ctx context.Context, uid int64) (domain.User, error) {
	// todo cache
	// todo 降级
	u, err := repo.dao.FindById(ctx, uid)
	if err != nil {
		return domain.User{}, err
	}
	du := repo.toDomain(u)
	// todo cache
	return du, nil
}

// NewCachedUserRepository todo 缓存模式
func NewCachedUserRepository(
	dao dao.UserDAO,
// c cache.UserCache
) UserRepository {
	return &CachedUserRepository{
		dao: dao,
		//cache: c,
		logFunc: func(format string, a ...any) {
			fmt.Printf(format, a)
		},
	}
}
func (repo *CachedUserRepository) toDomain(u dao.UserInfo) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Phone:    u.Phone.String,
		Password: u.Password,
		Name:     u.Name,
		Status:   u.Status,
		CreateAt: u.CTime,
		UpdateAt: u.UTime,
	}
}

func (repo *CachedUserRepository) toEntity(u domain.User) dao.UserInfo {
	return dao.UserInfo{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password: u.Password,
		Name:     u.Name,
		Status:   u.Status,
		CTime:    u.CreateAt,
		UTime:    u.UpdateAt,
	}
}
