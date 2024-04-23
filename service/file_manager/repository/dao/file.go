package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type FileManagerDao interface {
	InsertFileMeta(context.Context, FileMetaInfo) (int64, error)
	InsertUserFile(ctx context.Context, info UserFileInfo) (int64, error)
	UpdateFileMetaById(ctx context.Context, id int64, columns map[string]interface{}) error
	UpdateUserFileInfo(ctx context.Context, id int64, columns map[string]interface{}) error
	FindFileMetaById(ctx context.Context, id int64) (FileMetaInfo, error)
	FindUserFileById(ctx context.Context, id int64) (UserFileInfo, error)
	GetFileMetasByUserId(ctx context.Context, uid int64) ([]FileMetaInfo, error)
	GetFileMeta(ctx context.Context, sha1 string) (FileMetaInfo, error)
	GetUserIdsByFileSha1(ctx context.Context, sha1 string) ([]int64, error)
}

type FileMetaInfo struct {
	Id      int64  `gorm:"column:id;primaryKey;not null;autoIncrement;comment:自增ID"`
	Sha1    string `gorm:"column:sha1;not null;unique;comment:文件sha1"`
	Size    int64  `gorm:"column:size;not null;comment:文件size"`
	Address string `gorm:"column:address;comment:文件存储地址"`
	Type    string `gorm:"column:type;comment:文件类型"`
}

type UserFileInfo struct {
	Id       int64      `gorm:"column:id;primaryKey;not null;autoIncrement;comment:自增ID"`
	UserId   int64      `gorm:"column:user_id;not null;comment:用户id"`
	FileName string     `gorm:"column:file_name;not null;comment:文件名"`
	FileSha1 string     `gorm:"column:file_sha1;not null;comment:文件sha1"`
	FileSize int64      `gorm:"column:file_size;not null;comment:文件size"`
	UpdateAt *time.Time `gorm:"column:update_at;not null;comment:更新时间"`
}

func NewOrmFileManager(db *gorm.DB) FileManagerDao {
	return &OrmFileManager{
		db: db,
	}
}

type OrmFileManager struct {
	db *gorm.DB
}

func (o *OrmFileManager) InsertFileMeta(ctx context.Context, info FileMetaInfo) (int64, error) {
	err := o.db.WithContext(ctx).Create(&info).Error
	return info.Id, err
}

func (o *OrmFileManager) InsertUserFile(ctx context.Context, info UserFileInfo) (int64, error) {
	now := time.Now()
	info.UpdateAt = &now
	err := o.db.WithContext(ctx).Create(&info).Error
	return info.Id, err
}

func (o *OrmFileManager) UpdateFileMetaById(ctx context.Context, id int64, columns map[string]interface{}) error {
	return o.db.WithContext(ctx).Model(&FileMetaInfo{}).Where("id = ?", id).Updates(columns).Error
}

func (o *OrmFileManager) UpdateUserFileInfo(ctx context.Context, id int64, columns map[string]interface{}) error {
	return o.db.WithContext(ctx).Model(&UserFileInfo{}).Where("id = ?", id).Updates(columns).Error
}

func (o *OrmFileManager) FindFileMetaById(ctx context.Context, id int64) (FileMetaInfo, error) {
	var res FileMetaInfo
	err := o.db.WithContext(ctx).Where("id = ?", id).First(&res).Error
	return res, err
}

func (o *OrmFileManager) FindUserFileById(ctx context.Context, id int64) (UserFileInfo, error) {
	var res UserFileInfo
	err := o.db.WithContext(ctx).Where("id = ?", id).First(&res).Error
	return res, err
}
func (o *OrmFileManager) GetFileMetasByUserId(ctx context.Context, uid int64) ([]FileMetaInfo, error) {
	var fileMetas []FileMetaInfo
	err := o.db.WithContext(ctx).Joins("UserFileInfo").
		Where("user_file_info.file_sha1 = file_meta_info.sha1 and user_file_info.user_id = ?", uid).
		Find(&fileMetas).Error
	return fileMetas, err
}
func (o *OrmFileManager) GetUserIdsByFileSha1(ctx context.Context, sha1 string) ([]int64, error) {
	var userIds []int64
	// 查询满足条件的数据
	err := o.db.WithContext(ctx).Model(&UserFileInfo{}).
		Where("file_sha1 = ?", "sha1").
		Pluck("user_id", &userIds).Error
	return userIds, err
}
func (o *OrmFileManager) GetFileMeta(ctx context.Context, sha1 string) (FileMetaInfo, error) {
	var fileMeta FileMetaInfo
	err := o.db.WithContext(ctx).Where("sha1 = ?", sha1).First(&fileMeta).Error
	return fileMeta, err
}
