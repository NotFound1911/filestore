package domain

import "time"

type FileMeta struct {
	Id          int64  `json:"id"`           // 文件id
	Sha1        string `json:"sha1"`         // 文件sha1
	Size        int64  `json:"size"`         // 文件size
	Address     string `json:"address"`      // 文件本地存储位置
	Type        string `json:"type"`         // 文件类别
	Bucket      string `json:"bucket"`       // 桶
	StorageName string `json:"storage_name"` // 存储名称
}

type UserFile struct {
	Id       int64      `json:"id"` // 用户文件id
	UserId   int64      `json:"user_id"`
	FileName string     `json:"file_name"` // 文件名
	FileSha1 string     `json:"file_sha1"` // 文件hash
	FileSize int64      `json:"file_size"`
	UpdateAt *time.Time `json:"update_at"` // 更新时间
}
