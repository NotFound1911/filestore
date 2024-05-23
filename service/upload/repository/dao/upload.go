package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type UploadDao interface {
	Insert(ctx context.Context, u UploadInfo) (int64, error)
	UpdateById(ctx context.Context, id int64, columns map[string]interface{}) error
	FindById(ctx context.Context, id int64) (UploadInfo, error)
	GetUploadInfosByUid(ctx context.Context, uid int64) ([]UploadInfo, error)
}

type UploadInfo struct {
	Id       int64      `gorm:"column:id;primaryKey;not null;autoIncrement;comment:自增ID"`
	UId      int64      `gorm:"column:uid;not null;comment:用户id"`
	FileName string     `gorm:"column:file_name;not null;comment:文件名称"`
	FileSha1 string     `gorm:"column:file_sha1;not null;comment:文件sha1"` // 不应该unique 应该允许同一文件多次上传
	FileSize int64      `gorm:"column:file_size;not null;comment:文件size"`
	CTime    *time.Time `gorm:"column:c_time;comment:上传开始时间"`
	UTime    *time.Time `gorm:"column:u_time;comment:上传结束时间"`
	Status   string     `gorm:"column:status;comment:上传状态"`
	Type     string     `gorm:"column:type;comment:上传类别"`
}

func NewOrmUpload(db *gorm.DB) UploadDao {
	return &OrmUpload{
		db: db,
	}
}

type OrmUpload struct {
	db *gorm.DB
}

func (o *OrmUpload) Insert(ctx context.Context, u UploadInfo) (int64, error) {
	now := time.Now()
	u.CTime = &now
	err := o.db.WithContext(ctx).Create(&u).Error
	return u.Id, err
}

func (o *OrmUpload) UpdateById(ctx context.Context, id int64, columns map[string]interface{}) error {
	return o.db.WithContext(ctx).Model(&UploadInfo{}).Where("id = ?", id).Updates(columns).Error
}

func (o *OrmUpload) FindById(ctx context.Context, id int64) (UploadInfo, error) {
	var res UploadInfo
	err := o.db.WithContext(ctx).Where("id = ?", id).First(&res).Error
	return res, err
}

func (o *OrmUpload) GetUploadInfosByUid(ctx context.Context, uid int64) ([]UploadInfo, error) {
	var res []UploadInfo
	err := o.db.WithContext(ctx).Where("uid = ?", uid).Find(&res).Error
	return res, err
}
