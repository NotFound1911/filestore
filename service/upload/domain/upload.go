package domain

import "time"

type Upload struct {
	Id       int64      `json:"id"`        // 上传id
	UId      int64      `json:"u_id"`      // 用户id
	FileName string     `json:"file_name"` // 文件名称
	FileSha1 string     `json:"file_sha1"` // 文件sha1
	FileSize int64      `json:"file_size"` // 文件size
	CreateAt *time.Time `json:"create_at"` // 上传开始时间
	UpdateAt *time.Time `json:"update_at"` // 上传更新时间
	Status   string     `json:"status"`    // 上传状态
}
