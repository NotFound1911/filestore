package dao

import (
	"context"
	"errors"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDuplicateSha1 = errors.New("文件sha1冲突")
)

type FileDAO interface {
	Insert(ctx context.Context, f FileMeta) error                 // 插入文件元数据
	UpdateBySha1(ctx context.Context, f FileMeta) error           // 更新元数据
	GetBySha1(ctx context.Context, sha1 string) (FileMeta, error) // 获取元数据
	SoftDeleteBySha1(ctx context.Context, sha1 string) error      // 软删除元数据
}

var _ FileDAO = &OrmFile{}

type OrmFile struct {
	db *gorm.DB
}

func (o *OrmFile) Insert(ctx context.Context, f FileMeta) error {
	now := time.Now()
	f.CTime = &now
	f.UTime = &now
	err := o.db.WithContext(ctx).Create(&f).Error
	if pe, ok := err.(*pq.Error); ok {
		const UniqueViolation pq.ErrorCode = "23505"
		if pe.Code == UniqueViolation {
			return ErrDuplicateSha1
		}
	}
	return err
}

func (o *OrmFile) UpdateBySha1(ctx context.Context, f FileMeta) error {
	return o.db.WithContext(ctx).Model(&f).Where("sha1 = ?", f.Sha1).Updates(
		map[string]any{
			"u_time": time.Now(),
			"name":   f.Name,
			"status": f.Status,
		},
	).Error
}

func (o *OrmFile) GetBySha1(ctx context.Context, sha1 string) (FileMeta, error) {
	var res FileMeta
	err := o.db.WithContext(ctx).Where("sha1 = ?", sha1).First(&res).Error
	return res, err
}

func (o *OrmFile) SoftDeleteBySha1(ctx context.Context, sha1 string) error {
	var res FileMeta
	err := o.db.WithContext(ctx).Model(&res).Update("status", "soft_delete").Error
	return err
}

type FileMeta struct {
	Id      int64      `gorm:"column:id;primaryKey;not null;autoIncrement;comment:自增ID"`
	Sha1    string     `gorm:"column:sha1;unique;not null;comment:元数据sha1"`
	Name    string     `gorm:"column:name;comment:文件名"`
	Size    int64      `gorm:"column:size;comment:文件大小"`
	Address string     `gorm:"column:address;not null;comment:文件大小"`
	Status  string     `gorm:"column:status;comment:文件状态"`
	CTime   *time.Time `gorm:"column:c_time;comment:创建时间"`
	UTime   *time.Time `gorm:"column:u_time;comment:更新时间"`
}
