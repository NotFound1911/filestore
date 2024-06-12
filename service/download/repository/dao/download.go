package dao

import (
	"context"
	"github.com/NotFound1911/filestore/service/download/domain"
	"gorm.io/gorm"
	"time"
)

type DownloadDao interface {
	Insert(ctx context.Context, d DownloadInfo) (int64, error)
	UpdateById(ctx context.Context, id int64, columns map[string]interface{}) error
	FindById(ctx context.Context, id int64) (DownloadInfo, error)
	GetDownloadInfosByUid(ctx context.Context, uid int64) ([]DownloadInfo, error)
}

type DownloadInfo struct {
	Id       int64         `gorm:"column:id;primaryKey;not null;autoIncrement;comment:自增ID"`
	UId      int64         `gorm:"column:uid;not null;comment:用户id"`
	FileName string        `gorm:"column:file_name;not null;comment:文件名称"`
	FileSize int64         `gorm:"column:file_size;not null;comment:文件size"`
	CTime    *time.Time    `gorm:"column:c_time;comment:下载开始时间"`
	UTime    *time.Time    `gorm:"column:u_time;comment:下载结束时间"`
	Status   domain.Status `gorm:"column:status;comment:下载状态"`
}

func NewOrmDownload(db *gorm.DB) DownloadDao {
	return &OrmDownload{
		db: db,
	}
}

type OrmDownload struct {
	db *gorm.DB
}

func (o *OrmDownload) Insert(ctx context.Context, d DownloadInfo) (int64, error) {
	now := time.Now()
	d.CTime = &now
	err := o.db.WithContext(ctx).Create(&d).Error
	return d.Id, err
}

func (o *OrmDownload) UpdateById(ctx context.Context, id int64, columns map[string]interface{}) error {
	return o.db.WithContext(ctx).Model(&DownloadInfo{}).Where("id = ?", id).Updates(columns).Error

}

func (o *OrmDownload) FindById(ctx context.Context, id int64) (DownloadInfo, error) {
	var res DownloadInfo
	err := o.db.WithContext(ctx).Where("id = ?", id).First(&res).Error
	return res, err
}

func (o *OrmDownload) GetDownloadInfosByUid(ctx context.Context, uid int64) ([]DownloadInfo, error) {
	var res []DownloadInfo
	err := o.db.WithContext(ctx).Where("uid = ?", uid).Find(&res).Error
	return res, err
}
