package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDuplicateEmailOrPhone = errors.New("邮箱或电话冲突")
	ErrRecordNotFound        = sql.ErrNoRows
)

type UserDAO interface {
	Insert(ctx context.Context, u UserInfo) error
	FindByEmail(ctx context.Context, email string) (UserInfo, error)
	UpdateById(ctx context.Context, entity UserInfo) error
	FindById(ctx context.Context, uid int64) (UserInfo, error)
	FindByPhone(ctx context.Context, phone string) (UserInfo, error)
}

type UserInfo struct {
	Id       int64          `gorm:"column:id;primaryKey;not null;autoIncrement;comment:自增ID"`
	Email    sql.NullString `gorm:"column:email;unique;comment:邮件"`
	Password string         `gorm:"column:password;comment:加密过后的密码"`
	Name     string         `gorm:"column:name;comment:加密过后的密码"`
	Phone    sql.NullString `gorm:"column:phone;unique;comment:电话"`
	Status   string         `gorm:"column:status;comment:账户状态"`
	CTime    *time.Time     `gorm:"column:c_time;comment:创建时间"`
	UTime    *time.Time     `gorm:"column:u_time;comment:更新时间"`
}

func NewOrmUser(db *gorm.DB) UserDAO {
	return &OrmUser{
		db: db,
	}
}

type OrmUser struct {
	db *gorm.DB
}

func (o *OrmUser) Insert(ctx context.Context, u UserInfo) error {
	now := time.Now()
	u.CTime = &now
	u.UTime = &now
	err := o.db.WithContext(ctx).Create(&u).Error
	if pe, ok := err.(*pq.Error); ok {
		const UniqueViolation pq.ErrorCode = "23505"
		if pe.Code == UniqueViolation {
			// 用户冲突，邮箱冲突
			return ErrDuplicateEmailOrPhone
		}
	}
	return err
}

func (o *OrmUser) FindByEmail(ctx context.Context, email string) (UserInfo, error) {
	var u UserInfo
	err := o.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err
}

func (o *OrmUser) UpdateById(ctx context.Context, entity UserInfo) error {
	return o.db.WithContext(ctx).Model(&entity).Where("id = ?", entity.Id).Updates(
		map[string]any{
			"u_time": time.Now(),
			"name":   entity.Name,
		},
	).Error
}

func (o *OrmUser) FindById(ctx context.Context, uid int64) (UserInfo, error) {
	var res UserInfo
	err := o.db.WithContext(ctx).Where("id = ?", uid).First(&res).Error
	return res, err
}

func (o *OrmUser) FindByPhone(ctx context.Context, phone string) (UserInfo, error) {
	var res UserInfo
	err := o.db.WithContext(ctx).Where("phone = ?", phone).First(&res).Error
	return res, err
}

var _ UserDAO = &OrmUser{}
